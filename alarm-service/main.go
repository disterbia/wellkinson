// /alarm-service/main.go

package main

import (
	"alarm-service/db"
	"alarm-service/endpoint"
	"alarm-service/service"
	"alarm-service/transport"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "alarm-service/docs"

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

	alarmSvc := service.NewAlarmService(database)

	saveAlarmEndpoint := endpoint.SaveAlarmEndpoint(alarmSvc)
	removeeAlarmEndpoint := endpoint.RemoveAlarmEndpoint(alarmSvc)
	getAlarmEndpoint := endpoint.GetEndpoint(alarmSvc)

	router := gin.Default()

	router.POST("/save-alarm", transport.SaveAlarmHandler(saveAlarmEndpoint))
	router.POST("/remove-alarm/:id", transport.RemoveAlarmHandler(removeeAlarmEndpoint))
	router.GET("/get-alarms", transport.GetHandler(getAlarmEndpoint))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44444")
}
