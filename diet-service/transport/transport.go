// /diet-service/transport/transport.go
package transport

import (
	"common/util"
	"diet-service/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	kitEndpoint "github.com/go-kit/kit/endpoint"
)

// @Summary 식단
// @Tags 식단생성/수정
// @Description 식단생성시 Id 생략
// @Produce  json
// @Param request body dto.DietPresetRequest true "요청 DTO - 답변데이터"
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

// @Summary 식단
// @Tags 식단조회
// @Description 식단조회시 호출 (10개씩)
// @Produce  json
// @Param  page  query int false  "페이지 번호 default 0"
// @Success 200 {object} []dto.DietPresetResponse "성공시 200 반환"
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
		pageParam := c.Query("page")
		var page int
		if pageParam != "" {
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
				return
			}
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

		resp := response.([]dto.DietPresetResponse)
		c.JSON(http.StatusOK, resp)
	}
}
