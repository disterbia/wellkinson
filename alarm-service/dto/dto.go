package dto

type AlarmRequest struct {
	Id        int    `json:"id"`
	Uid       int    `json:"-"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      string `json:"week" example:"0,4,6 (sunday:0,...)"`
}

type AlarmResponse struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      string `json:"week" example:"0,4,6 (sunday:0,...)"`
	Created   string `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated   string `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

type ErrorResponse struct {
	Err string `json:"err" `
}

type BasicResponse struct {
	Code string `json:"code"`
}
