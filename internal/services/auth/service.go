package auth

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Magic-Kot/kode/internal/models"
	"github.com/Magic-Kot/kode/pkg/utils/hash"
	"github.com/Magic-Kot/kode/pkg/utils/jwt_token"

	"github.com/rs/zerolog"
)

var (
	errAutorizationUser = errors.New("invalid refresh token")
)

type AuthPostgresRepository interface {
	SignIn(ctx context.Context, user *models.UserAuthorization) (*models.UserAuthorization, error)
}

type AuthRedisRepository interface {
	CreateSession(ctx context.Context, key string, value interface{}) (string, error)
	GetSession(ctx context.Context, key string) (string, error)
	DeleteSession(ctx context.Context, key string) error
}

type AuthService struct {
	AuthPostgresRepository AuthPostgresRepository
	AuthRepository         AuthRedisRepository
	token                  *jwt_token.Manager
}

func NewAuthService(authPostgresRepository AuthPostgresRepository, authRepository AuthRedisRepository, token *jwt_token.Manager) *AuthService {
	return &AuthService{
		AuthPostgresRepository: authPostgresRepository,
		AuthRepository:         authRepository,
		token:                  token,
	}
}

// SignIn - user authentication
func (s *AuthService) SignIn(ctx context.Context, user *models.UserAuthorization) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'SignIn' service")

	var res models.Tokens

	passwordHash := hash.GenerateHash(user.Password)

	user, err := s.AuthPostgresRepository.SignIn(ctx, user)
	if err != nil {
		return res, err
	}

	if passwordHash != user.Password {
		logger.Debug().Msg("invalid password")

		err = errors.New("invalid password")
		return res, err
	}

	res.AccessToken, err = s.token.NewJWT(strconv.Itoa(user.ID))
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.token.NewRefreshToken(strconv.Itoa(user.ID))

	RefreshTokenHash, err := hash.GenerateHashBcrypt(res.RefreshToken)
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: RefreshTokenHash,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
	}

	_, err = s.AuthRepository.CreateSession(ctx, strconv.Itoa(user.ID), session)
	if err != nil {
		logger.Debug().Msgf("failed to create session: %s", err)
		return res, errors.New("failed to create session")
	}

	return res, nil
}

// RefreshToken - getting new refresh and access tokens
func (s *AuthService) RefreshToken(ctx context.Context, refresh string) (models.Tokens, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'RefreshToken' service")

	var (
		res   models.Tokens
		token models.Session
	)

	passwordDecode, err := s.token.ParseRefreshToken(refresh)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	userId := strings.Fields(passwordDecode)

	session, err := s.AuthRepository.GetSession(ctx, userId[1])

	err = json.Unmarshal([]byte(session), &token)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	err = hash.CompareHashBcrypt(refresh, token.RefreshToken)
	if err != nil {
		logger.Debug().Msgf("invalid refresh token: %s", err)
		// redirecting to the authorization page
		return res, errAutorizationUser
	}

	// generation of new access and refresh tokens
	res.AccessToken, err = s.token.NewJWT(userId[1])
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.token.NewRefreshToken(userId[1])

	RefreshTokenHash, err := hash.GenerateHashBcrypt(res.RefreshToken)
	if err != nil {
		return res, err
	}

	newSession := models.Session{
		RefreshToken: RefreshTokenHash,
		ExpiresAt:    time.Now().Add(s.token.RefreshTokenTTL()),
	}

	// update session
	_, err = s.AuthRepository.CreateSession(ctx, userId[1], newSession)
	if err != nil {
		logger.Debug().Msgf("failed to update session: %s", err)
		return res, errors.New("failed to update session")
	}

	return res, nil
}
