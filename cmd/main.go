package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/KarmaBeLike/jwt-auth-service/config"
	"github.com/KarmaBeLike/jwt-auth-service/internal/database"
	"github.com/KarmaBeLike/jwt-auth-service/internal/handler"
	"github.com/KarmaBeLike/jwt-auth-service/internal/repository"
	"github.com/KarmaBeLike/jwt-auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
		return
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("error running migrations", slog.Any("error", err))
		return
	}
	jwtSecret := (os.Getenv("JWT_SECRET"))
	accessTTL := time.Minute * 15
	refreshTTL := time.Hour * 24 * 7

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtSecret, accessTTL, refreshTTL)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	authHandler.RegisterRoutes(r)

	// serverPort := os.Getenv("SERVER_PORT")
	r.Run("localhost:8080")
}
