package pkg

import (
	"errors"
	"fmt"
	"net/smtp"
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

func SendEmailWarning(email string, userID, oldIP string, newIP string) error {
	from := "your.email@gmail.com"
	password := "your_generated_app_pass"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := fmt.Sprintf("Subject: IP Change Warning\n\nUser %s has changed IP from %s to %s", userID, oldIP, newIP)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
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
