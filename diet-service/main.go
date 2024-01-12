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
	"os"
	"time"

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
	}
	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
	}

	svc := service.NewDietPresetService(database)

	savePresetEndpoint := endpoint.SavePresetEndpoint(svc)
	getPresetsEndpoint := endpoint.GetPresetsEndpoint(svc)
	removePresetsEndpoint := endpoint.RemovePresetEndpoint(svc)

	router := gin.Default()
	rateLimiter := util.NewRateLimiter(rate.Every(1*time.Minute), 100)
	router.Use(rateLimiter.Middleware())

	router.POST("/save-preset", transport.SavePresetHandler(savePresetEndpoint))
	router.GET("/get-preset", transport.GetPresetsHandler(getPresetsEndpoint))
	router.POST("/remove-preset/:id", transport.RemovePresetHandler(removePresetsEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
