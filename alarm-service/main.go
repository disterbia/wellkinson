// /alarm-service/main.go

package main

import (
	"alarm-service/db"
	"alarm-service/endpoint"
	"alarm-service/service"
	"alarm-service/transport"
	"common/util"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

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

	rateLimiter := util.NewRateLimiter(rate.Every(1*time.Minute), 100)
	router.Use(rateLimiter.Middleware())

	router.POST("/save-alarm", transport.SaveAlarmHandler(saveAlarmEndpoint))
	router.POST("/remove-alarm/:id", transport.RemoveAlarmHandler(removeeAlarmEndpoint))
	router.GET("/get-alarms", transport.GetHandler(getAlarmEndpoint))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44444")
}
