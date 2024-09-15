package config

import "time"

type Config struct {
	ServerDeps   ServerDeps   `env:"SERVER" yaml:"server"`
	PostgresDeps PostgresDeps `env:"POSTGRES" yaml:"postgres"`
	RedisDeps    RedisDeps    `env:"REDIS" yaml:"redis"`
	LoggerDeps   LoggerDeps   `env:"LOGGER" yaml:"logger"`
	AuthDeps     AuthDeps     `env:"AUTH" yaml:"auth"`
	Speller      Speller      `env:"SPELLER" yaml:"speller"`
}

type ServerDeps struct {
	Host    string        `env:"HOST"  yaml:"host" env-default:"localhost"`
	Port    string        `env:"PORT" yaml:"port" env-default:":8000"`
	Timeout time.Duration `env:"TIMEOUT" yaml:"timeout" env-default:"5s"`
}

type PostgresDeps struct {
	MaxAttempts int           `env:"MAX_ATTEMPTS" yaml:"maxAttempts" env-default:"3"`
	Delay       time.Duration `env:"DELAY" yaml:"delay" env-default:"10s"`
	Username    string        `env:"USERNAME_POSTGRES"  yaml:"username" env-default:"postgres"`
	Password    string        `env:"PASSWORD" yaml:"password" env-default:"postgres"`
	Host        string        `env:"HOST" yaml:"host" env-default:"127.0.0.1"`
	Port        string        `env:"PORT" yaml:"port" env-default:"5432"`
	Database    string        `env:"DATABASE" yaml:"database" env-default:"postgres"`
	SSLMode     string        `env:"MODELESS" yaml:"sslMode" env-default:"disable"`
}

type RedisDeps struct {
	Username string `env:"USERNAME_REDIS"  yaml:"username" env-default:"reds"`
	Password string `env:"PASSWORD_REDIS" yaml:"password" env-default:""`
	Host     string `env:"HOST_REDIS" yaml:"host" env-default:"127.0.0.1"`
	Port     string `env:"PORT_REDIS" yaml:"port" env-default:"6379"`
	Database string `env:"DATABASE_REDIS" yaml:"database" env-default:"0"`
}

type LoggerDeps struct {
	LogLevel string `env:"LOG_LEVEL" yaml:"logLevel" env-default:"info"`
}

type AuthDeps struct {
	SigningKey      string        `env:"SIGNING_KEY" yaml:"signingKey" env-default:""`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" yaml:"accessTokenTTL" env-default:"1h"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" yaml:"refreshTokenTTL" env-default:"4h"`
}

type Speller struct {
	Url string `env:"SPELLER_URL" yaml:"url" env-default:"https://speller.yandex.net/services/spellservice.json/checkTexts"`
}
