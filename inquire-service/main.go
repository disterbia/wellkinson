// /inquire-service/main.go

package main

import (
	"inquire-service/db"
	"inquire-service/endpoint"
	"inquire-service/service"
	"inquire-service/transport"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "inquire-service/docs"

	swaggerFiles "github.com/swaggo/files"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
	}

	// gRPC 클라이언트 연결 생성
	conn, err := grpc.Dial("email:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	inquireSvc := service.NewInquireService(database, conn)
	answerEndpoint := endpoint.AnswerEndpoint(inquireSvc)
	sendEndpoint := endpoint.SendEndpoint(inquireSvc)
	getEndpoint := endpoint.GetEndpoint(inquireSvc)
	allEndpoint := endpoint.GetAllEndpoint(inquireSvc)
	removeInquireEndpoint := endpoint.RemoveInquireEndpoint(inquireSvc)
	reomoveReplyEndpoint := endpoint.RemoveReplyEndpoint(inquireSvc)

	router := gin.Default()
	config := cors.Config{
		AllowAllOrigins: true, // 모든 출처 허용
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}
	router.Use(cors.New(config))

	router.POST("/inquire-reply", transport.AnswerHandler(answerEndpoint))
	router.POST("/send-inquire", transport.SendHandler(sendEndpoint))
	router.POST("/remove-inquire/:id", transport.RemoveInquireHandler(removeInquireEndpoint))
	router.POST("/remove-reply/:id", transport.RemoveReplyHandler(reomoveReplyEndpoint))
	router.GET("/get-inquires", transport.GetHandler(getEndpoint))
	router.GET("/all-inquires", transport.GetAllHandler(allEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44406")

}
