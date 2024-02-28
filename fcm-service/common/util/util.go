package util

import (
	"encoding/json"
	"errors"
	"fcm-service/common/model"
	"fmt"
	"log"
	"regexp"
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

func VerifyJWT(c *gin.Context) (uint, string, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
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

	id := uint((*claims)["id"].(float64))
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
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), // 한달 유효 기간
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

func ValidatePhoneNumber(phone string) error {
	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil || !matched {
		return errors.New("invalid phone format, should be 01000000000")
	}
	return nil
}

func CopyStruct(input interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, output)
	if err != nil {
		return err
	}

	return nil
}

func DecodeJwt(tokenString string) string {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}

	// MapClaims 타입으로 주장(claims) 확인
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 'iss' 주장 확인
		if iss, ok := claims["iss"].(string); ok {
			fmt.Println("Issuer (iss):", iss)
			return iss
		} else {
			fmt.Println("'iss' 주장이 없습니다.")
			return ""
		}
	} else {
		log.Fatal("주장을 MapClaims로 변환할 수 없습니다.")
		return ""
	}
}
