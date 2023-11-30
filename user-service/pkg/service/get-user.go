// /user-service/pkg/service/set-user.go

package service

import (
	"common/model"

	"gorm.io/gorm"
)

type GetUserService interface {
	GetUser(id int) (model.User, error)
}

type getUserService struct {
	db *gorm.DB
}

func NewGetUserService(db *gorm.DB) GetUserService {
	return &getUserService{db: db}
}

func (gu *getUserService) GetUser(id int) (model.User, error) {
	var user model.User
	result := gu.db.First(&user, id)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}
