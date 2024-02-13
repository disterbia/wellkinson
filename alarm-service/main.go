// /alarm-service/main.go

package main

import (
	"alarm-service/db"
	"alarm-service/endpoint"
	"alarm-service/service"
	"alarm-service/transport"
	"log"
	"net"
	"os"

	_ "alarm-service/docs"
	pb "alarm-service/proto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"

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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	alarmServer := &service.AlarmServer{Db: database}
	pb.RegisterAlarmServiceServer(grpcServer, alarmServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	alarmSvc := service.NewAlarmService(database)

	saveAlarmEndpoint := endpoint.SaveAlarmEndpoint(alarmSvc)
	removeAlarmsEndpoint := endpoint.RemoveAlarmEndpoint(alarmSvc)
	getAlarmsEndpoint := endpoint.GetEndpoint(alarmSvc)

	router := gin.Default()

	router.POST("/save-alarm", transport.SaveAlarmHandler(saveAlarmEndpoint))
	router.POST("/remove-alarms", transport.RemoveAlarmsHandler(removeAlarmsEndpoint))
	router.GET("/get-alarms", transport.GetHandler(getAlarmsEndpoint))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44444ç")
}
