// /user-service/endpoint/endpoint.go

package endpoint

import (
	"common/model"
	"context"
	"user-service/dto"
	"user-service/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeAdminLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		email := reqMap["email"].(string)
		password := reqMap["password"].(string)

		token, err := s.AdminLogin(email, password)

		if err != nil {
			return dto.LoginResponse{Err: err.Error()}, err
		}
		return dto.LoginResponse{Jwt: token}, nil
	}
}

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
		id := request.(uint)
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

func GetMainServicesEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		code, err := s.GetMainServices()
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return code, nil
	}
}

func SendCodeEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		number := request.(string)
		code, err := s.SendAuthCode(number)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func VerifyEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		veri := request.(dto.VerifyRequest)
		code, err := s.VerifyAuthCode(veri.PhoneNumber, veri.Code)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func RemoveEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		code, err := s.RemoveUser(uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func LinkEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(dto.LinkRequest)

		code, err := s.LinkEmail(reqMap.Id, reqMap.IdToken)

		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func MakeAppleLoginEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.LoginRequest)
		token, err := s.AppleLogin(req.IdToken, req.UserRequest)
		if err != nil {
			return dto.LoginResponse{Err: err.Error()}, err
		}
		return dto.LoginResponse{Jwt: token}, nil
	}
}

func GetVersionEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		version, err := s.GetVersion()
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return version, nil
	}
}
