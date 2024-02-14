// /exercise-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"exercise-service/dto"
	"exercise-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveExerciseEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(dto.ExerciseRequest)
		code, err := s.SaveExercise(exercise)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetExercisesEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetParams)
		inquires, err := s.GetExercises(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func RemoveExercisesEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveExercises(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func DoExerciseEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(dto.ExerciseDo)
		code, err := s.DoExercise(exercise)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetProjectsEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		projects, err := s.GetProjects()
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return projects, nil
	}
}

func GetVideosEndpoint(s service.ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(dto.GetVideoParams)
		videos, err := s.GetVideos(reqMap.ProjectId, reqMap.Page)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return videos, nil
	}
}
