// /diet-service/transport/transport.go
package transport

import (
	"common/util"
	"diet-service/dto"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

var userLocks sync.Map

// @Tags 식단
// @Summary 추가한 식단 생성/수정
// @Description 추가한 식단 생성시 Id 생략
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.DietPresetRequest true "요청 DTO - 추가한 식단 데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-preset [post]
func SavePresetHandler(savePresetEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var req dto.DietPresetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.Uid = id
		response, err := savePresetEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 식단
// @Summary 추가한 식단 조회
// @Description 추가한 식단 조회시 호출 (10개씩)
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  page  query  uint  false  "페이지 번호 default 0" (10개씩)
// @Param  start_date  query string  false  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  false  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []dto.DietPresetResponse "식단정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-presets [get]
func GetPresetsHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var queryParams dto.GetPresetParams
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// id와 queryParams를 함께 전달
		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.DietPresetResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 식단
// @Summary 추가한 식단 삭제
// @Description 추가한 식단 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-preset [post]
func RemovePresetHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

// @Tags 식단
// @Summary 식단생성/수정
// @Description 식단생성시 Id 생략
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.DietRequest true "요청 DTO - 식단데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-diet [post]
func SaveDietHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var req dto.DietRequest
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

// @Tags 식단
// @Summary 식단 조회
// @Description 식단 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  false  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  false  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []dto.DietResponse "식단정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-diets [get]
func GetDietsHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var queryParams dto.GetPresetParams
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// id와 queryParams를 함께 전달
		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.DietResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 식단
// @Summary 추가한 식단 삭제
// @Description 추가한 식단 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-diet [post]
func RemoveDietHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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
