package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"

	"github.com/gin-gonic/gin"
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

		c.Set("userID", claims.ID)

		if claims.GameID != -1 {
			c.Set("gameID", claims.GameID)
		}

		c.Next()
	}
}
