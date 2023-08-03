package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int64
	jwt.RegisteredClaims
}

const (
	lifeTime = 1 * time.Hour
)

func GenerateToken(userID int64, key []byte) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(lifeTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func ParseToken(token string, key []byte) (int64, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if claims.Valid() != nil {
		return 0, claims.Valid()
	}
	return claims.UserID, err
}

func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
