// /alarm-service/endpoint/remove-endpoint.go

package endpoint

import (
	"alarm-service/service"
	"common/model"
	"context"

	"github.com/go-kit/kit/endpoint"
)

func RemoveAlarmEndpoint(s service.RemoveAlarmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		alarm := request.(model.Alarm)
		code, err := s.RemoveAlarm(alarm)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}
