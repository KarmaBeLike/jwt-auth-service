package handler

import (
	"context"
	"net/http"

	"github.com/KarmaBeLike/jwt-auth-service/internal/dto"
	"github.com/KarmaBeLike/jwt-auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/token", h.GenerateToken)
	r.POST("/refresh", h.RefreshToken)
}

// Генерация нового токена
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	ipAddress := c.ClientIP()

	tokens, err := h.authService.GenerateTokens(context.Background(), userID, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}

// Обновление токенов
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ipAddress := c.ClientIP()

	tokens, err := h.authService.RefreshTokens(context.Background(), body.RefreshToken, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}
