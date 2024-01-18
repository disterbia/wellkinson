// /diet-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"diet-service/dto"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"gorm.io/gorm"
)

type DietService interface {
	SavePreset(presetRequest dto.DietPresetRequest) (string, error)
	GetPresets(id int, page int, startDate, endDate string) ([]dto.DietPresetResponse, error)
	RemovePreset(ids []int, uid int) (string, error)
	SaveDiet(diet dto.DietRequest) (string, error)
	GetDiets(id int, startDate, endDate string) ([]dto.DietResponse, error)
	RemoveDiet(ids []int, uid int) (string, error)
}

type dietService struct {
	db        *gorm.DB
	s3svc     *s3.S3
	bucket    string
	bucketUrl string
}

func NewDietService(db *gorm.DB, s3svc *s3.S3, bucket string, bucketUrl string) DietService {
	return &dietService{db: db, s3svc: s3svc, bucket: bucket, bucketUrl: bucketUrl}
}

func (service *dietService) GetDiets(id int, startDate, endDate string) ([]dto.DietResponse, error) {

	var diet []model.Diet

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")
	result := query.Preload("Images", "level!= 10").Find(&diet)

	if result.Error != nil {
		return nil, result.Error
	}

	var dietResponses []dto.DietResponse
	if err := util.CopyStruct(diet, &dietResponses); err != nil {
		return nil, err
	}

	for i := range dietResponses {
		for j := range dietResponses[i].Images {
			// S3 객체 키를 추출 (URL에서)
			urlkey := extractKeyFromUrl(dietResponses[i].Images[j].Url, service.bucket, service.bucketUrl)
			thumbnailUrlkey := extractKeyFromUrl(dietResponses[i].Images[j].ThumbnailUrl, service.bucket, service.bucketUrl)
			// 사전 서명된 URL을 생성
			url, _ := service.s3svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(service.bucket),
				Key:    aws.String(urlkey),
			})
			thumbnailUrl, _ := service.s3svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(service.bucket),
				Key:    aws.String(thumbnailUrlkey),
			})
			urlStr, err := url.Presign(1 * time.Second) // URL은 1초 동안 유효
			if err != nil {
				return nil, err
			}
			thumbnailUrlStr, err := thumbnailUrl.Presign(1 * time.Second) // URL은 1초 동안 유효
			if err != nil {
				return nil, err
			}
			dietResponses[i].Images[j].Url = urlStr // 사전 서명된 URL로 업데이트
			dietResponses[i].Images[j].ThumbnailUrl = thumbnailUrlStr
		}
	}
	return dietResponses, nil
}

