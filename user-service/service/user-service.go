// /user-service/service/user-service.go

package service

import (
	"common/model"
	"common/util"
	"errors"
	"log"
	"user-service/dto"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserService interface {
	AutoLogin(email string, user model.User) (string, error)                 //자동로그인
	KakaoLogin(idToken string, userReqeust dto.UserRequest) (string, error)  //카카오로그인
	GoogleLogin(idToken string, userRequest dto.UserRequest) (string, error) //구글로그인
	findOrCreateUser(user model.User) (model.User, error)                    //로그인처리
	SetUser(user dto.UserRequest) (string, error)                            //유저업데이트
	GetUser(id int) (dto.UserResponse, error)                                //유저조회
	AdminLogin(email string, password string) (string, error)
}

type userService struct {
	db *gorm.DB
}

type KakaoPublicKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []KakaoPublicKey `json:"keys"`
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}
func (service *userService) AdminLogin(email string, password string) (string, error) {
	log.Println("fff")
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

func (su *userService) SetUser(userRequest dto.UserRequest) (string, error) {

	// 유효성 검사 수행
	if err := util.ValidateDate(userRequest.Birthday); err != nil {
		return "", err
	}
	if err := util.ValidatePhoneNumber(userRequest.PhoneNum); err != nil {
		return "", err
	}

	var user model.User
	if err := util.CopyStruct(userRequest, &user); err != nil {
		return "", err
	}

	result := su.db.Model(&model.User{}).Where("id = ?", userRequest.Id).Updates(user)
	if result.Error != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}

func (gu *userService) GetUser(id int) (dto.UserResponse, error) {
	var user model.User
	result := gu.db.First(&user, id)
	if result.Error != nil {
		return dto.UserResponse{}, result.Error
	}
	var userResponse dto.UserResponse
	log.Println(user.Created)
	if err := util.CopyStruct(user, &userResponse); err != nil {
		return dto.UserResponse{}, err
	}

	return userResponse, nil
}
