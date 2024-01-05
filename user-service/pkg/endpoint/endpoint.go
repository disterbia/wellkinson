// /user-service/pkg/endpoint/endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"log"
	"user-service/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeAutoLoginEndpoint(s service.UserService) endpoint.Endpoint {
	log.Println("endpoint: 호출")
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		email := request.(string)
		token, err := s.AutoLogin(email, model.User{})
		if err != nil {
			return model.LoginResponse{Err: err.Error()}, err
		}
		log.Println("endpoint: 완료")
		return model.LoginResponse{Jwt: token}, nil
	}
}

func MakeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(int)
		result, err := s.GetUser(id)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return result, nil
	}
}
func MakeGoogleLoginEndpoint(s service.UserService) endpoint.Endpoint {
	log.Println("endpoint: 호출")
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.LoginRequest)
		token, err := s.GoogleLogin(req.IdToken, req.User)
		if err != nil {
			return model.LoginResponse{Err: err.Error()}, err
		}
		log.Println("endpoint: 완료")
		return model.LoginResponse{Jwt: token}, nil
	}
}
func MakeKakaoLoginEndpoint(s service.UserService) endpoint.Endpoint {
	log.Println("endpoint: 호출")
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.LoginRequest)
		token, err := s.KakaoLogin(req.IdToken, req.User)
		if err != nil {
			return model.LoginResponse{Err: err.Error()}, err
		}
		log.Println("endpoint: 완료")
		return model.LoginResponse{Jwt: token}, nil
	}
}
func MakeSetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := request.(model.User)
		code, err := s.SetUser(user)
		if err != nil {
			return model.BasicResponse{Code: err.Error()}, err
		}
		return model.BasicResponse{Code: code}, nil
	}
}
