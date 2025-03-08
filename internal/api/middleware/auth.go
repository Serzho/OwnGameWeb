package middleware

import (
	"OwnGameWeb/config"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type UserClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectToLogin := func() {
			fmt.Println("Redirect to login page")
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		}

		tokenString, err := c.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Println("No cookies at request!")
			redirectToLogin()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return cfg.Global.SecretPhrase, nil
		})

		if err != nil || !token.Valid {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				fmt.Println("Invalid token format")
			case errors.Is(err, jwt.ErrTokenExpired):
				fmt.Println("Token expired")
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				fmt.Println("Token not active yet")
			default:
				fmt.Printf("Validation error: %v\n", err)
			}
			redirectToLogin()
			return
		}

		claims, ok := token.Claims.(*UserClaims)

		if !ok {
			fmt.Println("Invalid token claims")
			redirectToLogin()
			return
		}

		fmt.Printf("Claims: %+v\n", claims)

		c.Set("userId", claims.Id)
		c.Next()
	}
}
