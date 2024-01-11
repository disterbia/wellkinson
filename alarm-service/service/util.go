// /alarm-service/service/util.go
package service

import (
	"alarm-service/dto"
	"common/util"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

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

func validateAlarm(alarm dto.AlarmRequest) error {
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
