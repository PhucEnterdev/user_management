package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 10)
		// 5 request/sec, brust : 10 (số request tối đa có thể xử lý ngay lập tức)
		// nghĩa là nếu sử dụng 1 request thì mất 1 token và mỗi giây sau đó nó sẽ cấp lại 1 token
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

/** Hàm này chạy liên tục => sử dụng goroutines*/
func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

/** Để test thì ta sử dụng apache benchmark để gửi nhiều request cùng lúc*/
// câu lệnh: ab -n 20 -c 1 -H "X-API-Key:phuccongtu" "http://localhost:9999/api/v1/news/phuc-cong-tu"
// -n : số lượng rq, c: core, h: header
func RateLimitingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)
		if ip == "" {
			ip = ctx.Request.RemoteAddr
		}
		limiter := getRateLimiter(ip)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many request",
				"msg":   "Bạn đa gửi quá nhiều request. Hãy thử lại sau",
			})
			return
		}

		ctx.Next()
	}
}
