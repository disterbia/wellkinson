// /user-service/service/service.go

package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"user-service/common/model"
	"user-service/common/util"
	"user-service/dto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserService interface {
	AutoLogin(autoLoginRequest dto.AutoLoginRequest) (string, error) //자동로그인
	SnsLogin(idToken string, userRequest dto.UserRequest) (string, error)
	SetUser(user dto.UserRequest) (string, error) //유저업데이트
	GetUser(id uint) (dto.UserResponse, error)    //유저조회
	AdminLogin(email string, password string) (string, error)
	GetMainServices() ([]dto.MainServiceResponse, error)
	SendAuthCode(number string) (string, error)
	VerifyAuthCode(number, code string) (string, error)
	RemoveUser(id uint) (string, error)
	LinkEmail(uid uint, idToken string) (string, error) // 0:카카오 1:구글 2:애플
	GetVersion() (dto.AppVersionResponse, error)
	RemoveProfile(uid uint) (string, error)
}

type userService struct {
	db        *gorm.DB
	s3svc     *s3.S3
	bucket    string
	bucketUrl string
}

type PublicKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []PublicKey `json:"keys"`
}

func (service *userService) SendAuthCode(number string) (string, error) {
	//존재하는 번호인지 체크
	result := service.db.Debug().Where("phone_number=?", number).Find(&model.VerifiedNumbers{})
	if result.Error != nil {
		return "", errors.New("db error")

	} else if result.RowsAffected > 0 {
		// 레코드가 존재할 때
		return "", errors.New("-1")
	}

	err := util.ValidatePhoneNumber(number)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10)) // 0부터 9까지의 숫자를 무작위로 선택
	}

	apiURL := "https://kakaoapi.aligo.in/akv10/alimtalk/send/"
	data := url.Values{}
	data.Set("apikey", os.Getenv("API_KEY"))
	data.Set("userid", os.Getenv("USER_ID"))
	data.Set("token", os.Getenv("TOKEN"))
	data.Set("senderkey", os.Getenv("SENDER_KEY"))
	data.Set("tpl_code", os.Getenv("TPL_CODE"))
	data.Set("sender", os.Getenv("SENDER"))
	data.Set("subject_1", os.Getenv("SUBJECT_1"))

	data.Set("receiver_1", number)
	data.Set("message_1", "인증번호는 ["+sb.String()+"]"+" 입니다.")

	// HTTP POST 요청 실행
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		fmt.Printf("HTTP Request Failed: %s\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(fmt.Errorf("server returned non-200 status: %d, body: %s", resp.StatusCode, string(body)))

	if err := service.db.Create(&model.AuthCode{PhoneNumber: number, Code: sb.String()}).Error; err != nil {
		return "", err
	}
	return "200", nil
}

func (service *userService) VerifyAuthCode(number, code string) (string, error) {
	var authCode model.AuthCode

	if err := service.db.Where("phone_number = ? ", number).Last(&authCode).Error; err != nil {
		return "", errors.New("db error")
	}
	if authCode.Code != code {
		return "", errors.New("-1")
	}
	if err := service.db.Create(&model.VerifiedNumbers{PhoneNumber: authCode.PhoneNumber}).Error; err != nil {
		return "", errors.New("db error2")
	}

	return "200", nil
}

func NewUserService(db *gorm.DB, s3svc *s3.S3, bucket string, bucketUrl string) UserService {
	return &userService{db: db, s3svc: s3svc, bucket: bucket, bucketUrl: bucketUrl}
}

