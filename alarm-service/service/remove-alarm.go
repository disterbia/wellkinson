// /alarm-service/service/remove-alarm.go

package service

import (
	"common/model"

	"gorm.io/gorm"
)

type RemoveAlarmService interface {
	RemoveAlarm(alarm model.Alarm) (string, error)
}

type removeAlarmService struct {
	db *gorm.DB
}

func NewRemoveAlarmService(db *gorm.DB) RemoveAlarmService {
	return &removeAlarmService{db: db}
}

func (ra *removeAlarmService) RemoveAlarm(alarm model.Alarm) (string, error) {

	result := ra.db.Where("id=? AND uid= ?", alarm.Id, alarm.Uid).Delete(&alarm)

	if result.Error != nil {
		return "", result.Error
	}
	return "200", nil
}
