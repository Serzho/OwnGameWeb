package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func JwtParse(tokenString string, secretPhrase string) (*UserClaims, error) {
	JwtParseErr := fmt.Errorf("jwt parse error")
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("unexpected signing method: %v\n", token.Header["alg"])
			return nil, JwtParseErr
		}
		return secretPhrase, nil
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

		return nil, JwtParseErr
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		fmt.Println("Invalid token claims")
		return nil, JwtParseErr
	}

	return claims, nil
}
