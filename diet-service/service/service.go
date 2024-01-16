// /diet-service/service/service.go
package service

import (
	"bytes"
	"common/model"
	"common/util"
	"diet-service/dto"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"

	"gorm.io/gorm"
)

type DietService interface {
	SavePreset(presetRequest dto.DietPresetRequest) (string, error)
	GetPresets(id int, page int, startDate, endDate string) ([]dto.DietPresetResponse, error)
	RemovePreset(id int, uid int) (string, error)
	SaveDiet(diet dto.DietRequest) (string, error)
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
		return "", result.Error
	}

	// DietRequest의 복사본 생성
	dietRequestCopy := dietRequest
	dietRequestCopy.Images = []string{} // Images 필드 초기화

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
			fileName, thumbnailFileName, err := uploadImagesToS3(imgData, thumbnailData, contentType, ext, service.s3svc, service.bucket, service.bucketUrl)
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
				deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl)
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
					deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl)
				}
			}()
			return "", err
		}
	} else {
		// 기존 레코드 업데이트
		if err := tx.Model(&diet).Updates(diet).Error; err != nil {
			tx.Rollback()
			// 이미 업로드된 파일들을 S3에서 삭제
			go func() {
				for file := range uploadedFiles {
					deleteFromS3(file, service.s3svc, service.bucket, service.bucketUrl)
				}
			}()
			return "", err
		}

	}

	// 트랜잭션 커밋
	tx.Commit()
	return "200", nil
}

func deleteFromS3(fileKey string, s3Client *s3.S3, bucket string, bucketUrl string) error {

	// URL에서 객체 키 추출
	key := extractKeyFromUrl(fileKey, bucket, bucketUrl)
	log.Println("key", fileKey)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket), // 실제 S3 버킷 이름으로 대체
		Key:    aws.String(key),
	})

	// 에러 발생 시 처리 로직
	if err != nil {
		fmt.Printf("Failed to delete object from S3: %s, error: %v\n", fileKey, err)
	}

	return err
}

// URL에서 S3 객체 키를 추출하는 함수
func extractKeyFromUrl(url, bucket string, bucketUrl string) string {
	prefix := fmt.Sprintf("https://%s.%s/", bucket, bucketUrl)
	return strings.TrimPrefix(url, prefix)
}
func uploadImagesToS3(imgData []byte, thumbnailData []byte, contentType string, ext string, s3Client *s3.S3, bucket string, bucketUrl string) (string, string, error) {
	// 이미지 파일 이름과 썸네일 파일 이름 생성
	imgFileName := "images/diet/" + uuid.New().String() + ext
	thumbnailFileName := "images/diet/thumbnail/" + uuid.New().String() + ext

	// S3에 이미지 업로드
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imgFileName),
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", err
	}

	// S3에 썸네일 업로드
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(thumbnailFileName),
		Body:        bytes.NewReader(thumbnailData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", err
	}

	// 업로드된 이미지와 썸네일의 URL 생성 및 반환
	imgURL := "https://" + bucket + "." + bucketUrl + "/" + imgFileName
	thumbnailURL := "https://" + bucket + "." + bucketUrl + "/" + thumbnailFileName

	return imgURL, thumbnailURL, nil
}
func reduceImageSize(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	log.Println("image size: ", len(data))
	// 원본 이미지의 크기를 절반씩 줄이면서 10MB 이하로 만듦
	for len(data) > 10*1024*1024 {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		err := jpeg.Encode(&buf, resizedImg, nil)
		if err != nil {
			return nil, err
		}

		data = buf.Bytes()
		img = resizedImg
	}

	return data, nil
}

func createThumbnail(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 썸네일의 크기를 절반씩 줄이면서 1MB 이하로 만듦
	for {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		thumbnail := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, thumbnail, nil)
		if err != nil {
			return nil, err
		}

		thumbnailData := buf.Bytes()
		log.Println("thumbnailData size: ", len(thumbnailData))
		if len(thumbnailData) < 1024*1024 {
			return thumbnailData, nil
		}

		img = thumbnail
	}
}

func getImageFormat(imgData []byte) (contentType, extension string, err error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return "", "", err
	}

	contentType = "image/" + format
	extension = "." + format
	if format == "jpeg" {
		extension = ".jpg"
	}

	return contentType, extension, nil
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

func (service *dietService) RemovePreset(id int, uid int) (string, error) {

	result := service.db.Where("id=? AND uid= ?", id, uid).Delete(&model.DietPreset{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
