// /user-service/pkg/service/user-service.go

package service

import (
	"common/model"
	"common/util"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserService interface {
	AutoLogin(email string, user model.User) (string, error)            //자동로그인
	KakaoLogin(idToken string, user model.User) (string, error)         //카카오로그인
	GoogleLogin(idToken string, user model.User) (string, error)        //구글로그인
	findOrCreateUser(email string, user model.User) (model.User, error) //로그인처리
	SetUser(user model.User) (string, error)                            //유저업데이트
	GetUser(id int) (model.User, error)                                 //유저조회
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

func (service *userService) KakaoLogin(idToken string, user model.User) (string, error) {
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
		u, err := service.findOrCreateUser(email, user)
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

func (service *userService) GoogleLogin(idToken string, user model.User) (string, error) {
	email, err := validateGoogleIDToken(idToken)
	if err != nil {
		return "", errors.New(err.Error())
	}

	u, err := service.findOrCreateUser(email, user)
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

func (service *userService) findOrCreateUser(email string, user model.User) (model.User, error) {
	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return model.User{}, err
	}
	// 데이터베이스에서 사용자 조회 및 없으면 생성
	err := service.db.Where(model.User{Email: email}).FirstOrCreate(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (su *userService) SetUser(user model.User) (string, error) {

	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return "", err
	}

	result := su.db.Model(&model.User{}).Where("id = ?", user.Id).Updates(user)
	if result.Error != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}

func (gu *userService) GetUser(id int) (model.User, error) {
	var user model.User
	result := gu.db.First(&user, id)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}
