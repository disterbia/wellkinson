// /sleep-service/transport/transport.go
package transport

import (
	"net/http"
	"sleep-service/common/util"
	"sleep-service/dto"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

var userLocks sync.Map

// @Tags 수면 /sleep
// @Summary 수면알림 생성/수정
// @Description 수면알림 생성시 Id 생략
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.SleepAlarmRequest true "수면알림 DTO - 수면알림 데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-sleep [post]
func SaveSleepHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var req dto.SleepAlarmRequest
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

// @Tags 수면 /sleep
// @Summary 수면알람 조회
// @Description 수만알람 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []dto.SleepAlarmResponse "수면 알람정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-sleep-alarms [get]
func GetSleepAlarmsHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		resp := response.([]dto.SleepAlarmResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 수면 /sleep
// @Summary 수면알림 삭제
// @Description 수면 알림 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-sleep-alarms [post]
func RemoveSleepAlarmsHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

// @Tags 수면 /sleep
// @Summary 기간별 수면시간 조회
// @Description 기간별 수면시간 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  true  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  true  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []dto.SleepTimeResponse "수면시간 정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-sleep-times [get]
func GetSleepTimesHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		// id와 queryParams를 함께 전달
		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.SleepTimeResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 수면 /sleep
// @Summary 수면시간 등록/수정
// @Description 수면시간 등록/수정 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.SleepTimeRequest true "수면시간 DTO - 수면시간 데이터"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-sleep-time [post]
func SaveSleepTimeHandler(doEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var param dto.SleepTimeRequest
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

// @Tags 수면 /sleep
// @Summary 수면시간 삭제
// @Description 수면 시간 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path uint ture "수면시간 ID"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-sleep-times/{id} [post]
func RemoveSleepTimeHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tid := c.Param("id")
		id, _ := strconv.ParseUint(tid, 10, 32)
		response, err := removeEndpoint(c.Request.Context(), map[string]interface{}{
			"uid": uid,
			"id":  uint(id),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)

	}
}
