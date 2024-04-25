// /medicine-service/service/util.go
package service

import (
	"encoding/json"
	"errors"
	"medicine-service/common/model"
	"medicine-service/common/util"
	"medicine-service/dto"
)

func validateMedicine(medicine dto.MedicineRequest) error {
	if medicine.IntervalType != 1 {
		if medicine.StartAt != "" {
			if err := util.ValidateDate(medicine.StartAt); err != nil {
				return err
			}
		}
		if medicine.EndAt != "" {
			if err := util.ValidateDate(medicine.EndAt); err != nil {
				return err
			}
		}

		if medicine.Timestamp != nil && len(medicine.Timestamp) != 0 {
			for _, v := range medicine.Timestamp {
				if err := util.ValidateTime(v); err != nil {
					return err
				}
			}
		} else {
			return errors.New(("invalid time format, should be HH:MM"))
		}
	}
	return nil
}

func validateWeek(medicine model.Medicine) (json.RawMessage, []int32, error) {
	if medicine.IntervalType == 1 {
		return json.RawMessage("[]"), []int32{}, nil
	}
	// JSON 배열을 Go 슬라이스로 변환
	weekdays := medicine.Weekdays
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
		} else {
			return nil, nil, errors.New("weekday mus 0-6")
		}
	}

	newWeekdays, err := json.Marshal(unique)
	if err != nil {
		return nil, nil, err
	}
	return newWeekdays, unique, nil
}
