package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	expHours, err := strconv.Atoi(os.Getenv("JWT_EXP_HOURS"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseJWT(tokensStr string) (string, error) {
	token, err := jwt.Parse(tokensStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token string: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid subject claim")
		}
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
