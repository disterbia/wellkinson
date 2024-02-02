// /face-service/main.go
package main

import (
	"face-service/db"
	_ "face-service/docs"
	"face-service/endpoint"
	"face-service/service"
	"face-service/transport"
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
		log.Fatal("Error loading .env file")
		return
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}

	svc := service.NewFaceService(database)

	savefaceScoresEndpoint := endpoint.SaveScoresEndpoint(svc)
	getfaceScoresEndpoint := endpoint.GetScoresEndpoint(svc)
	getfaceExamsEndpoint := endpoint.GetFaceExamsEndpoint(svc)

	router := gin.Default()
	router.POST("/save-faces", transport.SaveScoresHandler(savefaceScoresEndpoint))
	router.GET("/get-faces", transport.GetScoresHandler(getfaceScoresEndpoint))
	router.GET("/get-face-exams", transport.GetFaceExamsHandler(getfaceExamsEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
