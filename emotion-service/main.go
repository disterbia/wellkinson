// /emotion-service/main.go
package main

import (
	"emotion-service/db"
	_ "emotion-service/docs"
	"emotion-service/endpoint"
	"emotion-service/service"
	"emotion-service/transport"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	svc := service.NewEmotionService(database)

	saveEmotionEndpoint := endpoint.SaveEmotionEndpoint(svc)
	getEmotionsEndpoint := endpoint.GetEmotionsEndpoint(svc)
	removeEmotionsEndpoint := endpoint.RemoveEmotionsEndpoint(svc)

	router := gin.Default()
	router.POST("/save-emotion", transport.SaveEmotionHandler(saveEmotionEndpoint))
	router.POST("/remove-emotions", transport.RemoveEmotionsHandler(removeEmotionsEndpoint))
	router.GET("/get-emotions", transport.GetEmotionsHandler(getEmotionsEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44403")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
