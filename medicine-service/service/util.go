// /medicine-service/service/util.go
package service

import (
	"common/util"
	"encoding/json"
	"medicine-service/dto"
)

func validateMedicine(medicine dto.MedicineRequest) error {
	if err := util.ValidateDate(medicine.StartAt); err != nil {
		return err
	}
	if err := util.ValidateDate(medicine.EndAt); err != nil {
		return err
	}
	for _, v := range medicine.Timestamp {
		if err := util.ValidateTime(v); err != nil {
			return err
		}
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
