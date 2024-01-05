// /inquire-service/main.go

package main

import (
	"common/util"
	"inquire-service/db"
	"inquire-service/endpoint"
	"inquire-service/service"
	"inquire-service/transport"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

	_ "inquire-service/docs"

	swaggerFiles "github.com/swaggo/files"
)

type User struct {
	ID                   int    `gorm:"primaryKey;autoIncrement"`
	Birthday             string `gorm:"size:40;default:''"`
	DeviceID             string `gorm:"size:40;default:''"`
	Gender               bool
	FCMToken             string    `gorm:"size:255;default:''"`
	IsFirst              bool      `gorm:"default:false"`
	Name                 string    `gorm:"size:40;default:''"`
	PhoneNum             string    `gorm:"size:40;default:''"`
	UseAutoLogin         bool      `gorm:"default:false"`
	UsePrivacyProtection bool      `gorm:"default:false"`
	UseSleepTracking     bool      `gorm:"default:false"`
	UserType             string    `gorm:"size:40;default:''"`
	Email                string    `gorm:"size:40;default:''"`
	Created              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	Updated              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

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

	inquireSvc := service.NewInquireService(database)
	answerEndpoint := endpoint.AnswerEndpoint(inquireSvc)
	sendEndpoint := endpoint.SendEndpoint(inquireSvc)
	getEndpoint := endpoint.GetEndpoint(inquireSvc)

	router := gin.Default()

	rateLimiter := util.NewRateLimiter(rate.Every(1*time.Minute), 10)
	router.Use(rateLimiter.Middleware())

	router.POST("/inquire-answer", transport.AnswerHandler(answerEndpoint))
	router.POST("/send-inquire", transport.SendHandler(sendEndpoint))
	router.GET("/get-inquires", transport.GetHandler(getEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":44444")

}
