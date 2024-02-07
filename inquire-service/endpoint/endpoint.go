// /inquire-service/endpoint/endpoint.go

package endpoint

import (
	"context"
	"inquire-service/dto"
	"inquire-service/service"

	"github.com/go-kit/kit/endpoint"
)

func AnswerEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(dto.InquireReplyRequest)
		code, err := s.AnswerInquire(answer)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func SendEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(dto.InquireRequest)
		code, err := s.SendInquire(inquire)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func RemoveInquireEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveInquire(id, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
func RemoveReplyEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveReply(id, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetInquireParams)
		inquires, err := s.GetMyInquires(id, queryParams.Page, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetAllEndpoint(s service.InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetInquireParams)
		inquires, err := s.GetAllInquires(id, queryParams.Page, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}
