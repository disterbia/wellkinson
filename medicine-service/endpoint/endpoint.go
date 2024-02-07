// /medicine-service/endpoint/endpoint.go
package endpoint

import (
	"context"
	"medicine-service/dto"
	"medicine-service/service"

	"github.com/go-kit/kit/endpoint"
)

func SaveEndpoint(s service.MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		medicine := request.(dto.MedicineRequest)
		code, err := s.SaveMedicine(medicine)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func RemoveEndpoint(s service.MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveMedicines(ids, uid)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return dto.BasicResponse{Code: code}, nil
	}
}

func GetTakensEndpoint(s service.MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(dto.GetParams)
		inquires, err := s.GetTakens(id, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetMedicinesEndpoint(s service.MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		medicines, err := s.GetMedicines(id)
		if err != nil {
			return dto.BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}
