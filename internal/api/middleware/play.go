package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Play() gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectToMain := func() {
			slog.Warn("redirecting to main page")
			c.Redirect(http.StatusTemporaryRedirect, "/main")
			c.Abort()
		}

		value, exists := c.Get("gameID")

		if !exists {
			slog.Warn("game id not found in context")
			redirectToMain()
			return
		}

		_, ok := value.(int)
		if !ok {
			slog.Warn("game id has incorrect type")
			redirectToMain()
			return
		}

		c.Next()
	}
}
