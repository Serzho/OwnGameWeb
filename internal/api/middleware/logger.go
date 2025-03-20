package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()

		if c.Writer.Status() >= 500 {
			slog.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())

			return
		}

		slog.Info(
			"HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"statusCode", c.Writer.Status(),
			"latency", end.Sub(start),
			"clientIP", c.ClientIP(),
			"body", c.Request.Body,
			"params", c.Params,
		)
	}
}
