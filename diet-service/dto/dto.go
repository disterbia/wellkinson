// /diet-service/dto/dto.go
package dto

type DietPresetRequest struct {
	Id    int
	Uid   int `json:"-"`
	Name  string
	Foods []FoodRequest
}

type DietPresetResponse struct {
	Id      int
	Uid     int `json:"-"`
	Name    string
	Foods   []FoodResponse
	Created string
	Updated string
}

type FoodRequest struct {
	Id   int
	Name string
}

type FoodResponse struct {
	Id       int
	PresetId int `json:"preset_id"`
	Name     string
	Created  string
	Updated  string
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
