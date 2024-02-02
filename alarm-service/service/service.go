// /alarm-service/service/alarm-service.go

package service

import (
	"alarm-service/dto"
	"common/model"
	"common/util"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type AlarmService interface {
	SaveAlarm(alarmRequest dto.AlarmRequest) (string, error)
	RemoveAlarm(ids []int, uid int) (string, error)
	GetAlarms(id int, page int) ([]dto.AlarmResponse, error)
}

type alarmService struct {
	db *gorm.DB
}

func NewAlarmService(db *gorm.DB) AlarmService {
	return &alarmService{db: db}
}

func (service *alarmService) GetAlarms(id int, page int) ([]dto.AlarmResponse, error) {
	pageSize := 10
	var alarms []model.Alarm
	offset := page * pageSize

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
	// JSON 배열을 Go 슬라이스로 변환
	var weekdaySlice []int
	err := json.Unmarshal(alarm.Week, &weekdaySlice)
	if err != nil {
		return "", err
	}

	seen := make(map[int]bool)
	unique := []int{}

	for _, v := range weekdaySlice {
		// 숫자가 0과 6 사이인지 확인
		if v >= 0 && v <= 6 {
			// 중복되지 않은 경우, 결과 슬라이스에 추가
			if !seen[v] {
				seen[v] = true
				unique = append(unique, v)
			}
		}
	}

	newWeekdays, err := json.Marshal(unique)
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

func (ra *alarmService) RemoveAlarm(ids []int, uid int) (string, error) {

	result := ra.db.Where("id IN ? AND uid= ?", ids, uid).Delete(&model.Alarm{})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
