package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.Printf(
			"[%s] %s | %s | %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
		)
	}
}
