package dto

type GetInquireParams struct {
	Page      int    `form:"page"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type InquireRequest struct {
	Uid     int    `json:"-"`
	Email   string `json:"email"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InquireResponse struct {
	Id      int                    `json:"id"`
	Email   string                 `json:"email"`
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	Replies []InquireReplyResponse `json:"replies"`
}

type InquireReplyRequest struct {
	Id        int    `json:"-"`
	Uid       int    `json:"-"`
	InquireId int    `json:"inquire_id"`
	Content   string `json:"content"`
	ReplyType bool   `json:"reply_type"`
}

type InquireReplyResponse struct {
	Id        int    `json:"id"`
	InquireId int    `json:"inquire_id"`
	Content   string `json:"content"`
	ReplyType bool   `json:"reply_type"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
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
