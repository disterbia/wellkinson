package dto

type GetInquireParams struct {
	Page      uint   `form:"page"`
	StartDate string `form:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `form:"end_date" example:"YYYY-MM-DD"`
}

type InquireRequest struct {
	Uid     uint   `json:"-"`
	Email   string `json:"email"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InquireResponse struct {
	Id      uint                   `json:"id"`
	Email   string                 `json:"email"`
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Created string                 `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string                 `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Replies []InquireReplyResponse `json:"replies"`
}

type InquireReplyRequest struct {
	Id        uint   `json:"-"`
	Uid       uint   `json:"-"`
	InquireId uint   `json:"inquire_id"`
	Content   string `json:"content"`
	ReplyType bool   `json:"reply_type"`
}

type InquireReplyResponse struct {
	Id        uint   `json:"id"`
	InquireId uint   `json:"inquire_id"`
	Content   string `json:"content"`
	ReplyType bool   `json:"reply_type"`
	Created   string `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated   string `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}
type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

type ErrorResponse struct {
	Err string `json:"err"` // wwwwww
}

type BasicResponse struct {
	Code string `json:"code"`
}