func (service *userService) AdminLogin(email string, password string) (string, error) {
	var u model.User
	if err := service.db.Debug().Where("email=? AND phone_num=?", email, password).First(&u).Error; err != nil {
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

func (service *userService) AutoLogin(autoLoginRequest dto.AutoLoginRequest) (string, error) {
	if autoLoginRequest.FcmToken == "" || autoLoginRequest.DeviceId == "" {
		return "", errors.New("check fcm_token,device_id")
	}
	// 데이터베이스에서 사용자 조회
	var u model.User
	if err := service.db.Where("email = ?", autoLoginRequest.Email).First(&u).Error; err != nil {
		return "", errors.New("db error")
	}
	if !u.UseAutoLogin {
		return "", errors.New("not use auto login")
	}
	// 새로운 JWT 토큰 생성
	tokenString, err := util.GenerateJWT(u)
	if err != nil {
		return "", err
	}

	if err := service.db.Model(&u).Updates(model.User{FCMToken: autoLoginRequest.FcmToken, DeviceID: autoLoginRequest.DeviceId}).Error; err != nil {
		return "", errors.New("db error2")
	}
	return tokenString, nil
}

func (service *userService) SnsLogin(idToken string, userRequest dto.UserRequest) (string, error) {
	iss := util.DecodeJwt(idToken)

	var user model.User
	var err error

	if strings.Contains(iss, "kakao") { // 카카오
		if user, err = KakaoLogin(idToken, userRequest); err != nil {
			return "", err
		}
	} else if strings.Contains(iss, "google") { // 구글
		if user, err = GoogleLogin(idToken, userRequest); err != nil {
			return "", err
		}
	} else if strings.Contains(iss, "apple") { // 애플
		if user, err = AppleLogin(idToken, userRequest); err != nil {
			return "", err
		}

	} //
	u, err := findOrCreateUser(user, service)
	if err != nil {
		return "", err
	}

	// JWT 토큰 생성
	tokenString, err := util.GenerateJWT(u)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func AppleLogin(idToken string, userRequest dto.UserRequest) (model.User, error) {
	if userRequest.FCMToken == "" || userRequest.DeviceID == "" {
		return model.User{}, errors.New("check fcm_token,device_id")
	}
	jwks, err := getApplePublicKeys()
	if err != nil {
		return model.User{}, err
	}

	parsedToken, err := verifyAppleIDToken(idToken, jwks)
	if err != nil {
		return model.User{}, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return model.User{}, errors.New("email not found in token claims")
		}

		var user model.User

		var temp dto.TempUser
		if err := util.CopyStruct(userRequest, &temp); err != nil {
			return model.User{}, err
		}
		if err := util.CopyStruct(temp, &user); err != nil {
			return model.User{}, err
		}

		user.Email = email
		return user, nil

	}
	return model.User{}, errors.New("invalid token")

}
func KakaoLogin(idToken string, userRequest dto.UserRequest) (model.User, error) {
	if userRequest.FCMToken == "" || userRequest.DeviceID == "" {
		return model.User{}, errors.New("check fcm_token,device_id")
	}
	jwks, err := getKakaoPublicKeys()
	if err != nil {
		return model.User{}, err
	}

	parsedToken, err := verifyKakaoTokenSignature(idToken, jwks)
	if err != nil {
		return model.User{}, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return model.User{}, errors.New("email not found in token claims")
		}

		var user model.User

		var temp dto.TempUser
		if err := util.CopyStruct(userRequest, &temp); err != nil {
			return model.User{}, err
		}
		if err := util.CopyStruct(temp, &user); err != nil {
			return model.User{}, err
		}

		user.Email = email
		return user, nil
	}
	return model.User{}, errors.New("invalid token")

}

func GoogleLogin(idToken string, userRequest dto.UserRequest) (model.User, error) {
	if userRequest.FCMToken == "" || userRequest.DeviceID == "" {
		return model.User{}, errors.New("check fcm_token,device_id")
	}
	email, err := validateGoogleIDToken(idToken)
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	var temp dto.TempUser
	if err := util.CopyStruct(userRequest, &temp); err != nil {
		return model.User{}, err
	}
	if err := util.CopyStruct(temp, &user); err != nil {
		return model.User{}, err
	}

	user.Email = email
	return user, nil
}

func findOrCreateUser(user model.User, service *userService) (model.User, error) {

	fcmToken := user.FCMToken
	deviceId := user.DeviceID

	result := service.db.Where("email = ? ", user.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 연동 로그인
		// 연동 이메일 목록에 없다면 유저생성 (같은번호 있으면 db에서 생성안됨. 같은번호 없으면 회원가입과 같음)
		var linkedEmail model.LinkedEmail
		result := service.db.Where("email = ?", user.Email).First(&linkedEmail)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := service.db.Where("phone_number = ?", user.PhoneNum).First(&model.VerifiedNumbers{}).Error; err != nil {

				return model.User{}, errors.New("-1") // 인증해야함
			}
			// 유효성 검사 수행
			if err := util.ValidateDate(user.Birthday); err != nil {
				return model.User{}, errors.New("-3")
			}
			if err := util.ValidatePhoneNumber(user.PhoneNum); err != nil {
				return model.User{}, errors.New("-3")
			}
			if user.UserType > 2 {
				return model.User{}, errors.New("-3")
			}
			err := service.db.Create(&user).Error
			if err != nil {
				return model.User{}, errors.New("-2") //이미 가입된 번호
			}
		} else if result.Error != nil {
			return model.User{}, errors.New("db error2")
			// 있다면 해당 이메일의 uid로 조회
		} else {
			if err := service.db.Where("id = ?", linkedEmail.Uid).First(&user).Error; err != nil {
				return model.User{}, errors.New("db error3")
			}
		}

	}

	if err := service.db.Model(&user).Updates(model.User{FCMToken: fcmToken, DeviceID: deviceId}).Error; err != nil {
		return model.User{}, errors.New("db error4")
	}

	return user, nil
}

