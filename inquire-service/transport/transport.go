// /inquire-service/transport/transport.go

package transport

import (
	"common/model"
	"common/util"
	"log"
	"net/http"

	kitEndpoint "github.com/go-kit/kit/endpoint"

	"github.com/gin-gonic/gin"
)

// @Summary 알람설정
// @Tags 알람생성
// @Description 답변등록시 호출
// @Accept  json
// @Produce  json
// @Param request body InquireReply true "요청 DTO - 답변데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /inquire-answer [post]
func AnswerHandler(answerEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	log.Println("transport:시작")
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var req *model.InquireReply
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.Uid = id
		response, err := answerEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.BasicResponse)
		c.JSON(http.StatusOK, resp)
		log.Println("transport:종료")
	}
}
