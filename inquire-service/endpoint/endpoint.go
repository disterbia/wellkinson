// /inquire-service/endpoint/endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"inquire-service/service"

	"github.com/go-kit/kit/endpoint"
)

func AnswerEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(model.InquireReply)
		code, err := s.AnswerInquire(answer)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}

func SendEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(model.Inquire)
		code, err := s.SendInquire(inquire)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}

func GetEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(int)
		inquires, err := s.GetMyInquires(id)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}
