// /emotion-service/dto/dto.go
package dto

type GetEmotionsParams struct {
	StartDate string `form:"start_date" example:"yyyy-mm-dd"`
	EndDate   string `form:"end_date" example:"yyyy-mm-dd"`
}
type EmotionRequest struct {
	Id      uint    `json:"id"`
	Uid     uint    `json:"-"`
	Emotion *uint   `json:"emotion"`
	State   *string `json:"state" example:"기분내용"`
}

type EmotionResponse struct {
	Id      uint   `json:"id"`
	Emotion uint   `json:"emotion" `
	State   string `json:"state" example:"기분내용"`
	Created string `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
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
