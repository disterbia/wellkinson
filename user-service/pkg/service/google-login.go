// /user-service/pkg/service/google-login.go

package service

import (
	"common/model"
	"common/util"
	"context"
	"errors"
	"log"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type googleLoginService struct {
	db *gorm.DB
}

func NewGoogleLoginService(db *gorm.DB) LoginService {
	return &googleLoginService{db: db}
}

func (gl *googleLoginService) Login(idToken string, user model.User) (string, error) {
	email, err := validateGoogleIDToken(idToken)
	if err != nil {
		return "", errors.New(err.Error())
	}

	u, err := gl.findOrCreateUser(email, user)
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

func (gl *googleLoginService) findOrCreateUser(email string, user model.User) (model.User, error) {
	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return model.User{}, err
	}
	// 데이터베이스에서 사용자 조회 및 없으면 생성
	err := gl.db.Where(model.User{Email: email}).FirstOrCreate(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// Google ID 토큰을 검증하고 이메일을 반환
func validateGoogleIDToken(idToken string) (string, error) {
	log.Print("idToken: ", idToken)
	// idtoken 패키지를 사용하여 토큰 검증
	payload, err := idtoken.Validate(context.Background(), idToken, "390432007084-1hqslpiclba2hucb6hl41acecv1qekbt.apps.googleusercontent.com")
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return "", err
	}

	// 이메일 추출
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", errors.New("email claim not found in token")
	}

	return email, nil
}
