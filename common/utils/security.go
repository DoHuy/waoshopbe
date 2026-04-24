package utils

import (
	"dropshipbe/gateway/config"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateGuestToken(guestID string, jwtConfig config.JwtConfig) (string, int64, error) {
	exp := time.Now().Add(time.Duration(jwtConfig.ExpireHours) * time.Hour).Unix()
	claims := jwt.MapClaims{
		"guest_id": guestID,
		"exp":      exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Secret))
	return signedToken, exp, err
}

func ValidateGuestToken(tokenString string, jwtConfig config.JwtConfig) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims["guest_id"].(string), nil
}

func SanitizeMessage(msg string, maxLength int) string {
	msg = strings.TrimSpace(msg)
	msg = html.EscapeString(msg)

	if len(msg) > maxLength {
		msg = msg[:maxLength]
	}

	return msg
}
