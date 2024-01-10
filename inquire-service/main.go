// /inquire-service/main.go

package main

import (
	"common/util"
	"inquire-service/db"
	"inquire-service/endpoint"
	"inquire-service/service"
	"inquire-service/transport"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "inquire-service/docs"

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

	// gRPC 클라이언트 연결 생성
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	inquireSvc := service.NewInquireService(database, conn)
	answerEndpoint := endpoint.AnswerEndpoint(inquireSvc)
	sendEndpoint := endpoint.SendEndpoint(inquireSvc)
	getEndpoint := endpoint.GetEndpoint(inquireSvc)
	allEndpoint := endpoint.GetAllEndpoint(inquireSvc)

	router := gin.Default()

	rateLimiter := util.NewRateLimiter(rate.Every(1*time.Minute), 20)
	router.Use(rateLimiter.Middleware())

	router.POST("/inquire-answer", transport.AnswerHandler(answerEndpoint))
	router.POST("/send-inquire", transport.SendHandler(sendEndpoint))
	router.GET("/get-inquires", transport.GetHandler(getEndpoint))
	router.GET("/all-inquires", transport.GetAllHandler(allEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44444")

}
