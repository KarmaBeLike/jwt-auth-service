package model

import "time"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID        string // ID записи в БД
	UserID    string // GUID пользователя
	TokenHash string // bcrypt-хеш Refresh токена
	IPAddress string // IP-адрес, на который выдан
	CreatedAt time.Time
	ExpiresAt time.Time
}
