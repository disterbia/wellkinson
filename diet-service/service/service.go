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
	GetDietPresets(id int, page int) ([]dto.DietPresetResponse, error)
	RemovePreset(id int, uid int) (string, error)
}

type dietPresetService struct {
	db *gorm.DB
}

func NewDietPresetService(db *gorm.DB) DietPresetService {
	return &dietPresetService{db: db}
}

func (service *dietPresetService) GetDietPresets(id int, page int) ([]dto.DietPresetResponse, error) {
	pageSize := 10
	var dietPresets []model.DietPreset
	offset := page * pageSize

	result := service.db.Where("uid = ? ", id).Order("id DESC").Offset(offset).Limit(pageSize).Find(&dietPresets)
	if result.Error != nil {
		return nil, result.Error
	}

	var dietPresetResponses []dto.DietPresetResponse
	if err := util.CopyStruct(dietPresets, &dietPresetResponses); err != nil {
		return nil, err
	}

	return dietPresetResponses, nil
}

func (service *dietPresetService) SaveDietPreset(presetRequest dto.DietPresetRequest) (string, error) {
	var dietPreset model.DietPreset

	result := service.db.First(&model.DietPreset{}, presetRequest.Id)

	if err := util.CopyStruct(presetRequest, &dietPreset); err != nil {
		return "", err
	}
	dietPreset.Uid = presetRequest.Uid

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		dietPreset.Id = 0
		if err := service.db.Create(&dietPreset).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&dietPreset).Updates(dietPreset).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}

func (service *dietPresetService) RemovePreset(id int, uid int) (string, error) {

	result := service.db.Where("id=? AND uid= ?", id, uid).Delete(&model.DietPreset{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
