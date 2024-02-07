// /face-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type FaceScoreRequest struct {
	Uid   uint `json:"-"`
	Score uint `json:"score"`
	Type  uint `json:"type"`
}

type FaceScoreResponse struct {
	Score   uint   `json:"score"`
	Type    uint   `json:"type"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type FaceExamResponse struct {
	Type    uint
	Title   string
	VideoId string `json:"video_id"`
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
