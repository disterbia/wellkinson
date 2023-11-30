// /alarm-service/service/save-alarm.go

package service

import (
	"common/model"
	"common/util"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type SaveAlarmService interface {
	SaveAlarm(alarm model.Alarm) (string, error)
}

type saveAlarmService struct {
	db *gorm.DB
}

func NewSaveAlarmService(db *gorm.DB) SaveAlarmService {
	return &saveAlarmService{db: db}
}

func validateWeek(week string) error {
	var weekPattern = regexp.MustCompile(`^[0-6](,[0-6])*$`)

	if !weekPattern.MatchString(week) {
		return errors.New("week must be a comma-separated list of numbers between 0 and 6")
	}

	weekNumbers := strings.Split(week, ",")
	weekMap := make(map[int]bool)

	for _, w := range weekNumbers {
		num, err := strconv.Atoi(w)
		if err != nil {

			return errors.New("week must contain only numbers")
		}

		if weekMap[num] {
			return errors.New("duplicate week numbers are not allowed")
		}
		weekMap[num] = true
	}
	return nil
}

func validateAlarm(alarm model.Alarm) error {
	if err := util.ValidateDate(alarm.StartAt); err != nil {
		return err
	}
	if err := util.ValidateDate(alarm.EndAt); err != nil {
		return err
	}
	if err := util.ValidateTime(alarm.Timestamp); err != nil {
		return err
	}
	if err := validateWeek(alarm.Week); err != nil {
		return err
	}
	return nil
}

func (sa *saveAlarmService) SaveAlarm(alarm model.Alarm) (string, error) {
	// 유효성 검사 수행
	if err := validateAlarm(alarm); err != nil {
		return "", err
	}
	var existingAlarm model.Alarm
	result := sa.db.First(&existingAlarm, alarm.Id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		alarm.Id = 0
		if err := sa.db.Create(&alarm).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", result.Error
	} else {
		// 레코드가 존재하면 업데이트
		if err := sa.db.Model(&existingAlarm).Updates(alarm).Error; err != nil {
			return "", err
		}
	}

	return "200", nil
}