func (service *dietService) SaveDiet(dietRequest dto.DietRequest) (string, error) {
	var diet model.Diet

	// 트랜잭션 시작
	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// 데이터베이스에서 기존 Diet 레코드 조회
	result := tx.Where("id=? AND uid=?", dietRequest.Id, dietRequest.Uid).First(&model.Diet{})
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return "", result.Error
	}

	// DietRequest의 복사본 생성
	dietRequestCopy := dietRequest
	dietRequestCopy.Images = []string{}

	// DietRequest 복사본에서 Diet으로 필드 복사
	if err := util.CopyStruct(dietRequestCopy, &diet); err != nil {
		tx.Rollback()
		return "", err
	}
	diet.Uid = dietRequest.Uid

	// Images 필드 별도 처리
	var temp []model.Image
	for _, imgStr := range dietRequest.Images {
		temp = append(temp, model.Image{Url: imgStr})
	}
	diet.Images = temp

	// 고루틴으로 이미지 처리 및 업로드
	var wg sync.WaitGroup
	images := make([]model.Image, len(dietRequest.Images))
	errorsChan := make(chan error, len(dietRequest.Images))
	uploadedFiles := make(chan string, len(dietRequest.Images)*2)

	uidString := strconv.Itoa(diet.Uid)

	for i, imgStr := range dietRequest.Images {
		wg.Add(1)
		go func(i int, imgStr string) {
			defer wg.Done()

			imgData, err := base64.StdEncoding.DecodeString(imgStr)
			if err != nil {
				errorsChan <- fmt.Errorf("base64 decoding error: %v", err)
				return
			}

			contentType, ext, err := getImageFormat(imgData)
			if err != nil {
				errorsChan <- fmt.Errorf("invalid image format: %v", err)
				return
			}

			// 이미지 크기 조정 (10MB 제한)
			if len(imgData) > 10*1024*1024 {
				imgData, err = reduceImageSize(imgData)
				if err != nil {
					errorsChan <- fmt.Errorf("error reducing image size: %v", err)
					return
				}
			}

			// 썸네일 이미지 생성
			thumbnailData, err := createThumbnail(imgData)
			if err != nil {
				errorsChan <- fmt.Errorf("thumbnail creation error: %v", err)
				return
			}

			// S3에 이미지 및 썸네일 업로드
			fileName, thumbnailFileName, err := uploadImagesToS3(imgData, thumbnailData, contentType, ext, service.s3svc, service.bucket, service.bucketUrl, uidString)
			if err != nil {
				errorsChan <- fmt.Errorf("error uploading images to S3: %v", err)
				return
			}

			uploadedFiles <- fileName
			uploadedFiles <- thumbnailFileName

			images[i] = model.Image{Uid: dietRequest.Uid, Url: fileName, ThumbnailUrl: thumbnailFileName}
		}(i, imgStr)
	}

	wg.Wait()
	close(errorsChan)
	close(uploadedFiles)

	// 업로드 중 에러 확인 및 처리
	var uploadErrorOccurred bool
	for err := range errorsChan {
		uploadErrorOccurred = true
		fmt.Println(err) // 에러 로깅
	}

	if uploadErrorOccurred {
		tx.Rollback()

		// 이미 업로드된 파일들을 S3에서 삭제
		go func() {
			for file := range uploadedFiles {
				deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl, uidString)
			}
		}()
		return "", fmt.Errorf("error occurred during image upload")
	}

	// Diet 객체에 이미지 정보 추가
	diet.Images = images

	// 데이터베이스 작업
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 새 레코드 생성
		diet.Id = 0
		if err := tx.Create(&diet).Error; err != nil {
			tx.Rollback()
			// 이미 업로드된 파일들을 S3에서 삭제
			go func() {
				for file := range uploadedFiles {
					deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl, uidString)
				}
			}()
			return "", err
		}
	} else {
		// 기존 이미지 레코드 논리삭제
		result := service.db.Model(&model.Image{}).Where("diet_id = ?", diet.Id).Select("level").Updates(map[string]interface{}{"level": 10})
		if result.Error != nil {
			tx.Rollback()
			return "", errors.New("db error")
		}
		// 기존 레코드 업데이트
		if err := tx.Model(&diet).Updates(diet).Error; err != nil {
			tx.Rollback()
			// 이미 업로드된 파일들을 S3에서 삭제
			go func() {
				for file := range uploadedFiles {
					deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl, uidString)
				}
			}()
			return "", err
		}

	}

	// 트랜잭션 커밋
	tx.Commit()
	return "200", nil
}

func (service *dietService) RemoveDiet(ids []int, uid int) (string, error) {
	tx := service.db.Begin()
	result := tx.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.Diet{})

	if result.Error != nil {
		tx.Rollback()
		return "", errors.New("db error")
	}

	result = tx.Model(&model.Image{}).Where("diet_id =IN (?)", ids).Select("level").Updates(map[string]interface{}{"level": 10})

	if result.Error != nil {
		tx.Rollback()
		return "", errors.New("db error2")
	}

	tx.Commit()
	return "200", nil
}

func (service *dietService) GetPresets(id int, page int, startDate, endDate string) ([]dto.DietPresetResponse, error) {
	pageSize := 10
	var dietPresets []model.DietPreset
	offset := page * pageSize

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")
	result := query.Offset(offset).Limit(pageSize).Find(&dietPresets)

	if result.Error != nil {
		return nil, result.Error
	}

	var dietPresetResponses []dto.DietPresetResponse
	if err := util.CopyStruct(dietPresets, &dietPresetResponses); err != nil {
		return nil, err
	}

	return dietPresetResponses, nil
}

func (service *dietService) SavePreset(presetRequest dto.DietPresetRequest) (string, error) {
	var dietPreset model.DietPreset

	result := service.db.First(&model.DietPreset{}, presetRequest.Id)

	if err := util.CopyStruct(presetRequest, &dietPreset); err != nil {
		return "", err
	}
	dietPreset.Uid = presetRequest.Uid

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		dietPreset.Id = 0
		if err := service.db.Create(&dietPreset).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&dietPreset).Updates(dietPreset).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}

func (service *dietService) RemovePreset(ids []int, uid int) (string, error) {
	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.DietPreset{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
