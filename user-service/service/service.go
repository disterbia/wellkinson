// /user-service/service/service.go

package service

import (
	"common/model"
	"common/util"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"
	"user-service/dto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserService interface {
	AutoLogin(email string, user model.User) (string, error)                 //자동로그인
	KakaoLogin(idToken string, userReqeust dto.UserRequest) (string, error)  //카카오로그인
	GoogleLogin(idToken string, userRequest dto.UserRequest) (string, error) //구글로그인
	findOrCreateUser(user model.User) (model.User, error)                    //로그인처리
	SetUser(user dto.UserRequest) (string, error)                            //유저업데이트
	GetUser(id uint) (dto.UserResponse, error)                               //유저조회
	AdminLogin(email string, password string) (string, error)
	GetMainServices() ([]dto.MainServiceResponse, error)
}

type userService struct {
	db        *gorm.DB
	s3svc     *s3.S3
	bucket    string
	bucketUrl string
}

type KakaoPublicKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []KakaoPublicKey `json:"keys"`
}

func NewUserService(db *gorm.DB, s3svc *s3.S3, bucket string, bucketUrl string) UserService {
	return &userService{db: db, s3svc: s3svc, bucket: bucket, bucketUrl: bucketUrl}
}

func (service *userService) AdminLogin(email string, password string) (string, error) {
	var u model.User
	if err := service.db.Where(model.User{Email: email, PhoneNum: password}).First(&u).Error; err != nil {
		return "", err
	}

	if !u.IsAdmin {
		return "", errors.New("not admin")
	}

	// 새로운 JWT 토큰 생성
	tokenString, err := util.GenerateJWT(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (service *userService) AutoLogin(email string, user model.User) (string, error) {

	// 데이터베이스에서 사용자 조회
	var u model.User
	if err := service.db.Where(model.User{Email: email}).First(&u).Error; err != nil {
		return "", err
	}
	if !u.UseAutoLogin {
		return "", errors.New("not use auto login")
	}
	// 새로운 JWT 토큰 생성
	tokenString, err := util.GenerateJWT(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *userService) KakaoLogin(idToken string, userRequest dto.UserRequest) (string, error) {
	jwks, err := getKakaoPublicKeys()
	if err != nil {
		return "", err
	}

	parsedToken, err := verifyKakaoTokenSignature(idToken, jwks)
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("email not found in token claims")
		}

		var user model.User
		if err := util.CopyStruct(userRequest, &user); err != nil {
			return "", err
		}

		user.Email = email
		u, err := service.findOrCreateUser(user)
		if err != nil {
			return "", errors.New(err.Error())
		}

		// JWT 토큰 생성
		tokenString, err := util.GenerateJWT(u)
		if err != nil {
			return "", errors.New(err.Error())
		}
		return tokenString, nil
	}
	return "", errors.New("invalid token")

}

func (service *userService) GoogleLogin(idToken string, userRequest dto.UserRequest) (string, error) {
	email, err := validateGoogleIDToken(idToken)
	if err != nil {
		return "", errors.New(err.Error())
	}

	var user model.User

	if err := util.CopyStruct(userRequest, &user); err != nil {
		return "", err
	}

	user.Email = email
	u, err := service.findOrCreateUser(user)
	if err != nil {
		return "", errors.New(err.Error())
	}

	// JWT 토큰 생성
	tokenString, err := util.GenerateJWT(u)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return tokenString, nil
}

func (service *userService) findOrCreateUser(user model.User) (model.User, error) {
	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return model.User{}, err
	}
	if err := util.ValidatePhoneNumber(user.PhoneNum); err != nil {
		return model.User{}, err
	}

	// 데이터베이스에서 사용자 조회 및 없으면 생성
	err := service.db.Where(model.User{Email: user.Email}).FirstOrCreate(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (service *userService) SetUser(userRequest dto.UserRequest) (string, error) {

	// 유효성 검사 수행
	if err := util.ValidateDate(userRequest.Birthday); err != nil {
		return "", err
	}
	if err := util.ValidatePhoneNumber(userRequest.PhoneNum); err != nil {
		return "", err
	}
	var fileName, thumbnailFileName string
	var image model.Image

	var tempUser dto.TempUser

	if err := util.CopyStruct(userRequest, &tempUser); err != nil {
		return "", err
	}
	var user model.User
	if err := util.CopyStruct(tempUser, &user); err != nil {
		return "", err
	}

	user.Id = userRequest.Id

	if userRequest.ProfileImage != "" {
		imgData, err := base64.StdEncoding.DecodeString(userRequest.ProfileImage)
		if err != nil {
			return "", err
		}

		contentType, ext, err := getImageFormat(imgData)
		if err != nil {
			return "", err
		}

		// 이미지 크기 조정 (10MB 제한)
		if len(imgData) > 10*1024*1024 {
			imgData, err = reduceImageSize(imgData)
			if err != nil {
				return "", err
			}
		}

		// 썸네일 이미지 생성
		thumbnailData, err := createThumbnail(imgData)
		if err != nil {
			return "", err
		}

		// S3에 이미지 및 썸네일 업로드
		fileName, thumbnailFileName, err = uploadImagesToS3(imgData, thumbnailData, contentType, ext, service.s3svc, service.bucket, service.bucketUrl, strconv.FormatUint(uint64(user.Id), 10))
		if err != nil {

			return "", err
		}
		image = model.Image{
			Uid:          user.Id,
			Url:          fileName,
			ThumbnailUrl: thumbnailFileName,
			ParentId:     user.Id,
			Type:         uint(util.UserProfileImageType),
		}
	}

	// 트랜잭션 시작
	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	//유저 정보 업데이트
	result := service.db.Model(&model.User{}).Where("id = ?", userRequest.Id).Updates(user)
	if result.Error != nil {
		tx.Rollback()

		// 이미 업로드된 파일들을 S3에서 삭제
		if userRequest.ProfileImage != "" {
			go func() {
				deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
				deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
			}()
		}

		return "", errors.New("db error")
	}

	if userRequest.ProfileImage != "" {
		// 기존 이미지 레코드 논리삭제
		result = service.db.Model(&model.Image{}).Where("parent_id = ? AND type =?", user.Id, util.UserProfileImageType).Select("level").Updates(map[string]interface{}{"level": 10})
		if result.Error != nil {
			tx.Rollback()
			if userRequest.ProfileImage != "" {
				go func() {
					deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
				}()
			}
			return "", errors.New("db error4")
		}
		// 이미지 레코드 재 생성

		if err := tx.Create(&image).Error; err != nil {
			tx.Rollback()
			if userRequest.ProfileImage != "" {
				go func() {
					deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
				}()
			}
			return "", errors.New("db error5")
		}
	}

	//서비스항목 조회
	mainServices := make([]model.MainService, 0)
	result = service.db.Where("id IN ? ", userRequest.UseServices).Find(&mainServices)
	if result.Error != nil {
		tx.Rollback()

		if userRequest.ProfileImage != "" {
			go func() {
				deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
				deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
			}()
		}
		return "", errors.New("db error3")
	}

	//유저별 사용 서비스 삭제 후
	result = service.db.Where("uid = ?", userRequest.Id).Delete(&model.UseService{})
	if result.Error != nil {
		tx.Rollback()

		if userRequest.ProfileImage != "" {
			go func() {
				deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
				deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
			}()
		}
		return "", errors.New("db error6")
	}

	//유저별 사용 서비스 생성
	if len(mainServices) > 0 {
		var useService []model.UseService
		for _, v := range mainServices {
			useService = append(useService, model.UseService{
				Uid:       user.Id,
				ServiceId: v.Id,
				Title:     v.Title,
			})
		}

		result = service.db.Create(&useService)
		if result.Error != nil {
			tx.Rollback()
			if userRequest.ProfileImage != "" {
				go func() {
					deleteFromS3(fileName, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(thumbnailFileName, service.s3svc, service.bucket, service.bucketUrl)
				}()
			}
			return "", errors.New("db error7")
		}

	}
	tx.Commit()
	return "200", nil
}

func (service *userService) GetUser(id uint) (dto.UserResponse, error) {
	var user model.User
	result := service.db.Joins("ProfileImage", "level != 10 AND type = ?", util.UserProfileImageType).First(&user, id)
	if result.Error != nil {
		return dto.UserResponse{}, result.Error
	}

	var useServices []model.UseService
	result = service.db.Where("uid = ?", id).Find(&useServices, id)
	if result.Error != nil {
		return dto.UserResponse{}, result.Error
	}

	var mainServices []dto.MainServiceResponse
	for _, v := range useServices {
		mainServices = append(mainServices, dto.MainServiceResponse{
			Id:    v.ServiceId,
			Title: v.Title,
		})
	}

	var userResponse dto.UserResponse
	if err := util.CopyStruct(user, &userResponse); err != nil {
		return dto.UserResponse{}, err
	}
	userResponse.UseServices = mainServices

	// 여기서 프로필 이미지가 비었는지 체크 해야하는거 아니냐
	if userResponse.ProfileImage.Url != "" {
		urlkey := extractKeyFromUrl(userResponse.ProfileImage.Url, service.bucket, service.bucketUrl)
		thumbnailUrlkey := extractKeyFromUrl(userResponse.ProfileImage.ThumbnailUrl, service.bucket, service.bucketUrl)
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
			return dto.UserResponse{}, err
		}
		thumbnailUrlStr, err := thumbnailUrl.Presign(1 * time.Second) // URL은 1초 동안 유효 CachedNetworkImage 에서 캐싱해서 쓰면됨
		if err != nil {
			return dto.UserResponse{}, err
		}
		userResponse.ProfileImage.Url = urlStr // 사전 서명된 URL로 업데이트
		userResponse.ProfileImage.ThumbnailUrl = thumbnailUrlStr
	}

	return userResponse, nil
}

func (service *userService) GetMainServices() ([]dto.MainServiceResponse, error) {
	var services []model.MainService
	result := service.db.Where("level != 10").Find(&services)
	if result.Error != nil {
		return nil, result.Error
	}
	var serviceResposnes []dto.MainServiceResponse
	if err := util.CopyStruct(services, &serviceResposnes); err != nil {
		return nil, err
	}

	return serviceResposnes, nil
}
