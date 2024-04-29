// /exercise-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `form:"end_date" example:"YYYY-MM-DD"`
}

type GetVideoParams struct {
	Page      uint   `form:"page"`
	ProjectId string `form:"project_id"`
}

type ExerciseRequest struct {
	Id              uint   `json:"id"`
	Uid             uint   `json:"-"`
	Title           string `json:"title"`
	ExerciseStartAt string `json:"exercise_start_at" example:"HH:mm"`
	ExerciseEndAt   string `json:"exercise_end_at" example:"HH:mm"`
	PlanStartAt     string `json:"plan_start_at" example:"YYYY-MM-DD"`
	PlanEndAt       string `json:"plan_end_at" example:"YYYY-MM-DD"`
	UseAlarm        *bool  `json:"use_alarm"`
	Weekdays        []uint `json:"weekdays"`
}

type ExerciseResponse struct {
	Id              uint   `json:"id"`
	Title           string `json:"title"`
	ExerciseStartAt string `json:"exercise_start_at" example:"HH:mm"`
	ExerciseEndAt   string `json:"exercise_end_at"  example:"HH:mm"`
	PlanStartAt     string `json:"plan_start_at"  example:"YYYY-MM-DD"`
	PlanEndAt       string `json:"plan_end_at"  example:"YYYY-MM-DD"`
	UseAlarm        bool   `json:"use_alarm"`
	Repeat          uint   `json:"repeat"`
	Weekdays        []uint `json:"weekdays"`
	Created         string `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated         string `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type ExerciseDateInfo struct {
	Date      string             `json:"date" example:"YYYY-MM-DD"`
	Exercises []ExerciseDoneInfo `json:"exercises"`
}

type ExerciseDoneInfo struct {
	Exercise ExerciseResponse `json:"exercise"`
	Done     bool             `json:"done"`
}

type ExerciseDo struct {
	Uid           uint   `json:"-"`
	ExerciseId    uint   `json:"exercise_id"`
	PerformedDate string `json:"performed_date"  example:"YYYY-MM-DD"`
}

type ProjectResponse struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
	Count     uint   `json:"count"`
}

type VideoResponse struct {
	Name         string `json:"name"`
	VideoId      string `json:"video_id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Duration     uint   `json:"duration"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
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
