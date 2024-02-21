// /user-service/main.go

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"user-service/db"
	"user-service/endpoint"
	"user-service/service"
	"user-service/transport"

	_ "user-service/docs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

var ipLimiters = make(map[string]*rate.Limiter)
var ipLimitersMutex sync.Mutex

func getClientIP(c *gin.Context) string {
	// X-Real-IP 헤더를 확인
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	// X-Forwarded-For 헤더를 확인
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0] // 여러 IP가 쉼표로 구분되어 있을 수 있음
	}
	// 헤더가 없는 경우 Gin의 기본 메서드 사용
	return c.ClientIP()
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)

		// IP별 리미터가 있는지 확인
		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			// 새로운 리미터 생성
			limiter = rate.NewLimiter(rate.Every(24*time.Hour/5), 5)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청 허용 여부 확인
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "요청 횟수 초과"})
			return
		}

		c.Next()
	}
}
func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")
	bucketUrl := os.Getenv("S3_BUCKET_URL")
	s3sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Println("aws connection error:", err)
	}

	s3svc := s3.New(s3sess)
	usvc := service.NewUserService(database, s3svc, bucket, bucketUrl)

	adminLoginEndpoint := endpoint.MakeAdminLoginEndpoint(usvc)
	googleLoginEndpoint := endpoint.MakeGoogleLoginEndpoint(usvc)
	kakoLoginEndpoint := endpoint.MakeKakaoLoginEndpoint(usvc)
	appleLoginEndpoint := endpoint.MakeAppleLoginEndpoint(usvc)
	autoLoginEndpoint := endpoint.MakeAutoLoginEndpoint(usvc)
	setUserEndpoint := endpoint.MakeSetUserEndpoint(usvc)
	getUserEndpoint := endpoint.MakeGetUserEndpoint(usvc)
	getversionEndpoint := endpoint.GetVersionEndpoint(usvc)
	getMainServicesEndpoint := endpoint.GetMainServicesEndpoint(usvc)
	sendCodeEndpoint := endpoint.SendCodeEndpoint(usvc)
	verifyEndpoint := endpoint.VerifyEndpoint(usvc)
	removeEndpoint := endpoint.RemoveEndpoint(usvc)
	linkEndpoint := endpoint.LinkEndpoint(usvc)

	router := gin.Default()
	router.Use(cors.Default())
	rateLimiterMiddleware := RateLimitMiddleware()

	router.POST("/admin-login", transport.AdminLoginHandler(adminLoginEndpoint))
	router.POST("/google-login", transport.GoogleLoginHandler(googleLoginEndpoint))
	router.POST("/kakao-login", transport.KakaoLoginHandler(kakoLoginEndpoint))
	router.POST("/apple-login", transport.AppleLoginHandler(appleLoginEndpoint))
	router.POST("/auto-login", transport.AutoLoginHandler(autoLoginEndpoint))
	router.POST("/set-user", transport.SetUserHandler(setUserEndpoint))
	router.POST("/send-code/:number", rateLimiterMiddleware, transport.SendCodeHandler(sendCodeEndpoint))
	router.POST("/verify-code", transport.VerifyHandler(verifyEndpoint))
	router.POST("/remove-user", transport.RemoveHandler(removeEndpoint))
	router.POST("/link-email", transport.LinkHandler(linkEndpoint))

	router.GET("/get-user", transport.GetUserHandler(getUserEndpoint))
	router.GET("/get-version", transport.GetVersionHandeler(getversionEndpoint))
	router.GET("/get-services", transport.GetMainServicesHandeler(getMainServicesEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44409")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
