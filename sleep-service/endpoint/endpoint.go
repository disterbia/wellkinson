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

func GetSleepAlarmsEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		SleepAlarms, err := s.GetSleepAlarms(id)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return SleepAlarms, nil
	}
}

func RemoveSleepAlarmsEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveSleepAlarms(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetSleepTimesEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetParams)
		sleepTimes, err := s.GetSleepTimes(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return sleepTimes, nil
	}
}

func SaveSleepTimeEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sleepTime := request.(dto.SleepTimeRequest)
		code, err := s.SaveSleepTime(sleepTime)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func RemoveSleepTimeEndpoint(s service.SleepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveSleepTime(id, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
