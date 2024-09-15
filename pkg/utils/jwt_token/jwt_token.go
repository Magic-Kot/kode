package jwt_token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token interface {
	NewJWT(id string) (string, error)
	ParseToken(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type TokenJWTDeps struct {
	SigningKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokenJWT(cfg *TokenJWTDeps) (*Manager, error) {
	if cfg.SigningKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{
		signingKey:      cfg.SigningKey,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}, nil
}

// NewJWT - генерация JWT токена
func (m *Manager) NewJWT(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   id,
	})

	return token.SignedString([]byte(m.signingKey))
}

// ParseToken - парсинг токена
func (m *Manager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

// NewRefreshToken - генерация Refresh токена
func (m *Manager) NewRefreshToken(id string) string {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "" //, err
	}

	data := fmt.Sprintf("\"id\": %s %s", id, string(b))
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return encoded
}

func (m *Manager) ParseRefreshToken(refreshToken string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (m *Manager) RefreshTokenTTL() time.Duration {
	return m.refreshTokenTTL
}
