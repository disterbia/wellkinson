// /diet-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"emotion-service/dto"
	"emotion-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveEmotionEndpoint(s service.EmotionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		emotion := request.(dto.EmotionRequest)
		code, err := s.SaveEmotion(emotion)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetEmotionsEndpoint(s service.EmotionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(int)
		queryParams := reqMap["queryParams"].(dto.GetEmotionsParams)
		inquires, err := s.GetEmotions(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func RemoveEmotionsEndpoint(s service.EmotionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]int)
		uid := reqMap["uid"].(int)
		code, err := s.RemoveEmotions(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
