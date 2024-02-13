// /sleep-service/service/util.go
package service

import (
	"common/util"
	"encoding/json"
	"errors"
	"sleep-service/dto"
)

func validateSleep(sleepRequest dto.SleepAlarmRequest) error {
	if err := util.ValidateTime(sleepRequest.StartTime); err != nil {
		return err
	}
	if err := util.ValidateTime(sleepRequest.EndTime); err != nil {
		return err
	}
	if sleepRequest.IsActive {
		if err := util.ValidateTime(sleepRequest.AlarmTime); err != nil {
			return err
		}
	}

	return nil
}

func validateSleepTime(sleepRequest dto.SleepTimeRequest) error {
	if err := util.ValidateTime(sleepRequest.StartTime); err != nil {
		return err
	}
	if err := util.ValidateTime(sleepRequest.EndTime); err != nil {
		return err
	}
	if err := util.ValidateDate(sleepRequest.DateSleep); err != nil {
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

func isDuplicateWeek(userWeekdays []uint, dbWeekdaysRaw []json.RawMessage) error {
	for _, rawDay := range dbWeekdaysRaw {
		var dbWeekdays []uint
		err := json.Unmarshal(rawDay, &dbWeekdays)
		if err != nil {
			return err // JSON 파싱 오류
		}

		// 두 배열 간 중복 요일 확인
		for _, uDay := range userWeekdays {
			for _, dbDay := range dbWeekdays {
				if uDay == dbDay {
					return errors.New("중복되는 요일이 있습니다")
				}
			}
		}
	}

	return nil // 중복 없음
}
