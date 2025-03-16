package middleware

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectToLogin := func() {
			fmt.Println("Redirect to login page")
			c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
			c.Abort()
		}

		tokenString, err := c.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Println("No cookies at request!")
			redirectToLogin()
			return
		}

		claims, err := utils.JwtParse(tokenString, cfg.Global.SecretPhrase)

		if err != nil {
			fmt.Println(err)
			redirectToLogin()
			return
		}

		fmt.Printf("Claims: %+v\n", claims)

		c.Set("userId", claims.Id)
		if claims.GameId != -1 {
			c.Set("gameId", claims.GameId)
		}
		c.Next()
	}
}
