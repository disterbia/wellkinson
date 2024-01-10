// /user-service/endpoint/endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"user-service/dto"
	"user-service/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeAutoLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		email := request.(string)
		token, err := s.AutoLogin(email, model.User{})
		if err != nil {
			return dto.LoginResponse{Err: err.Error()}, err
		}
		return dto.LoginResponse{Jwt: token}, nil
	}
}

func MakeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(int)
		result, err := s.GetUser(id)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return result, nil
	}
}
func MakeGoogleLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.LoginRequest)
		token, err := s.GoogleLogin(req.IdToken, req.UserRequest)
		if err != nil {
			return dto.LoginResponse{Err: err.Error()}, err
		}
		return dto.LoginResponse{Jwt: token}, nil
	}
}
func MakeKakaoLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.LoginRequest)
		token, err := s.KakaoLogin(req.IdToken, req.UserRequest)
		if err != nil {
			return dto.LoginResponse{Err: err.Error()}, err
		}
		return dto.LoginResponse{Jwt: token}, nil
	}
}
func MakeSetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := request.(dto.UserRequest)
		code, err := s.SetUser(user)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
