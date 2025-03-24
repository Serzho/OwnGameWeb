package utils

import (
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID     int `json:"userId"`
	GameID int `json:"gameId"`
	jwt.RegisteredClaims
}

func JwtParse(tokenString string, secretPhrase string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("unexpected signing method", "alg", token.Header["alg"])

			return nil, ErrJWTParse
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

		return nil, ErrJWTParse
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		slog.Warn("Invalid token claims", "claims", token.Claims)

		return nil, ErrJWTParse
	}

	slog.Info("UserClaims", "claims", claims)

	return claims, nil
}

func JwtCreate(userID int, gameID int, secretPhrase string) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256, UserClaims{ID: userID, GameID: gameID, RegisteredClaims: jwt.RegisteredClaims{}},
	)
	slog.Info("Creating jwt for claims", "claims", claims)

	tokenString, err := claims.SignedString([]byte(secretPhrase))
	if err != nil {
		return "", ErrCreatingToken
	}

	return tokenString, nil
}
