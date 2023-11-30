// /user-service/pkg/endpoint/auto-login-endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"user-service/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeGetUserEndpoint(s service.GetUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(int)
		result, err := s.GetUser(id)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return result, nil
	}
}
