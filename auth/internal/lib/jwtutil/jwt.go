package jwtutil

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWTSecret []byte

func init() {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Fatalf("NO JWT SECRET KEY PROVIDED")
	}

	JWTSecret = []byte(key)
}

func GenerateJWT(userId int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(JWTSecret)

	if err != nil {
		return "", fmt.Errorf("error during jwt token generation: %s", err.Error())
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return JWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
