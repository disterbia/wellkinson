// /admin-video-service/transport/transport.go
package transport

import (
	"admin-video-service/common/util"
	"admin-video-service/dto"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

var userLocks sync.Map

// @Tags 관리자 동영상 관리 /admin-video
// @Summary 최상위 레벨 조회
// @Description 최초에 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []dto.VimeoLevel1 "웰킨스 폴더 내용"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-items [get]
func GetVimeoLevel1sHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
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

		resp := response.([]dto.VimeoLevel1)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 관리자 동영상 관리 /admin-video
// @Summary 폴더 레벨2 조회
// @Description 폴더내부 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path string true "id"
// @Success 200 {object} []dto.VimeoLevel2 "해당 폴더 내용"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-items/{id} [get]
func GetVimeoLevel2sHandler(getEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		projectId := c.Param("id")
		response, err := getEndpoint(c.Request.Context(), map[string]interface{}{
			"id":        id,
			"projectId": projectId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.VimeoLevel2)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 관리자 동영상 관리 /admin-video
// @Summary 동영상 활성화
// @Description 활성화 동영상 변경시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []string true "활성화 할 id 배열"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-videos/{id} [post]
func SaveHandler(saveEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var videoIds []string // 삭제할 ID 배열
		if err := c.ShouldBindJSON(&videoIds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Concurrent request detected"})
			return
		}
		defer userLocks.Delete(id)

		response, err := saveEndpoint(c.Request.Context(), map[string]interface{}{
			"id":       id,
			"videoIds": videoIds,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}
