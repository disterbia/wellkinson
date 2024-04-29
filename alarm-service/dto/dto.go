package dto

type AlarmRequest struct {
	Id        uint   `json:"id"`
	Uid       uint   `json:"-"`
	Type      uint   `json:"type"`
	Body      string `json:"body" example:"알람내용"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      []uint `json:"week"`
}

type AlarmResponse struct {
	Id        uint   `json:"id"`
	Type      uint   `json:"type"`
	Body      string `json:"body" example:"알람내용"`
	StartAt   string `json:"start_at" example:"yyyy-mm-dd"`
	EndAt     string `json:"end_at" example:"yyyy-mm-dd"`
	Timestamp string `json:"timestamp" example:"HH:mm"`
	Week      []uint `json:"week"`
	Created   string `json:"created" example:"YYYY-mm-ddTHH:mm:ss "`
	Updated   string `json:"updated" example:"YYYY-mm-ddTHH:mm:ss "`
}

type NotificationResponse struct {
	Id       uint   `json:"id"`
	Type     uint   `json:"type"`
	Body     string `json:"body" example:"알람내용"`
	ParentId uint   `json:"parent_id"`
	IsRead   bool   `json:"is_read"`
	Created  string `json:"created" example:"YYYY-mm-ddTHH:mm:ss "`
	Updated  string `json:"updated" example:"YYYY-mm-ddTHH:mm:ss "`
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
