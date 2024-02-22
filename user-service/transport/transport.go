// /user-service/transport/transport.go

package transport

import (
	"common/util"
	"net/http"
	"user-service/dto"

	kitEndpoint "github.com/go-kit/kit/endpoint"

	"github.com/gin-gonic/gin"
)

// @Tags 로그인 /user
// @Summary 관리자 로그인
// @Description 관리자 로그인시 호출
// @Accept  json
// @Produce  json
// @Param email body string true "email"
// @Param password body string true "password"
// @Success 200 {object} dto.SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Security jwt
// @Router /admin-login [post]
func AdminLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := requestData["email"].(string)
		password := requestData["password"].(string)
		response, err := loginEndpoint(c.Request.Context(), map[string]interface{}{
			"email":    email,
			"password": password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.LoginResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 로그인 /user
// @Summary 구글로그인
// @Description 구글로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "요청 DTO - idToken,기본값 데이터 user_type: 0:해당없음 1:파킨슨 환자 2:보호자"
// @Success 200 {object} dto.SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 번호인증 필요"
// @Router /google-login [post]
func GoogleLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := loginEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.LoginResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 로그인 /user
// @Summary 카카오로그인
// @Description 카카오로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "요청 DTO - idToken,기본값 데이터 user_type: 0:해당없음 1:파킨슨 환자 2:보호자"
// @Success 200 {object} dto.SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 번호인증 필요"
// @Router /kakao-login [post]
func KakaoLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := loginEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.LoginResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 로그인 /user
// @Summary 애플로그인
// @Description 애플로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "요청 DTO - idToken,기본값 데이터 user_type: 0:해당없음 1:파킨슨 환자 2:보호자"
// @Success 200 {object} dto.SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 번호인증 필요"
// @Router /apple-login [post]
func AppleLoginHandler(loginEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := loginEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.LoginResponse)
		c.JSON(http.StatusOK, resp)

	}
}

// @Tags 로그인 /user
// @Summary 자동로그인
// @Description 최초 로그인 이후 앱 실행시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} dto.SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
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

		resp := response.(dto.LoginResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 회원상태 변경(본인)  /user
// @Summary 유저 데이터 변경
// @Description 유저 상태영구변경시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.UserRequest true "요청 DTO - 업데이트 할 데이터/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 번호인증 필요"
// @Router /set-user [post]
func SetUserHandler(setUserEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 토큰 검증 및 처리
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var req dto.UserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.Id = id

		response, err := setUserEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 회원조회(본인)  /user
// @Summary 유저 조회
// @Description 내 정보 조회시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} dto.UserResponse "성공시 유저 객체 반환/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자 sns_type- 0:카카오 1:구글 2:애플"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
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

		resp := response.(dto.UserResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 공통 /user
// @Summary 전체 서비스 목록 조회
// @Description 이용하고 싶은 서비스 목록 조회시 호출
// @Accept  json
// @Produce  json
// @Success 200 {object} []dto.MainServiceResponse "서비스 목록"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-services [get]
func GetMainServicesHandeler(getUserEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		response, err := getUserEndpoint(c.Request.Context(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.([]dto.MainServiceResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 인증번호 /user
// @Summary 인증번호 발송
// @Description 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param number path string true "휴대번호"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /send-code/{number} [post]
func SendCodeHandler(sendEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		number := c.Param("number")

		response, err := sendEndpoint(c.Request.Context(), number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 인증번호 /user
// @Summary 번호 인증
// @Description 인증번호 입력 후 호출
// @Accept  json
// @Produce  json
// @Param request body dto.VerifyRequest true "요청 DTO"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /verify-code [post]
func VerifyHandler(verifyEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req dto.VerifyRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := verifyEndpoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 회원탈퇴 /user
// @Summary 회원탈퇴
// @Description 회원탈퇴시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remvoe-user [post]
func RemoveHandler(removeEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 토큰 검증 및 처리
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := removeEndpoint(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.UserResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 계정 연동 /user
// @Summary 계정 연동
// @Description 계정 연동시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body dto.LinkRequest true "요청 DTO"
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /link-email [post]
func LinkHandler(endPoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var req dto.LinkRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = id
		response, err := endPoint(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.BasicResponse)
		c.JSON(http.StatusOK, resp)
	}
}

// @Tags 공통 /user
// @Summary 최신버전 조회
// @Description 최신버전 조회시 호출
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.AppVersionResponse "최신 버전 정보"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-version [get]
func GetVersionHandeler(getUserEndpoint kitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		response, err := getUserEndpoint(c.Request.Context(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := response.(dto.AppVersionResponse)
		c.JSON(http.StatusOK, resp)
	}
}
