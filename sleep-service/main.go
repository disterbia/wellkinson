// /sleep-service/main.go
package main

import (
	"log"
	"os"
	"sleep-service/db"
	_ "sleep-service/docs"
	"sleep-service/endpoint"
	"sleep-service/service"
	"sleep-service/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	// gRPC 클라이언트 연결 생성
	conn, err := grpc.Dial("alarm:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	svc := service.NewSleepService(database, conn)

	saveAlarmsEndpoint := endpoint.SaveSleepAlarmEndpoint(svc)
	getSleepAlarmsEndpoint := endpoint.GetSleepAlarmsEndpoint(svc)
	removeSleepAlarmsEndpoint := endpoint.RemoveSleepAlarmsEndpoint(svc)
	getSleepTimesEndpoint := endpoint.GetSleepTimesEndpoint(svc)
	removeSleepTimesEndpoint := endpoint.RemoveSleepTimeEndpoint(svc)
	saveSleepTimesEndpoint := endpoint.SaveSleepTimeEndpoint(svc)

	router := gin.Default()
	router.POST("/save-sleep-alarm", transport.SaveSleepHandler(saveAlarmsEndpoint))
	router.POST("/remove-sleep-alarms", transport.RemoveSleepAlarmsHandler(removeSleepAlarmsEndpoint))
	router.POST("/save-sleep-time", transport.SaveSleepTimeHandler(saveSleepTimesEndpoint))
	router.POST("/remove-sleep-time/:id", transport.RemoveSleepTimeHandler(removeSleepTimesEndpoint))
	router.GET("/get-sleep-alarms", transport.GetSleepAlarmsHandler(getSleepAlarmsEndpoint))
	router.GET("/get-sleep-times", transport.GetSleepTimesHandler(getSleepTimesEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44408")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
