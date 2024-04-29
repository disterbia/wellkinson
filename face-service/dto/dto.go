// /face-service/dto/dto.go
package dto

import (
	"face-service/common/model"
)

type GetParams struct {
	StartDate string `form:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `form:"end_date" example:"YYYY-MM-DD"`
}

type FaceScoreRequest struct {
	Uid   uint `json:"-"`
	Score uint `json:"score"`
	Type  uint `json:"type"`
}

type FaceScoreResponse struct {
	Score   uint   `json:"score"`
	Type    uint   `json:"type"`
	Created string `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated string `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type FaceExamResponse struct {
	Type    uint   `json:"type"`
	Title   string `json:"title"`
	VideoId string `json:"video_id"`
}

type SwaggerExercise struct {
	Id      uint
	Type    uint
	Title   string
	VideoId string `json:"video_id"`
	Created string
	Updated string
}

type FaceExerciseResponse struct {
	Type         uint                 `json:"type"`
	Title        string               `json:"title"`
	Count        uint                 `json:"count"`
	FaceExercise []model.FaceExercise `json:"videos" gorm:"-"`
}

type SwaggerResponse struct {
	Type         uint
	Title        string
	Count        uint
	FaceExercise SwaggerExercise `json:"videos"`
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
