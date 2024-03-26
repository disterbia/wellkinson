// /alarm-service/service/util.go
package service

import (
	"alarm-service/common/util"
	"alarm-service/dto"
	"encoding/json"
	"errors"
)

func validateAlarm(alarm dto.AlarmRequest) error {
	if alarm.StartAt != "" {
		if err := util.ValidateDate(alarm.StartAt); err != nil {
			return err
		}
	}
	if alarm.EndAt != "" {
		if err := util.ValidateDate(alarm.EndAt); err != nil {
			return err
		}
	}

	if err := util.ValidateTime(alarm.Timestamp); err != nil {
		return err
	}
	return nil
}

func validateWeek(weekdays json.RawMessage) (json.RawMessage, []int32, error) {
	// JSON 배열을 Go 슬라이스로 변환
	var weekdaySlice []int32
	err := json.Unmarshal(weekdays, &weekdaySlice)
	if err != nil {
		return nil, nil, err
	}

	if len(weekdaySlice) == 0 {
		return nil, nil, errors.New("must weekday")
	}
	seen := make(map[int32]bool)
	unique := []int32{}

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
		return nil, nil, err
	}
	return newWeekdays, unique, nil
}
