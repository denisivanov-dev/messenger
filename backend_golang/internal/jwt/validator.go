package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("U0QfbakniIkRyfbD4Mvg46Q3PQrcUvDoV2CMF_OlYX4EfE6PEUL7xX7xD7J8K0A5XWSAEspRWrdR4s9l_j1HHA") 

type Claims struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ParseToken(tokenStr string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	return claims.Sub, claims.Username, nil
}