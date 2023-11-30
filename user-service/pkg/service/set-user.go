// /user-service/pkg/service/set-user.go

package service

import (
	"common/model"
	"common/util"

	"gorm.io/gorm"
)

type SetUserService interface {
	SetUser(user model.User) (string, error)
}

type setUserService struct {
	db *gorm.DB
}

func NewSetUserService(db *gorm.DB) SetUserService {
	return &setUserService{db: db}
}

func (su *setUserService) SetUser(user model.User) (string, error) {

	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return "", err
	}

	result := su.db.Model(&model.User{}).Where("id = ?", user.Id).Updates(user)
	if result.Error != nil {
		return "", result.Error
	}

	return "200", nil
}
