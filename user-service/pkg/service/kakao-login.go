// /user-service/pkg/service/kakao-login.go

package service

import (
	"common/model"
	"common/util"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type KakaoPublicKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []KakaoPublicKey `json:"keys"`
}

type kakaoLoginService struct {
	db *gorm.DB
}

func NewKakaoLoginService(db *gorm.DB) LoginService {
	return &kakaoLoginService{db: db}
}

func (kl *kakaoLoginService) Login(idToken string, user model.User) (string, error) {
	jwks, err := getKakaoPublicKeys()
	if err != nil {
		return "", err
	}

	parsedToken, err := verifyKakaoTokenSignature(idToken, jwks)
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("email not found in token claims")
		}
		u, err := kl.findOrCreateUser(email, user)
		if err != nil {
			return "", errors.New(err.Error())
		}

		// JWT 토큰 생성
		tokenString, err := util.GenerateJWT(u)
		if err != nil {
			return "", errors.New(err.Error())
		}
		return tokenString, nil
	}
	return "", errors.New("invalid token")

}

func (kl *kakaoLoginService) findOrCreateUser(email string, user model.User) (model.User, error) {
	// 유효성 검사 수행
	if err := util.ValidateDate(user.Birthday); err != nil {
		return model.User{}, err
	}
	// 데이터베이스에서 사용자 조회 및 없으면 생성
	err := kl.db.Where(model.User{Email: email}).FirstOrCreate(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// 카카오 공개키 가져오기
func getKakaoPublicKeys() (JWKS, error) {
	resp, err := http.Get("https://kauth.kakao.com/.well-known/jwks.json")
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKS{}, err
	}
	return jwks, nil
}

// 카카오 공개키로 서명 검증
func verifyKakaoTokenSignature(token string, jwks JWKS) (*jwt.Token, error) {
	kid, err := extractKidFromToken(token)
	if err != nil {
		return nil, err
	}

	var key *rsa.PublicKey
	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
			if err != nil {
				return nil, err
			}

			n := big.NewInt(0).SetBytes(nBytes)
			e := big.NewInt(0).SetBytes(eBytes).Int64()
			key = &rsa.PublicKey{N: n, E: int(e)}
			break
		}
	}

	if key == nil {
		return nil, errors.New("appropriate public key not found")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

// ID 토큰에서 kid 추출
func extractKidFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}
	headerPart := parts[0]
	headerJson, err := base64.RawURLEncoding.DecodeString(headerPart)
	if err != nil {
		return "", err
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJson, &header); err != nil {
		return "", err
	}

	kid, ok := header["kid"].(string)
	if !ok {
		return "", errors.New("kid not found in token header")
	}
	return kid, nil
}
