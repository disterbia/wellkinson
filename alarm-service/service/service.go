// /alarm-service/service/alarm-service.go

package service

import (
	"alarm-service/common/model"
	"alarm-service/common/util"
	"alarm-service/dto"
	"errors"
	"log"

	"gorm.io/gorm"
)

type AlarmService interface {
	SaveAlarm(alarmRequest dto.AlarmRequest) (string, error)
	RemoveAlarm(ids []uint, uid uint) (string, error)
	GetAlarms(id uint, page uint) ([]dto.AlarmResponse, error)
	GetNotifications(uid uint) ([]dto.NotificationResponse, error)
	ReadAll(uid uint) (string, error)
	RemoveNotifications(ids []uint, uid uint) (string, error)
}

type alarmService struct {
	db *gorm.DB
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmService{db: db}
}
func (service *alarmService) GetNotifications(uid uint) ([]dto.NotificationResponse, error) {
	var notifications []model.Notification
	result := service.db.Where("uid = ? ", uid).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	var notificationResponses []dto.NotificationResponse
	if err := util.CopyStruct(notifications, &notificationResponses); err != nil {
		return nil, err
	}

	return notificationResponses, nil
}

func (service *alarmService) GetAlarms(id uint, page uint) ([]dto.AlarmResponse, error) {
	pageSize := 10
	var alarms []model.Alarm
	offset := int(page) * pageSize

	log.Println("fdffsd")
	result := service.db.Where("uid = ? ", id).Order("id DESC").Offset(offset).Limit(pageSize).Find(&alarms)
	if result.Error != nil {
		return nil, result.Error
	}

	var alarmResponses []dto.AlarmResponse
	if err := util.CopyStruct(alarms, &alarmResponses); err != nil {
		return nil, err
	}

	return alarmResponses, nil
}
func (service *alarmService) SaveAlarm(alarmRequest dto.AlarmRequest) (string, error) {
	// 유효성 검사 수행
	if err := validateAlarm(alarmRequest); err != nil {
		return "", err
	}

	var alarm model.Alarm
	result := service.db.Where("id=? AND uid=?", alarmRequest.Id, alarmRequest.Uid).First(&model.Alarm{})

	if err := util.CopyStruct(alarmRequest, &alarm); err != nil {
		return "", err
	}
	alarm.Uid = alarmRequest.Uid //  json: "-" 이라서
	newWeekdays, _, err := validateWeek(alarm.Week)
	if err != nil {
		return "", err
	}
	alarm.Week = newWeekdays

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		alarm.Id = 0
		if err := service.db.Create(&alarm).Error; err != nil {
			return "", errors.New("db error")
		}
	} else if result.Error != nil {
		return "", errors.New("db error2")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&alarm).Updates(alarm).Error; err != nil {
			return "", errors.New("db error3")
		}
	}

	return "200", nil
}

func (service *alarmService) RemoveAlarm(ids []uint, uid uint) (string, error) {

	result := service.db.Where("id IN ? AND uid= ?", ids, uid).Delete(&model.Alarm{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *alarmService) ReadAll(uid uint) (string, error) {
	var noti model.Notification
	result := service.db.Model(&noti).Where("uid = ?", uid).Select("is_read").Updates(map[string]interface{}{"is_read": true})
	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *alarmService) RemoveNotifications(ids []uint, uid uint) (string, error) {

	result := service.db.Where("id IN ? AND uid= ?", ids, uid).Delete(&model.Notification{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
