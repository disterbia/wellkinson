// /user-service/pkg/endpoint/auto-login-endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"user-service/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeSetUserEndpoint(s service.SetUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := request.(model.User)
		code, err := s.SetUser(user)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}
