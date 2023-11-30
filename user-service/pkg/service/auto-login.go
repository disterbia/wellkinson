// /user-service/pkg/service/auto-login.go
package service

import (
	"common/model"
	"common/util"
	"errors"

	"gorm.io/gorm"
)

type autoLoginService struct {
	db *gorm.DB
}

func NewAutoLoginService(db *gorm.DB) LoginService {
	return &autoLoginService{db: db}
}

func (al *autoLoginService) Login(email string, user model.User) (string, error) {

	// 데이터베이스에서 사용자 조회
	var u model.User
	if err := al.db.Where(model.User{Email: email}).First(&u).Error; err != nil {
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
