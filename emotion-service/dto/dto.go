// /emotion-service/dto/dto.go
package dto

type GetEmotionsParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type EmotionRequest struct {
	Id      int    `json:"id"`
	Uid     int    `json:"-"`
	Emotion string `json:"emotion"`
	State   string `json:"state"`
}

type EmotionResponse struct {
	Id      int    `json:"id"`
	Emotion string `json:"emotion"`
	State   string `json:"state"`
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
