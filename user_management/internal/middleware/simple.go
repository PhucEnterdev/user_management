package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Trước khi bắt đầu vào handler
		log.Println("Start func - Check from middleware")

		ctx.Next()

		// Sau khi handler xử lý
		log.Println("End func - Check from middleware")

	}
}
