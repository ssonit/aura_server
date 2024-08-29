package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret []byte, userID string, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    exp,
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodedToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func GenerateExpTime(d time.Duration) int64 {
	return time.Now().Add(d).Unix()
}
