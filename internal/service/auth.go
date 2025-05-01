package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/KarmaBeLike/jwt-auth-service/internal/dto"
	"github.com/KarmaBeLike/jwt-auth-service/internal/model"
	"github.com/KarmaBeLike/jwt-auth-service/internal/repository"
	"github.com/KarmaBeLike/jwt-auth-service/pkg"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.AuthRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// Генерация новых токенов
func (s *AuthService) GenerateTokens(ctx context.Context, userID, ipAddress string) (dto.TokenPair, error) {
	// 1. Генерируем access_token_id
	accessTokenID := pkg.GenerateUUID()

	// 2. Генерируем access токен с этим ID
	accessToken, err := pkg.GenerateJWT(dto.JWTData{
		UserID: userID,
		IP:     ipAddress,
		ID:     accessTokenID,
	}, s.jwtSecret, s.accessTTL)
	if err != nil {
		return dto.TokenPair{}, err
	}

	// 3. Генерация refresh токена и хеша
	rawRefreshToken, tokenHash, err := s.generateRefreshToken()
	if err != nil {
		return dto.TokenPair{}, err
	}

	// 4. Создание записи в БД
	refreshRecord := model.RefreshToken{
		ID:            pkg.GenerateUUID(),
		UserID:        userID,
		TokenHash:     tokenHash,
		IPAddress:     ipAddress,
		AccessTokenID: accessTokenID,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(s.refreshTTL),
	}

	err = s.repo.SaveRefreshToken(ctx, refreshRecord)
	if err != nil {
		return dto.TokenPair{}, err
	}

	// 5. Возврат обоих токенов
	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshRecord.ID + ":" + rawRefreshToken,
	}, nil
}

// Обновление токенов
func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenFull, currentIP string) (dto.TokenPair, error) {
	id, providedToken, err := pkg.SplitToken(refreshTokenFull)
	if err != nil {
		return dto.TokenPair{}, err
	}

	storedToken, err := s.repo.GetRefreshToken(ctx, id)
	if err != nil {
		return dto.TokenPair{}, err
	}
	if storedToken == nil {
		return dto.TokenPair{}, errors.New("refresh token not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedToken.TokenHash), []byte(providedToken))
	if err != nil {
		return dto.TokenPair{}, errors.New("invalid refresh token")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return dto.TokenPair{}, errors.New("refresh token expired")
	}

	if storedToken.IPAddress != currentIP {
		pkg.SendMockEmailWarning(storedToken.UserID, storedToken.IPAddress, currentIP)
	}

	err = s.repo.DeleteRefreshToken(ctx, id)
	if err != nil {
		return dto.TokenPair{}, err
	}

	return s.GenerateTokens(ctx, storedToken.UserID, currentIP)
}

func (s *AuthService) generateRefreshToken() (string, string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", err
	}
	rawToken := base64.StdEncoding.EncodeToString(randomBytes)

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return rawToken, string(hashBytes), nil
}
