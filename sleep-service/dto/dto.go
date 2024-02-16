// /sleep-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type SleepAlarmRequest struct {
	Id        uint   `json:"id"`
	Uid       uint   `json:"-"`
	StartTime string `json:"start_time" example:"HH:mm"`
	EndTime   string `json:"end_time" example:"HH:mm"`
	AlarmTime string `json:"alarm_time" example:"HH:mm"`
	Weekdays  []uint `json:"weekdays"`
	IsActive  bool   `json:"is_active"`
}

type SleepAlarmResponse struct {
	Id        uint   `json:"id"`
	StartTime string `json:"start_time" example:"HH:mm"`
	EndTime   string `json:"end_time" example:"HH:mm"`
	Weekdays  []uint `json:"weekdays"`
	IsActive  bool   `json:"is_active"`
	Created   string `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated   string `json:"updated"  example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type SleepTimeRequest struct {
	Uid       uint   `json:"-"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time" example:"HH:mm"`
	DateSleep string `json:"date_sleep" example:"YYYY-MM-DD"`
}

type SleepTimeResponse struct {
	Id        uint   `json:"id"`
	StartTime string `json:"start_time" example:"HH:mm"`
	EndTime   string `json:"end_time" example:"HH:mm"`
	DateSleep string `json:"date_sleep" example:"YYYY-MM-DD"`
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