//// 연동 내역 조회 해야함 ////

func (service *userService) LinkEmail(uid uint, idToken string) (string, error) {
	iss := util.DecodeJwt(idToken)

	if strings.Contains(iss, "kakao") { //카카오
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
			if err := saveLinkedEmail(uid, email, service, uint(util.KakaoSnsType)); err != nil {
				return "", err
			}
		}

	} else if strings.Contains(iss, "google") { // 구글
		email, err := validateGoogleIDToken(idToken)
		if err != nil {
			return "", err
		}
		if err := saveLinkedEmail(uid, email, service, uint(util.GoogleSnsType)); err != nil {
			return "", err
		}

	} else if strings.Contains(iss, "apple") { // 애플
		jwks, err := getApplePublicKeys()
		if err != nil {
			return "", err
		}

		parsedToken, err := verifyAppleIDToken(idToken, jwks)
		if err != nil {
			return "", err
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			email, ok := claims["email"].(string)
			if !ok {
				return "", errors.New("email not found in token claims")
			}
			if err := saveLinkedEmail(uid, email, service, uint(util.AppleSnsType)); err != nil {
				return "", err
			}
		}
	} else {
		return "", errors.New("invalid snsType")
	}
	return "200", nil
}

func saveLinkedEmail(uid uint, email string, service *userService, snsType uint) error {
	var user model.User
	if err := service.db.Where("email = ? ", email).First(&user).Error; err != nil {
		return errors.New("db error")
	}
	if user.Email == email {
		return errors.New("wrong request")
	}

	linkedEmail := model.LinkedEmail{Email: email, Uid: uid, SnsType: snsType}

	result := service.db.Where(linkedEmail).First(&model.LinkedEmail{})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		err := service.db.Create(&linkedEmail).Error
		if err != nil {
			return errors.New("db error2")
		}
	} else if result.Error != nil {
		return errors.New("db error3")
	} else {
		// 레코드가 존재하면 삭제
		if err := service.db.Where(linkedEmail).Delete(&model.LinkedEmail{}).Error; err != nil {
			return errors.New("db error4")
		}
	}

	return nil
}

