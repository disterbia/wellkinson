// /gateway/main.go

package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IP별 레이트 리미터를 저장할 맵과 이를 동기화하기 위한 뮤텍스
var (
	ips = make(map[string]*rate.Limiter)
	mu  sync.RWMutex
)

// 특정 IP 주소에 대한 레이트 리미터를 반환
func GetLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := ips[ip]
	if !exists {
		limiter = rate.NewLimiter(20, 20) // 레이트 리미팅 설정 조정
		ips[ip] = limiter
	}

	return limiter
}
func getClientIP(c *gin.Context) string {
	// X-Real-IP 헤더를 확인
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	// X-Forwarded-For 헤더를 확인
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0] // 여러 IP가 쉼표로 구분되어 있을 수 있음
	}
	// 헤더가 없는 경우 Gin의 기본 메서드 사용
	return c.ClientIP()
}

// IP 주소별로 레이트 리미팅을 적용
func IPRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Swagger UI에 대한 요청은 레이트 리미팅에서 제외
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}

		ip := getClientIP(c)
		limiter := GetLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "요청 수가 너무 많습니다",
			})
			return
		}

		c.Next()
	}
}
func main() {
	router := gin.Default()
	router.Use(IPRateLimitMiddleware())

	//서비스로의 리버스 프록시 설정
	adminServiceURL, _ := url.Parse("http://admin:44400")
	adminProxy := httputil.NewSingleHostReverseProxy(adminServiceURL)
	router.Any("/admin/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		adminProxy.ServeHTTP(c.Writer, c.Request)
	})

	alarmServiceURL, _ := url.Parse("http://alarm:44401")
	alarmProxy := httputil.NewSingleHostReverseProxy(alarmServiceURL)
	router.Any("/alarm/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		alarmProxy.ServeHTTP(c.Writer, c.Request)
	})

	dietServiceURL, _ := url.Parse("http://diet:44402")
	dietProxy := httputil.NewSingleHostReverseProxy(dietServiceURL)
	router.Any("/diet/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		dietProxy.ServeHTTP(c.Writer, c.Request)
	})

	emotionServiceURL, _ := url.Parse("http://emotion:44403")
	emotionProxy := httputil.NewSingleHostReverseProxy(emotionServiceURL)
	router.Any("/emotion/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		emotionProxy.ServeHTTP(c.Writer, c.Request)
	})

	exerciseServiceURL, _ := url.Parse("http://exercise:44404")
	exerciseProxy := httputil.NewSingleHostReverseProxy(exerciseServiceURL)
	router.Any("/exercise/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		exerciseProxy.ServeHTTP(c.Writer, c.Request)
	})

	faceServiceURL, _ := url.Parse("http://face:44405")
	faceProxy := httputil.NewSingleHostReverseProxy(faceServiceURL)
	router.Any("/face/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		faceProxy.ServeHTTP(c.Writer, c.Request)
	})

	inquireServiceURL, _ := url.Parse("http://inquire:44406")
	inquireProxy := httputil.NewSingleHostReverseProxy(inquireServiceURL)
	router.Any("/inquire/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		inquireProxy.ServeHTTP(c.Writer, c.Request)
	})

	medicineServiceURL, _ := url.Parse("http://medicine:44407")
	medicineProxy := httputil.NewSingleHostReverseProxy(medicineServiceURL)
	router.Any("/medicine/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		medicineProxy.ServeHTTP(c.Writer, c.Request)
	})

	sleepServiceURL, _ := url.Parse("http://sleep:44408")
	sleepProxy := httputil.NewSingleHostReverseProxy(sleepServiceURL)
	router.Any("/sleep/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		sleepProxy.ServeHTTP(c.Writer, c.Request)
	})

	userServiceURL, _ := url.Parse("http://user:44409")
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	router.Any("/user/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path") // '/user' 접두사 제거
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

	vocalServiceURL, _ := url.Parse("http://vocal:44410")
	vocalProxy := httputil.NewSingleHostReverseProxy(vocalServiceURL)
	router.Any("/vocal/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		vocalProxy.ServeHTTP(c.Writer, c.Request)
	})

	setupSwaggerUIProxy(router, "/admin-service/swagger/*proxyPath", "http://admin:44400/swagger/")
	setupSwaggerUIProxy(router, "/alarm-service/swagger/*proxyPath", "http://alarm:44401/swagger/")
	setupSwaggerUIProxy(router, "/diet-service/swagger/*proxyPath", "http://diet:44402/swagger/")
	setupSwaggerUIProxy(router, "/emotion-service/swagger/*proxyPath", "http://emotion:44403/swagger/")
	setupSwaggerUIProxy(router, "/exercise-service/swagger/*proxyPath", "http://exercise:44404/swagger/")
	setupSwaggerUIProxy(router, "/face-service/swagger/*proxyPath", "http://face:44405/swagger/")
	setupSwaggerUIProxy(router, "/inquire-service/swagger/*proxyPath", "http://inquire:44406/swagger/")
	setupSwaggerUIProxy(router, "/medicine-service/swagger/*proxyPath", "http://medicine:44407/swagger/")
	setupSwaggerUIProxy(router, "/sleep-service/swagger/*proxyPath", "http://sleep:44408/swagger/")
	setupSwaggerUIProxy(router, "/user-service/swagger/*proxyPath", "http://user:44409/swagger/")
	setupSwaggerUIProxy(router, "/vocal-service/swagger/*proxyPath", "http://vocal:44410/swagger/")

	// API 게이트웨이 서버 시작
	router.Run(":50000")
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
