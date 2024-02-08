// /sleep-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"sleep-service/dto"
	"sleep-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveSleepAlarmEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sleep := request.(dto.SleepAlarmRequest)
		code, err := s.SaveSleepAlarm(sleep)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

// func GetExercisesEndpoint(s service.ExerciseService) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		reqMap := request.(map[string]interface{})
// 		id := reqMap["id"].(uint)
// 		queryParams := reqMap["queryParams"].(dto.GetParams)
// 		inquires, err := s.GetExercises(id, queryParams.StartDate, queryParams.EndDate)
// 		if err != nil {
// 			return dto.BasicResponse{Code: err.Error()}, err
// 		}
// 		return inquires, nil
// 	}
// }

// func RemoveExercisesEndpoint(s service.ExerciseService) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		reqMap := request.(map[string]interface{})
// 		ids := reqMap["ids"].([]uint)
// 		uid := reqMap["uid"].(uint)
// 		code, err := s.RemoveExercises(ids, uid)
// 		if err != nil {
// 			return dto.BasicResponse{Code: err.Error()}, err
// 		}
// 		return dto.BasicResponse{Code: code}, nil
// 	}
// }

// func DoExerciseEndpoint(s service.ExerciseService) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		exercise := request.(dto.ExerciseDo)
// 		code, err := s.DoExercise(exercise)
// 		if err != nil {
// 			return dto.BasicResponse{Code: err.Error()}, err
// 		}
// 		return dto.BasicResponse{Code: code}, nil
// 	}
// }
