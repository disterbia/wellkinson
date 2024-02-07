// /face-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"face-service/dto"
	"face-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveScoresEndpoint(s service.FaceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		scores := request.([]dto.FaceScoreRequest)
		code, err := s.SaveFaceScores(scores)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetScoresEndpoint(s service.FaceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetParams)
		faceScores, err := s.GetFaceScores(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return faceScores, nil
	}
}

func GetFaceExamsEndpoint(s service.FaceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		faceScores, err := s.GetFaceExams()
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return faceScores, nil
	}
}