func (service *userService) SetUser(userRequest dto.UserRequest) (string, error) {

	// 유효성 검사 수행
	if userRequest.Birthday != "" {
		if err := util.ValidateDate(userRequest.Birthday); err != nil {
			return "", err
		}
	}
	if userRequest.PhoneNum != "" {
		if err := util.ValidatePhoneNumber(userRequest.PhoneNum); err != nil {
			return "", err
		}
	}

	if err := service.db.Where("phone_number = ?", userRequest.PhoneNum).First(&model.VerifiedNumbers{}).Error; err != nil {
		return "", errors.New("-1")
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

	updateFields := make(map[string]interface{})

	userRequestValue := reflect.ValueOf(userRequest)
	userRequestType := userRequestValue.Type()

	for i := 0; i < userRequestValue.NumField(); i++ {
		field := userRequestValue.Field(i)
		fieldName := userRequestType.Field(i).Tag.Get("json")

		if fieldName == "-" || fieldName == "user_services" || fieldName == "profile_image" {
			continue
		}
		if !field.IsZero() {
			updateFields[fieldName] = field.Interface()
		}
	}

	//유저 정보 업데이트
	result := service.db.Model(&model.User{}).Where("id = ?", userRequest.Id).Updates(updateFields)
	if result.Error != nil {
		log.Println(result.Error.Error())
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
			log.Println(result.Error.Error())
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
			log.Println(err)
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
	result = service.db.Where("id IN ? ", userRequest.UserServices).Find(&mainServices)
	if result.Error != nil {
		log.Println(result.Error.Error())
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
	result = service.db.Where("uid = ?", userRequest.Id).Delete(&model.UserService{})
	if result.Error != nil {
		log.Println(result.Error.Error())
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
		var useService []model.UserService
		for _, v := range mainServices {
			useService = append(useService, model.UserService{
				Uid:       user.Id,
				ServiceId: v.Id,
				Title:     v.Title,
			})
		}

		result = service.db.Create(&useService)
		if result.Error != nil {
			log.Println(result.Error.Error())
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
	result := service.db.Debug().Preload("ProfileImage", "level != ? AND type = ?", 10, util.UserProfileImageType).
		Preload("LinkedEmails").First(&user, id)
	if result.Error != nil {
		return dto.UserResponse{}, errors.New("db error")
	}

	var useServices []model.UserService
	result = service.db.Debug().Where("uid = ?", id).Find(&useServices)
	if result.Error != nil {
		return dto.UserResponse{}, errors.New("db error2")
	}

	mainServices := make([]dto.MainServiceResponse, 0)
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

	userResponse.UserServices = mainServices

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
		urlStr, err := url.Presign(5 * time.Second) // URL은 5초 동안 유효
		if err != nil {
			return dto.UserResponse{}, err
		}
		thumbnailUrlStr, err := thumbnailUrl.Presign(5 * time.Second) // URL은 5초 동안 유효 CachedNetworkImage 에서 캐싱해서 쓰면됨
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
		return nil, errors.New("db error")
	}
	var serviceResposnes []dto.MainServiceResponse
	if err := util.CopyStruct(services, &serviceResposnes); err != nil {
		return nil, err
	}

	return serviceResposnes, nil
}

func (service *userService) RemoveUser(id uint) (string, error) {
	if err := service.db.Delete(&model.User{Id: id}).Error; err != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *userService) GetVersion() (dto.AppVersionResponse, error) {
	var version model.AppVersion
	result := service.db.Last(&version)
	if result.Error != nil {
		return dto.AppVersionResponse{}, errors.New("db error")
	}
	var versionResponse dto.AppVersionResponse
	if err := util.CopyStruct(version, &versionResponse); err != nil {
		return dto.AppVersionResponse{}, err
	}
	return versionResponse, nil
}

func (service *userService) RemoveProfile(uid uint) (string, error) {

	// 기존 이미지 레코드 논리삭제
	result := service.db.Model(&model.Image{}).Where("parent_id = ? AND type =?", uid, util.UserProfileImageType).Select("level").Updates(map[string]interface{}{"level": 10})
	if result.Error != nil {
		return "", errors.New("db error2")
	}
	return "200", nil
}
