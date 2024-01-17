// /diet-service/dto/dto.go
package dto

type GetPresetParams struct {
	Page      int    `form:"page"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type DietPresetRequest struct {
	Id    int
	Uid   int `json:"-"`
	Name  string
	Foods []string
}

type DietPresetResponse struct {
	Id      int
	Name    string
	Foods   []string
	Created string
	Updated string
}

type DietRequest struct {
	Id     int
	Uid    int `json:"-"`
	Name   string
	Time   string `example:"HH:mm"`
	Type   int
	Images []string `example:"base64 encoding string"`
	Foods  []string
}

type DietResponse struct {
	Id      int
	Name    string
	Time    string
	Type    int
	Images  []ImageResponse
	Foods   []string
	Created string
	Updated string
}

type ImageResponse struct {
	Url          string
	ThumbnailUrl string `json:"thumbnail_url"`
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
