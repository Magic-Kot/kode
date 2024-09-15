package main

import (
	"context"

	"github.com/Magic-Kot/code/internal/config"
	"github.com/Magic-Kot/code/internal/controllers"
	"github.com/Magic-Kot/code/internal/delivery/httpecho"
	"github.com/Magic-Kot/code/internal/middleware"
	"github.com/Magic-Kot/code/internal/repository/postgres"
	"github.com/Magic-Kot/code/internal/repository/redis"
	"github.com/Magic-Kot/code/internal/services/auth"
	"github.com/Magic-Kot/code/internal/services/note"
	"github.com/Magic-Kot/code/internal/services/user"
	"github.com/Magic-Kot/code/pkg/client/postg"
	"github.com/Magic-Kot/code/pkg/client/reds"
	"github.com/Magic-Kot/code/pkg/httpserver"
	"github.com/Magic-Kot/code/pkg/logging"
	"github.com/Magic-Kot/code/pkg/speller"
	"github.com/Magic-Kot/code/pkg/utils/jwt_token"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	// read config
	var cfg config.Config

	err := cleanenv.ReadConfig("internal/config/config.yml", &cfg) // Local: internal/config/config.yml Docker: config.yml
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
	}

	// create logger
	logCfg := logging.LoggerDeps{
		LogLevel: cfg.LoggerDeps.LogLevel,
	}

	logger, err := logging.NewLogger(&logCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	logger.Info().Msg("init logger")

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	logger.Debug().Msgf("config: %+v", cfg)

	// create server
	serv := httpserver.ConfigDeps{
		Host:    cfg.ServerDeps.Host,
		Port:    cfg.ServerDeps.Port,
		Timeout: cfg.ServerDeps.Timeout,
	}

	server := httpserver.NewServer(&serv)

	// create client Postgres
	repo := postg.ConfigDeps{
		MaxAttempts: cfg.PostgresDeps.MaxAttempts,
		Delay:       cfg.PostgresDeps.Delay,
		Username:    cfg.PostgresDeps.Username,
		Password:    cfg.PostgresDeps.Password,
		Host:        cfg.PostgresDeps.Host,
		Port:        cfg.PostgresDeps.Port,
		Database:    cfg.PostgresDeps.Database,
		SSLMode:     cfg.PostgresDeps.SSLMode,
	}

	pool, err := postg.NewClient(ctx, &repo)
	if err != nil {
		logger.Fatal().Err(err).Msgf("NewClient: %s", err)
	}

	// create client Redis for refresh tokens
	redisCfg := reds.ConfigDeps{
		Username: cfg.RedisDeps.Username,
		Password: cfg.RedisDeps.Password,
		Host:     cfg.RedisDeps.Host,
		Port:     cfg.RedisDeps.Port,
		Database: cfg.RedisDeps.Database,
	}

	clientRedis, err := reds.NewClientRedis(ctx, &redisCfg)
	if err != nil {
		logger.Fatal().Err(err).Msgf("redis refresh tokens: %s", err)
	}

	// create tokenJWT
	tokenCfg := jwt_token.TokenJWTDeps{
		SigningKey:      cfg.AuthDeps.SigningKey,
		AccessTokenTTL:  cfg.AuthDeps.AccessTokenTTL,
		RefreshTokenTTL: cfg.AuthDeps.RefreshTokenTTL,
	}

	tokenJWT, err := jwt_token.NewTokenJWT(&tokenCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init tokenJWT")
	}

	// create validator
	validate := validator.New()

	rds := redis.NewAuthRepository(clientRedis)
	middlewareUser := middleware.NewMiddleware(logger, tokenJWT)

	// Auth
	authRepository := postgres.NewAuthPostgresRepository(pool)
	authService := auth.NewAuthService(authRepository, rds, tokenJWT)
	authController := controllers.NewApiAuthController(authService, logger, validate)
	httpecho.SetAuthRoutes(server.Server(), authController)

	// User
	userRepository := postgres.NewUserRepository(pool)
	userService := user.NewUserService(userRepository, rds)
	userController := controllers.NewApiController(userService, logger, validate)
	httpecho.SetUserRoutes(server.Server(), userController, middlewareUser)

	// Note
	noteRepository := postgres.NewNoteRepository(pool)
	spell := speller.NewSpeller(cfg.Speller.Url)
	noteService := note.NewNoteService(noteRepository, spell)
	noteController := controllers.NewApiNoteController(noteService, logger, validate)
	httpecho.SetNoteRoutes(server.Server(), noteController, middlewareUser)

	// start server
	logger.Info().Msg("starting server")

	if err := server.Start(); err != nil {
		logger.Fatal().Msgf("serverStart: %v", err)
	}
}
