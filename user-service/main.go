// /user-service/main.go

package main

import (
	"common/util"
	"log"
	"os"
	"time"
	"user-service/db"
	"user-service/endpoint"
	"user-service/service"
	"user-service/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

	_ "user-service/docs"

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

	usvc := service.NewUserService(database)

	googleLoginEndpoint := endpoint.MakeGoogleLoginEndpoint(usvc)
	kakoLoginEndpoint := endpoint.MakeKakaoLoginEndpoint(usvc)
	autoLoginEndpoint := endpoint.MakeAutoLoginEndpoint(usvc)
	setUserEndpoint := endpoint.MakeSetUserEndpoint(usvc)
	getUserEndpoint := endpoint.MakeGetUserEndpoint(usvc)

	router := gin.Default()
	rateLimiter := util.NewRateLimiter(rate.Every(1*time.Minute), 5)
	router.Use(rateLimiter.Middleware())
	router.POST("/google-login", transport.GoogleLoginHandler(googleLoginEndpoint))
	router.POST("/kakao-login", transport.KakaoLoginHandler(kakoLoginEndpoint))
	router.POST("/auto-login", transport.AutoLoginHandler(autoLoginEndpoint))
	router.POST("/set-user", transport.SetUserHandler(setUserEndpoint))
	router.GET("/get-user", transport.GetUserHandler(getUserEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
