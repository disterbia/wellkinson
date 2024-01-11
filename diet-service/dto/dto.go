// /diet-service/dto/dto.go
package dto

type DietPresetRequest struct {
	Id    int
	Uid   int `json:"-"`
	Name  string
	Foods []string
}

type DietPresetUpdate struct {
	Name  string
	Foods string
}

type DietPresetResponse struct {
	Id      int
	Name    string
	Foods   []string
	Created string
	Updated string
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
