// /vocal-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"vocal-service/dto"
	"vocal-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveScoresEndpoint(s service.VocalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		scores := request.([]dto.VocalScoreRequest)
		code, err := s.SaveVocalScores(scores)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetScoresEndpoint(s service.VocalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetParams)
		faceScores, err := s.GetVocalScores(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return faceScores, nil
	}
}

func GetVocalTablesEndpoint(s service.VocalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		faceScores, err := s.GetVoiceTables()
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return faceScores, nil
	}
}
