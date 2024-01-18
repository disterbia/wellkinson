// /diet-service/dto/dto.go
package dto

type GetPresetParams struct {
	Page      int    `form:"page"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type DietPresetRequest struct {
	Id    int      `json:"id"`
	Uid   int      `json:"-"`
	Name  string   `json:"name"`
	Foods []string `json:"foods"`
}

type DietPresetResponse struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Foods   []string `json:"foods"`
	Created string   `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string   `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type DietRequest struct {
	Id     int      `json:"id"`
	Uid    int      `json:"-"`
	Name   string   `json:"name"`
	Time   string   `json:"time" example:"HH:mm"`
	Type   int      `json:"type"`
	Images []string `json:"images" example:"base64 encoding string"`
	Foods  []string `json:"foods"`
}

type DietResponse struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Time    string          `json:"time"`
	Type    int             `json:"type"`
	Images  []ImageResponse `json:"images"`
	Foods   []string        `json:"foods"`
	Created string          `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string          `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type ImageResponse struct {
	Url          string `json:"url"`
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
