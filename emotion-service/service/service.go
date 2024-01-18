// /emotion-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"emotion-service/dto"
	"errors"

	"gorm.io/gorm"
)

type EmotionService interface {
	SaveEmotion(emotionRequest dto.EmotionRequest) (string, error)
	GetEmotions(id int, startDate, endDate string) ([]dto.EmotionResponse, error)
	RemoveEmotions(ids []int, uid int) (string, error)
}

type emotionService struct {
	db *gorm.DB
}

func NewEmotionService(db *gorm.DB) EmotionService {
	return &emotionService{db: db}
}

func (service *emotionService) SaveEmotion(emotionRequest dto.EmotionRequest) (string, error) {
	var emotion model.Emotion

	result := service.db.First(&model.Emotion{}, emotionRequest.Id)

	if err := util.CopyStruct(emotionRequest, &emotion); err != nil {
		return "", err
	}
	emotion.Uid = emotionRequest.Uid

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		emotion.Id = 0
		if err := service.db.Create(&emotion).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&emotion).Updates(emotion).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}
func (service *emotionService) GetEmotions(id int, startDate, endDate string) ([]dto.EmotionResponse, error) {
	var emotions []model.Emotion

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")
	result := query.Find(&emotions)

	if result.Error != nil {
		return nil, result.Error
	}

	var emotionResponses []dto.EmotionResponse
	if err := util.CopyStruct(emotions, &emotionResponses); err != nil {
		return nil, err
	}

	return emotionResponses, nil
}

func (service *emotionService) RemoveEmotions(ids []int, uid int) (string, error) {
	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.Emotion{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
