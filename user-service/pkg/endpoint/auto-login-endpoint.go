// /user-service/pkg/endpoint/auto-login-endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"log"
	"user-service/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeAutoLoginEndpoint(s service.LoginService) endpoint.Endpoint {
	log.Println("endpoint: 호출")
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		email := request.(string)
		token, err := s.Login(email, model.User{})
		if err != nil {
			return model.LoginResponse{Err: err.Error()}, err
		}
		log.Println("endpoint: 완료")
		return model.LoginResponse{Jwt: token}, nil
	}
}
