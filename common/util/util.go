package util

import (
	"common/model"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT secret key
var jwtSecretKey = []byte("adapfit_mark")

type LoginService interface {
	Login(token string, user model.User) (string, error)
}

func VerifyJWT(c *gin.Context) (int, string, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return 0, "", errors.New("authorization header is required")
	}

	// 'Bearer ' 접두사 제거
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	id := int((*claims)["id"].(float64))
	email := (*claims)["email"].(string)
	if email == "" || id == 0 {
		return 0, "", errors.New("id or email not found in token")
	}
	return id, email, nil
}

func GenerateJWT(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 1주일 유효 기간
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseCustomDateTime(input string) (time.Time, error) {
	// "YYYY-MM-DD HH:MM:SS" 형식 파싱
	if t, err := time.Parse("2006-01-02 15:04:05", input); err == nil {
		return t, nil
	}
	// "YYYY-MM-DD" 형식 파싱
	if t, err := time.Parse("2006-01-02", input); err == nil {
		return t, nil
	}
	// "YYYY/MM/DD HH-MM-SS" 형식 파싱
	if t, err := time.Parse("2006/01/02 15-04-05", input); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("invalid date format")
}

func ValidateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.New("invalid date format, should be YYYY-MM-DD")
	}
	return nil
}

func ValidateTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return errors.New("invalid time format, should be HH:MM")
	}
	return nil
}
