// /admin-video-service/endpoint/endpoint.go
package endpoint

import (
	"admin-video-service/dto"
	"admin-video-service/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

func GetVimeoLevel1sEndpoint(s service.AdminVideoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		level1s, err := s.GetLevel1s(id)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return level1s, nil
	}
}

func GetVimeoLevel2sEndpoint(s service.AdminVideoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		projectId := reqMap["projectId"].(string)
		level1s, err := s.GetLevel2s(id, projectId)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return level1s, nil
	}
}

func SaveEndpoint(s service.AdminVideoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		code, err := s.SaveVideos(request.(dto.VideoData))
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}
