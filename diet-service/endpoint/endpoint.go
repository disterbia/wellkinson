// /diet-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"diet-service/dto"
	"diet-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SavePresetEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		dietPreset := request.(dto.DietPresetRequest)
		code, err := s.SavePreset(dietPreset)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetPresetsEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(int)
		queryParams := reqMap["queryParams"].(dto.GetPresetParams)
		inquires, err := s.GetPresets(id, queryParams.Page, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func RemovePresetsEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]int)
		uid := reqMap["uid"].(int)
		code, err := s.RemovePresets(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func SaveDietEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		diet := request.(dto.DietRequest)
		code, err := s.SaveDiet(diet)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetDietsEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(int)
		queryParams := reqMap["queryParams"].(dto.GetPresetParams)
		inquires, err := s.GetDiets(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func RemoveDietsEndpoint(s service.DietService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]int)
		uid := reqMap["uid"].(int)
		code, err := s.RemoveDiets(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
