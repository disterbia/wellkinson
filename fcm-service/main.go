// /fcm-service/main.go
package main

import (
	"fcm-service/db"
	"fcm-service/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	service.StartCentralCronScheduler(database)

	router := gin.Default()
	router.Run(":6767")
}
