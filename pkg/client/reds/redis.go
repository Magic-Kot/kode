package reds

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type ConfigDeps struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewClientRedis - создает клиента, подключаемый к базе данных по URL: reds://reds:12345@127.0.0.1:6379/reds
func NewClientRedis(ctx context.Context, cfg *ConfigDeps) (*redis.Client, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("creating a Redis client")
	logger.Debug().Msgf("reds config: %+v", cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		logger.Debug().Msgf("error connecting to Redis: %v", err)
		return nil, err
	}

	logger.Info().Msg("successful connection to Redis")

	return rdb, nil
}
