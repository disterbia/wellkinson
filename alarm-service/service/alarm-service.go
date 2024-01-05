// /alarm-service/service/alarm-service.go

package service

import (
	"common/model"
	"errors"

	"gorm.io/gorm"
)

type AlarmService interface {
	RemoveAlarm(alarm model.Alarm) (string, error)
	SaveAlarm(alarm model.Alarm) (string, error)
}

type alarmService struct {
	db *gorm.DB
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmService{db: db}
}

func (service *alarmService) SaveAlarm(alarm model.Alarm) (string, error) {
	// 유효성 검사 수행
	if err := validateAlarm(alarm); err != nil {
		return "", err
	}
	var existingAlarm model.Alarm
	result := service.db.First(&existingAlarm, alarm.Id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		alarm.Id = 0
		if err := service.db.Create(&alarm).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&existingAlarm).Updates(alarm).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}

func (ra *alarmService) RemoveAlarm(alarm model.Alarm) (string, error) {

	result := ra.db.Where("id=? AND uid= ?", alarm.Id, alarm.Uid).Delete(&alarm)

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
