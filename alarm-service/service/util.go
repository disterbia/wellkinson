// /alarm-service/service/util.go
package service

import (
	"alarm-service/common/util"
	"alarm-service/dto"
)

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
	return nil
}
