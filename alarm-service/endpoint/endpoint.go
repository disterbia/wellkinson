// /alarm-service/endpoint/endpoint.go

package endpoint

import (
	"alarm-service/dto"
	"alarm-service/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveAlarmEndpoint(s service.AlarmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		alarm := request.(dto.AlarmRequest)
		code, err := s.SaveAlarm(alarm)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func RemoveAlarmEndpoint(s service.AlarmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(int)
		uid := reqMap["uid"].(int)
		code, err := s.RemoveAlarm(id, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetEndpoint(s service.AlarmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(int)
		page := reqMap["page"].(int)
		inquires, err := s.GetAlarms(id, page)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}
