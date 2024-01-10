// /diet-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"diet-service/dto"
	"errors"

	"gorm.io/gorm"
)

type DietPresetService interface {
	SaveDietPreset(presetRequest dto.DietPresetRequest) (string, error)
}

type dietPresetService struct {
	db *gorm.DB
}

func NewDietPresetService(db *gorm.DB) DietPresetService {
	return &dietPresetService{db: db}
}

func (service *dietPresetService) SaveDietPreset(presetRequest dto.DietPresetRequest) (string, error) {
	var dietPreset model.DietPreset
	if err := util.CopyStruct(presetRequest, &dietPreset); err != nil {
		return "", err
	}

	dietPreset.Uid = presetRequest.Uid
	// DietPreset 저장
	if err := service.db.Save(&dietPreset).Error; err != nil {
		return "", errors.New("db error")
	}

	// DietPreset과 연관된 Foods 저장
	if len(dietPreset.Foods) > 0 {
		err := service.db.Model(&dietPreset).Association("Foods").Replace(dietPreset.Foods)
		if err != nil {
			// 에러 처리
			return "", err
		}
	}

	return "200", nil
}
