// /vocal-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type VocalScoreRequest struct {
	Uid   uint `json:"-"`
	Score uint `json:"score"`
	Type  uint `json:"type"`
}

type VocalScoreResponse struct {
	Score   uint   `json:"score"`
	Type    uint   `json:"type"`
	Created string `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated string `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type VoiceWordResponse struct {
	Type  uint   `json:"type"`
	Title string `json:"title"`
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
