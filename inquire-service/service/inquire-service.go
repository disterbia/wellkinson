// /inquire-service/service/inquire-service.go

package service

import (
	"common/model"
	"errors"

	"gorm.io/gorm"
)

type InquireService interface {
	AnswerInquire(answer model.InquireReply) (string, error)
	SendInquire(inquire model.Inquire) (string, error)
	GetMyInquires(id int) ([]model.Inquire, error)
}

type inquireService struct {
	db *gorm.DB
}

func NewInquireService(db *gorm.DB) InquireService {
	return &inquireService{db: db}
}

func (service *inquireService) GetMyInquires(id int) ([]model.Inquire, error) {
	var inquires []model.Inquire
	result := service.db.Where("uid = ?", id).Find(&inquires)
	if result.Error != nil {
		return nil, result.Error
	}
	return inquires, nil
}

func (service *inquireService) SendInquire(inquire model.Inquire) (string, error) {
	result := service.db.Save(&inquire)

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) AnswerInquire(answer model.InquireReply) (string, error) {
	var user model.User
	var inquire model.Inquire

	result := service.db.First(&user, answer.Uid)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	result2 := service.db.First(&inquire, answer.InquireId)

	if result2.Error != nil {
		return "", result2.Error
	}

	if answer.ReplyType { //답변
		if !user.IsAdmin() { // IsAdmin 메서드를 사용하여 확인
			return "", errors.New("unauthorized: user is not an admin")
		}
	} else { // 추가문의
		if user.IsAdmin() { // IsAdmin 메서드를 사용하여 확인
			return "", errors.New("unauthorized: user is admin")
		}
	}

	answer.Id = 0
	result = service.db.Save(answer)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	err := sendEmail(inquire, answer)
	if err != nil {
		return "", err
	}
	return "200", nil
}
