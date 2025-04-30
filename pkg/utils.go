package pkg

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/KarmaBeLike/jwt-auth-service/internal/dto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func SplitToken(token string) (string, string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid refresh token format")
	}
	return parts[0], parts[1], nil
}

func SendMockEmailWarning(email string, oldIP string, newIP string) {
	log.Printf("[MOCK EMAIL] Security warning to %s: IP address changed from %s to %s\n", email, oldIP, newIP)
}

func GenerateJWT(data dto.JWTData, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": data.UserID,
		"ip":      data.IP,
		"id":      data.ID, // access_token_id
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
