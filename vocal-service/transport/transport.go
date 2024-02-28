// /vocal-service/transport/transport.go
package transport

import (
	"net/http"
	"sync"
	"vocal-service/common/util"
	"vocal-service/dto"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

var userLocks sync.Map

// @Tags 발성 /vocal
// @Summary 발성 점수 저장
// @Description 발성 검사 완료 후 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []dto.VocalScoreRequest true "요청 DTO - 발성검사 데이터 ( type: 1:a 2:e 3:i 4:o 5:u )"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-vocals [post]
func SaveScoresHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		var req []dto.VocalScoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req[0].Uid = id
		response, err := saveEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 발성 /vocal
// @Summary 발성 점수 조회
// @Description 발성 점수 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  false  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  false  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []dto.VocalScoreResponse "표정검사 점수 정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-vocal-scores [get]
func GetScoresHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		resp := response.([]dto.VocalScoreResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 발성 /vocal
// @Summary 발성운동 단어 조회
// @Description 발성운동 단어 조회시 호출
// @Produce  json
// @Success 200 {object} []dto.VoiceWordResponse "발성운동 단어 데이터"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-voice-tables [get]
func GetVocalTablesHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		response, err := getEndpoint(c.Request.Context(), false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.VoiceWordResponse)
		c.JSON(http.StatusOK, resp)
	}
}
