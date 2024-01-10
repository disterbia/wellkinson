// /diet-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"diet-service/dto"
	"diet-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SavePresetEndpoint(s service.DietPresetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		dietPreset := request.(dto.DietPresetRequest)
		code, err := s.SaveDietPreset(dietPreset)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
