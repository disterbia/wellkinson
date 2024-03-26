// /admin-video-service/main.go
package main

import (
	"admin-video-service/db"
	_ "admin-video-service/docs"
	"admin-video-service/endpoint"
	"admin-video-service/service"
	"admin-video-service/transport"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	svc := service.NewAdminVideoService(database)

	getVimeoLevel1sEndpoint := endpoint.GetVimeoLevel1sEndpoint(svc)
	getVimeoLevel2sEndpoint := endpoint.GetVimeoLevel2sEndpoint(svc)
	saveEndpoint := endpoint.SaveEndpoint(svc)

	router := gin.Default()
	config := cors.Config{
		AllowAllOrigins: true, // 모든 출처 허용
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}
	router.Use(cors.New(config))

	router.GET("/get-items", transport.GetVimeoLevel1sHandler(getVimeoLevel1sEndpoint))
	router.GET("/get-videos/:id", transport.GetVimeoLevel2sHandler(getVimeoLevel2sEndpoint))
	router.POST("/save-videos", transport.SaveHandler(saveEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":44400")
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
