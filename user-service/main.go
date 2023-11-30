// /user-service/main.go

package main

import (
	"log"
	"os"
	"user-service/pkg/db"
	"user-service/pkg/endpoint"
	"user-service/pkg/service"
	"user-service/pkg/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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

	gsvc := service.NewGoogleLoginService(database)
	googleLoginEndpoint := endpoint.MakeLoginEndpoint(gsvc)

	ksvc := service.NewKakaoLoginService(database)
	kakoLoginEndpoint := endpoint.MakeLoginEndpoint(ksvc)

	autoLoginSvc := service.NewAutoLoginService(database)
	autoLoginEndpoint := endpoint.MakeAutoLoginEndpoint(autoLoginSvc)

	setUserSvc := service.NewSetUserService(database)
	setUserEndpoint := endpoint.MakeSetUserEndpoint(setUserSvc)

	getUserSvc := service.NewGetUserService(database)
	getUserEndpoint := endpoint.MakeGetUserEndpoint(getUserSvc)

	router := gin.Default()

	router.POST("/google-login", transport.GoogleLoginHandler(googleLoginEndpoint))
	router.POST("/kakao-login", transport.KakaoLoginHandler(kakoLoginEndpoint))
	router.POST("/auto-login", transport.AutoLoginHandler(autoLoginEndpoint))
	router.POST("/set-user", transport.SetUserHandler(setUserEndpoint))
	router.POST("/get-user", transport.GetUserHandler(getUserEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
