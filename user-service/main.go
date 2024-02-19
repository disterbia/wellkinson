// /user-service/main.go

package main

import (
	"log"
	"os"
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

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

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
	autoLoginEndpoint := endpoint.MakeAutoLoginEndpoint(usvc)
	setUserEndpoint := endpoint.MakeSetUserEndpoint(usvc)
	getUserEndpoint := endpoint.MakeGetUserEndpoint(usvc)
	getMainServicesEndpoint := endpoint.GetMainServicesEndpoint(usvc)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/admin-login", transport.AdminLoginHandler(adminLoginEndpoint))
	router.POST("/google-login", transport.GoogleLoginHandler(googleLoginEndpoint))
	router.POST("/kakao-login", transport.KakaoLoginHandler(kakoLoginEndpoint))
	router.POST("/auto-login", transport.AutoLoginHandler(autoLoginEndpoint))
	router.POST("/set-user", transport.SetUserHandler(setUserEndpoint))
	router.GET("/get-user", transport.GetUserHandler(getUserEndpoint))
	router.GET("/get-services", transport.GetMainServicesHandeler(getMainServicesEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44409")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
