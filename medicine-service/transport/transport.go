// /medicine-service/transport/transport.go
package transport

import (
	"medicine-service/common/util"
	"medicine-service/dto"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

var userLocks sync.Map

// @Tags 약물 /medicine
// @Summary 약물 저장
// @Description 약물등록 및 수정시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.MedicineRequest true "요청 DTO - 약물데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-medicine [post]
func SaveHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Concurrent request detected"})
			return
		}
		defer userLocks.Delete(id)

		var req dto.MedicineRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Uid = id
		response, err := saveEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 약물 /medicine
// @Summary 약물 삭제
// @Description 약물 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-medicines [post]
func RemoveHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var ids []uint // 삭제할 ID 배열
		if err := c.ShouldBindJSON(&ids); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := removeEndpoint(c.Request.Context(), map[string]interface{}{
			"uid": uid,
			"ids": ids,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 약물 /medicine
// @Summary 약물 복용내역 조회
// @Description 약물 복용내역 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  true  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  true  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []dto.MedicineDateInfo "운동정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-takens [get]
func GetTakensHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var queryParams dto.GetParams
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.MedicineDateInfo)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 약물 /medicine
// @Summary 등록 약물 조회
// @Description 등록 약물 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []dto.MedicineResponse "등록 약물 정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-medicines [get]
func GetMedicinesHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := getEndpoint(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.MedicineResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 약물 /medicine
// @Summary 약물 복용
// @Description 약물 복용시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.TakeMedicine true "약물 복용 데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /take-medicine [post]
func TakeHandler(doEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(uid, true); loaded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Concurrent request detected"})
			return
		}
		defer userLocks.Delete(uid)

		var param dto.TakeMedicine
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		param.Uid = uid
		response, err := doEndpoint(c.Request.Context(), param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 약물 /medicine
// @Summary 약물 복용취소
// @Description 약물 복용취소시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.UnTakeMedicine true "약물 취소 데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /untake-medicine [post]
func UnTakeHandler(doEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(uid, true); loaded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Concurrent request detected"})
			return
		}
		defer userLocks.Delete(uid)

		var param dto.UnTakeMedicine
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		param.Uid = uid
		response, err := doEndpoint(c.Request.Context(), param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 약물 /medicine
// @Summary 약물 찾기
// @Description 약물 검색 키워드 입력시 호출
// @Produce  json
// @Param  keyword  query string  true  "키워드"
// @Success 200 {object} []string "약물명"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /search-medicines [get]
func SearchHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		keyword := c.Query("keyword")

		response, err := getEndpoint(c.Request.Context(), keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]string)
		c.JSON(http.StatusOK, resp)
	}
}
