// /sleep-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type SleepAlarmRequest struct {
	Id        uint   `json:"id"`
	Uid       uint   `json:"-"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	AlarmTime string `json:"alarm_time"`
	Weekdays  []uint `json:"weekdays"`
	IsActive  bool   `json:"is_active"`
}

type SleepAlarmResponse struct {
	Id        uint   `json:"id"`
	StartTime string `json:"date_taken"`
	EndTime   string `json:"time_taken"`
	Weekdays  []uint `json:"weekdays"`
	IsActive  bool   `json:"is_active"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

type SleepTimeRequest struct {
	Uid       uint   `json:"-"`
	StartTime string `json:"date_taken"`
	EndTime   string `json:"time_taken"`
	DateSleep string `json:"date_sleep"`
}

type SleepTimeResponse struct {
	StartTime string `json:"date_taken"`
	EndTime   string `json:"time_taken"`
	DateSleep string `json:"date_sleep"`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}
type ErrorResponse struct {
	Err string `json:"err"`
}

type BasicResponse struct {
	Code string `json:"code"`
}
