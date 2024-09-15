package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	errCreateSession = errors.New("failed to create session")
	errGetSession    = errors.New("failed to get session")
	errDeleteSession = errors.New("failed to delete session")
)

type AuthRedisRepository struct {
	client redis.Client
}

func NewAuthRepository(client *redis.Client) *AuthRedisRepository {
	return &AuthRedisRepository{client: *client}
}

// CreateSession - creating a user session or updating it
func (a *AuthRedisRepository) CreateSession(ctx context.Context, key string, value interface{}) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'CreateSession' method")
	logger.Debug().Msgf("redis: create session by key: %s, value: %s", key, value)

	json := a.client.JSONSet(ctx, key, ".", value)

	a.client.Expire(ctx, key, 4*time.Hour) // Hard code

	if json.Err() == nil && json.Val() == "" {
		logger.Debug().Msgf("failed to create session. redis: %s", json.Err())
		return "", errCreateSession
	}

	return json.String(), nil
}

// GetSession - getting a user session by id
func (a *AuthRedisRepository) GetSession(ctx context.Context, key string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'GetSession' method")
	logger.Debug().Msgf("redis: get session by key: %s", key)

	session := a.client.JSONGet(ctx, key)

	if session.Err() == nil && session.Val() == "" {
		logger.Debug().Msgf("failed to get session: %s", session.Err())
		return "", errGetSession
	}

	return session.Result()
}

// DeleteSession - deleting a user session by id
func (a *AuthRedisRepository) DeleteSession(ctx context.Context, key string) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Redis using the 'DeleteSession' method")
	logger.Debug().Msgf("redis: delete session by key: %s", key)

	session := a.client.JSONDel(ctx, key, ".")

	if session.Err() != nil || session.Val() == 0 {
		logger.Debug().Msgf("failed to delete session. val: %d, err: %s", session.Val(), session.Err())
		return errDeleteSession
	}

	return nil
}
