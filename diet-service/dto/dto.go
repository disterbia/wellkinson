// /diet-service/dto/dto.go
package dto

type GetPresetParams struct {
	Page      uint   `form:"page"`
	StartDate string `form:"start_date" example:"yyyy-mm-dd"`
	EndDate   string `form:"end_date" example:"yyyy-mm-dd"`
}

type DietPresetRequest struct {
	Id    uint     `json:"id"`
	Uid   uint     `json:"-"`
	Name  string   `json:"name"`
	Foods []string `json:"foods"`
}

type DietPresetResponse struct {
	Id      uint     `json:"id"`
	Name    string   `json:"name"`
	Foods   []string `json:"foods"`
	Created string   `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string   `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type DietRequest struct {
	Id     uint     `json:"id"`
	Uid    uint     `json:"-"`
	Memo   string   `json:"memo"`
	Time   string   `json:"time" example:"HH:mm"`
	Date   string   `json:"date" example:"YYYY-MM-DD"`
	Type   uint     `json:"type"`
	Images []string `json:"images" example:"base64 encoding string"`
	Foods  []string `json:"foods"`
}

type DietCopy struct {
	Id      uint            `json:"id"`
	Memo    string          `json:"memo"`
	Time    string          `json:"time"`
	Type    uint            `json:"type"`
	Date    string          `json:"date" example:"YYYY-MM-DD"`
	Images  []ImageResponse `json:"images"`
	Foods   []string        `json:"foods"`
	Created string          `json:"created" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
	Updated string          `json:"updated" example:"YYYY-mm-ddTHH:mm:ssZ (ISO8601) "`
}

type DietResponse struct {
	Date  string     `json:"date"  example:"YYYY-MM-DD"`
	Diets []DietCopy `json:"diets"`
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
