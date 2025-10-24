package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	User string `json:"user"`
}

func GenerateUserToken(user, secret string, expiration int) (string, error) {
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + int64(expiration), // Add expiration in seconds
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func ValidateToken(tokenString, secret string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
