package dto

type AlarmRequest struct {
	Id        int    `json:"id"`
	Uid       int    `json:"-"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      []int  `json:"week"`
}

type AlarmResponse struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      []int  `json:"week"`
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
