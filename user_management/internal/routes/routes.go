package routes

import (
	"enterdev.com.vn/user_management/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

// Khi tạo các route khác mà có func RegisterRoutes
// thì có nghĩa là implement interface route
func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimitingMiddleware())
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
