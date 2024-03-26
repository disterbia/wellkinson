// /exercise-service/main.go
package main

import (
	"exercise-service/db"
	_ "exercise-service/docs"
	"exercise-service/endpoint"
	"exercise-service/service"
	"exercise-service/transport"
	"log"
	"os"

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

	svc := service.NewExerciseService(database, conn)

	saveExerciseEndpoint := endpoint.SaveExerciseEndpoint(svc)
	getExercisesEndpoint := endpoint.GetExercisesEndpoint(svc)
	removeExercisesEndpoint := endpoint.RemoveExercisesEndpoint(svc)
	doExerciseEndpoint := endpoint.DoExerciseEndpoint(svc)
	getProjectsEndpoint := endpoint.GetProjectsEndpoint(svc)
	getVideosEndpoint := endpoint.GetVideosEndpoint(svc)

	router := gin.Default()
	router.POST("/save-exercise", transport.SaveExerciseHandler(saveExerciseEndpoint))
	router.POST("/remove-exercises", transport.RemoveExercisesHandler(removeExercisesEndpoint))
	router.POST("/do-exercise", transport.DoExerciseHandler(doExerciseEndpoint))
	router.GET("/get-exercises", transport.GetExercisesHandler(getExercisesEndpoint))
	router.GET("/get-projects", transport.GetProjectsHandler(getProjectsEndpoint))
	router.GET("/get-videos", transport.GetVideosHandler(getVideosEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44404")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
