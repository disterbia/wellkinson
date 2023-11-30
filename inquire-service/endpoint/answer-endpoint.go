// /inquire-service/endpoint/save-endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"inquire-service/service"

	"github.com/go-kit/kit/endpoint"
)

func AnswerEndpoint(s service.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(model.InquireReply)
		code, err := s.AnswerInquire(answer)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}
