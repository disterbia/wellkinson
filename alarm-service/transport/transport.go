// /alarm-service/transport/transport.go

package transport

import (
	"alarm-service/dto"
	"common/util"
	"net/http"
	"strconv"
	"sync"

	kitEndpoint "github.com/go-kit/kit/endpoint"

	"github.com/gin-gonic/gin"
)

var userLocks sync.Map

// @Tags 알람 /alarm
// @Summary  알람생성/수정
// @Description 알람생성시 id생략 / 알람수정시 id 포함
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.AlarmRequest true "요청 DTO - 알람데이터, type 1:운동 2:약 3:수면"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-alarm [post]
func SaveAlarmHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var req dto.AlarmRequest
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

// @Tags 알람 /alarm
// @Summary 알람삭제
// @Description 알람삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-alarms [post]
func RemoveAlarmsHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

// @Tags 알람 /alarm
// @Summary 등록된 알람 조회
// @Description 등록된 알람조회시 호출 (10개씩)
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  page  query uint false  "페이지 번호 default 0"
// @Success 200 {object} []dto.AlarmResponse "알람정보 - type 1:운동 2:약 3:수면"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-presets [get]
func GetHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		pageParam := c.Query("page")
		var page uint
		if pageParam != "" {
			parsed, err := strconv.ParseUint(pageParam, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
				return
			}
			page = uint(parsed)
		}

		// id와 queryParams를 함께 전달
		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":   id,
			"page": page,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.AlarmResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 알람 /alarm
// @Summary 수신 알림 조회
// @Description 수신 알림 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []dto.AlarmResponse "알람정보 - type 1:운동 2:약 3:수면"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-notis [get]
func GetNotisHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		resp := response.([]dto.NotificationResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 알람 /alarm
// @Summary  알림 모두 읽기
// @Description 수신 알림 조회시 자동 읽음 처리라면 수신 알림 조회완료 후 함께 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /read-notis [post]
func ReadAllHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response, err := saveEndpoint(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 알람 /alarm
// @Summary 수신된 알림 삭제
// @Description 수신된 알림 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-notis [post]
func RemoveNoitsHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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
