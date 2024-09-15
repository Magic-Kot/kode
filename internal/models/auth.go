package models

import "time"

type UserAuthorization struct {
	ID        int    `json:"id"`
	GUID      string `json:"guid"`
	IPAddress string `json:"ip_address"`
	Username  string `json:"login"        validate:"required,min=1,max=20"`
	Password  string `json:"password"     validate:"required,min=1,max=20"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Session struct {
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
