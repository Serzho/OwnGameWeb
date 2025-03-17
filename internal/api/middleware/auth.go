package middleware

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectToLogin := func() {
			slog.Warn("redirecting to login page")
			c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
			c.Abort()
		}

		tokenString, err := c.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			slog.Warn("No token cookie")
			redirectToLogin()
			return
		}

		claims, err := utils.JwtParse(tokenString, cfg.Global.SecretPhrase)

		if err != nil {
			slog.Warn("Invalid token", "err", err)
			redirectToLogin()
			return
		}

		slog.Info("Setting claims", "claims", claims)

		c.Set("userId", claims.Id)
		if claims.GameId != -1 {
			c.Set("gameId", claims.GameId)
		}
		c.Next()
	}
}
