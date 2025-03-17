package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
)

type UserClaims struct {
	Id     int `json:"id"`
	GameId int `json:"gameId"`
	jwt.RegisteredClaims
}

func JwtParse(tokenString string, secretPhrase string) (*UserClaims, error) {
	JwtParseErr := fmt.Errorf("jwt parse error")
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("unexpected signing method", "alg", token.Header["alg"])
			return nil, JwtParseErr
		}
		return []byte(secretPhrase), nil
	})

	if err != nil || !token.Valid {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			slog.Warn("Invalid token format")
		case errors.Is(err, jwt.ErrTokenExpired):
			slog.Warn("Token expired")
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			slog.Warn("Token not active yet")
		default:
			slog.Error("Validation error: ", "err", err)
		}

		return nil, JwtParseErr
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		slog.Warn("Invalid token claims", "claims", token.Claims)
		return nil, JwtParseErr
	}

	slog.Info("UserClaims", "claims", claims)
	return claims, nil
}

func JwtCreate(userId int, gameId int, secretPhrase string) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256, UserClaims{Id: userId, GameId: gameId, RegisteredClaims: jwt.RegisteredClaims{}},
	)
	slog.Info("Creating jwt for claims", "claims", claims)
	tokenString, err := claims.SignedString([]byte(secretPhrase))

	if err != nil {
		return "", errors.New("error create token")
	}

	return tokenString, nil
}
