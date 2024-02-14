// /vocal-service/main.go
package main

import (
	"log"
	"os"
	"vocal-service/db"
	_ "vocal-service/docs"
	"vocal-service/endpoint"
	"vocal-service/service"
	"vocal-service/transport"

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
	getfaceExerciseEndPoint := endpoint.GetFaceExercisesEndpoint(svc)

	router := gin.Default()
	router.POST("/save-faces", transport.SaveScoresHandler(savefaceScoresEndpoint))
	router.GET("/get-face-scores", transport.GetScoresHandler(getfaceScoresEndpoint))
	router.GET("/get-face-exams", transport.GetFaceExamsHandler(getfaceExamsEndpoint))
	router.GET("/get-face-exercises", transport.GetFaceExercisesHandler(getfaceExerciseEndPoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44444")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
