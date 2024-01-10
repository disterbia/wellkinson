// /alarm-service/service/alarm-service.go

package service

import (
	"alarm-service/dto"
	"common/model"
	"common/util"
	"errors"

	"gorm.io/gorm"
)

type AlarmService interface {
	SaveAlarm(alarmRequest dto.AlarmRequest) (string, error)
	RemoveAlarm(id int, uid int) (string, error)
}

type alarmService struct {
	db *gorm.DB
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmService{db: db}
}

func (service *alarmService) SaveAlarm(alarmRequest dto.AlarmRequest) (string, error) {
	// 유효성 검사 수행
	if err := validateAlarm(alarmRequest); err != nil {
		return "", err
	}
	var existingAlarm model.Alarm
	var alarm model.Alarm
	result := service.db.First(&existingAlarm, alarmRequest.Id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		alarmRequest.Id = 0
		if err := util.CopyStruct(alarmRequest, &alarm); err != nil {
			return "", err
		}
		alarm.Uid = alarmRequest.Uid
		if err := service.db.Create(&alarm).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&existingAlarm).Updates(alarmRequest).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}

func (ra *alarmService) RemoveAlarm(id int, uid int) (string, error) {

	result := ra.db.Where("id=? AND uid= ?", id, uid).Delete(&model.Alarm{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
