package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}


func GenerateToken(userId int) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),

			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)


	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}


func ValidateToken(tokenString string) (int, error) {
	secretKey := os.Getenv("JWT_SECRET")


	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})


	if err != nil {
		return 0, err
	}


	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserId, nil
	}


	return 0, jwt.ErrTokenInvalidClaims
}