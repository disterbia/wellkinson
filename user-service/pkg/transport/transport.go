// /user-service/pkg/transport/transport.go

package transport

import (
	"common/model"
	"common/util"
	"net/http"

	kitEndpoint "github.com/go-kit/kit/endpoint"

	"github.com/gin-gonic/gin"
)

// @Summary 구글로그인
// @Tags 로그인
// @Description 구글로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "요청 DTO - idToken,기본값 데이터"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /google-login [post]
func GoogleLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := loginEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.LoginResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 카카오로그인
// @Tags 로그인
// @Description 카카오로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "요청 DTO - dToken,기본값 데이터"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /kakao-login [post]
func KakaoLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := loginEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.LoginResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 로그인
// @Summary 자동로그인
// @Description 최초 로그인 이후 앱 실행시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Security jwt
// @Router /auto-login [post]
func AutoLoginHandler(autoLoginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 토큰 검증 및 처리

		_, email, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := autoLoginEndpoint(c.Request.Context(), email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.LoginResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 유저 데이터 변경
// @Tags 회원상태 변경(본인)
// @Description 유저 상태영구변경시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body User true "요청 DTO - 업데이트 할 데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /set-user [post]
func SetUserHandler(setUserEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 토큰 검증 및 처리
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var req model.User

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != req.Id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 id"})
			return
		}

		response, err := setUserEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 유저 조회
// @Tags 회원조회(본인)
// @Description 내 정보 조회시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} User "성공시 유저 객체 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-user [post]
func GetUserHandler(getUserEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 토큰 검증 및 처리
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := getUserEndpoint(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(model.User)
		c.JSON(http.StatusOK, resp)
	}
}
