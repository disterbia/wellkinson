package service

import (
	"common/model"

	"gorm.io/gorm"
)

type SendService interface {
	SendInquire(inquire model.Inquire) (string, error)
}

type sendService struct {
	db *gorm.DB
}
