// /gateway/main.go

package main

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 유저 서비스로의 리버스 프록시 설정
	userServiceURL, _ := url.Parse("http://localhost:8083")
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	router.Any("/login", func(c *gin.Context) {
		log.Printf("API Gateway: Forwarding request to user service")
		userProxy.ServeHTTP(c.Writer, c.Request)
		log.Printf("API Gateway: Request forwarded")
	})

	// fcm 서비스로의 리버스 프록시 설정
	fcmServiceURL, _ := url.Parse("http://localhost:44444")
	fcmProxy := httputil.NewSingleHostReverseProxy(fcmServiceURL)
	router.Any("/register", func(c *gin.Context) {
		log.Printf("API Gateway: Forwarding request to fcm service")
		fcmProxy.ServeHTTP(c.Writer, c.Request)
		log.Printf("API Gateway: Request forwarded")
	})

	// 알람 서비스로의 리버스 프록시 설정
	alarmrServiceURL, _ := url.Parse("http://localhost:44444")
	alarmProxy := httputil.NewSingleHostReverseProxy(alarmrServiceURL)
	router.Any("/user", func(c *gin.Context) {
		log.Printf("API Gateway: Forwarding request to alarm service ")
		alarmProxy.ServeHTTP(c.Writer, c.Request)
		log.Printf("API Gateway: Request forwarded")
	})

	// 문의 서비스로의 리버스 프록시 설정
	InquireServiceURL, _ := url.Parse("http://localhost:44444")
	inquireProxy := httputil.NewSingleHostReverseProxy(InquireServiceURL)
	router.Any("/user", func(c *gin.Context) {
		log.Printf("API Gateway: Forwarding request to inquire service ")
		inquireProxy.ServeHTTP(c.Writer, c.Request)
		log.Printf("API Gateway: Request forwarded")
	})

	setupSwaggerUIProxy(router, "/user-service/swagger/*proxyPath", "http://localhost:44444/swagger/")
	setupSwaggerUIProxy(router, "/fcm-service/swagger/*proxyPath", "http://localhost:44444/swagger/")
	setupSwaggerUIProxy(router, "/alarm-service/swagger/*proxyPath", "http://localhost:44444/swagger/")
	setupSwaggerUIProxy(router, "/inquire-service/swagger/*proxyPath", "http://localhost:44444/swagger/")
	// API 게이트웨이 서버 시작
	router.Run(":8085")
}

// Swagger 문서에 대한 리버스 프록시를 설정
func setupSwaggerUIProxy(router *gin.Engine, path string, target string) {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	router.Any(path, func(c *gin.Context) {
		// Swagger 경로 재설정
		c.Request.URL.Path = c.Param("proxyPath")
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
