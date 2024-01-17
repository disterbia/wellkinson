// /diet-service/main.go
package main

import (
	"common/util"
	_ "diet-service/docs"
	"diet-service/endpoint"
	"diet-service/service"
	"diet-service/transport"
	"fcm-service/db"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
		return
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
		log.Println("Database connection error:", err)
		return
	}

	s3svc := s3.New(s3sess)
	svc := service.NewDietService(database, s3svc, bucket, bucketUrl)

	userRateLimiter := util.NewUserRateLimiter()

	savePresetEndpoint := endpoint.SavePresetEndpoint(svc)
	getPresetsEndpoint := endpoint.GetPresetsEndpoint(svc)
	removePresetsEndpoint := endpoint.RemovePresetEndpoint(svc)
	saveDietEndpoint := endpoint.SaveDietEndpoint(svc)
	getDietsEndpoint := endpoint.GetDietsEndpoint(svc)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		id, _, err := util.VerifyJWT(c)
		if err != nil {
			// ip별 Rate Limiter 가 들어가야함.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 사용자별 Rate Limiter 확인
		userLimiter := userRateLimiter.GetUserLimiter(id, rate.Every(1*time.Minute), 10)

		if !userLimiter.Limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	})

	router.POST("/save-preset", transport.SavePresetHandler(savePresetEndpoint))
	router.POST("/remove-preset/:id", transport.RemovePresetHandler(removePresetsEndpoint))
	router.POST("/save-diet", transport.SaveDietHandler(saveDietEndpoint))

	router.GET("/get-presets", transport.GetPresetsHandler(getPresetsEndpoint))
	router.GET("/get-diets", transport.GetDietsHandler(getDietsEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
