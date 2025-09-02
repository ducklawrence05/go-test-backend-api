package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTP     HTTP     `envPrefix:"HTTP_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
	Redis    Redis    `envPrefix:"REDIS_"`
	Logger   Logger   `envPrefix:"LOGGER_"`
	JWT      JWT      `envPrefix:"JWT_"`
	OTP      OTP      `envPrefix:"OTP_"`
	SMTP     SMTP     `envPrefix:"SMTP_"`
}

type HTTP struct {
	Url             string        `env:"URL"`
	Port            int           `env:"PORT"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

type Postgres struct {
	Host            string        `env:"HOST"`
	Port            int           `env:"PORT"`
	Username        string        `env:"USERNAME"`
	Password        string        `env:"PASSWORD"`
	Dbname          string        `env:"DBNAME"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
}

type Redis struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Password string `env:"PASSWORD"`
	Database int    `env:"DATABASE"`
}

type Logger struct {
	LogLevel    string `env:"LOG_LEVEL"`
	FileLogName string `env:"FILE_LOG_NAME"`
	MaxBackups  int    `env:"MAX_BACKUPS"`
	MaxAge      int    `env:"MAX_AGE"`
	MaxSize     int    `env:"MAX_SIZE"`
	Compress    bool   `env:"COMPRESS"`
}

type JWT struct {
	AccessTokenKey       string        `env:"ACCESS_TOKEN_KEY"`
	AccessTokenExpiresIn time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN"`

	RefreshTokenKey       string        `env:"REFRESH_TOKEN_KEY"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN"`

	RegisterTokenKey       string        `env:"REGISTER_TOKEN_KEY"`
	RegisterTokenExpiresIn time.Duration `env:"REGISTER_TOKEN_EXPIRES_IN"`

	RestoreAccountTokenKey       string        `env:"RESTORE_ACCOUNT_TOKEN_KEY"`
	RestoreAccountTokenExpiresIn time.Duration `env:"RESTORE_ACCOUNT_TOKEN_EXPIRES_IN"`
}

type OTP struct {
	RegisterKey          string        `env:"REGISTER_KEY"`
	RegisterTTL          time.Duration `env:"REGISTER_TTL"`
	RegisterRateLimit    int           `env:"REGISTER_RATE_LIMIT"`
	RegisterRateLimitTTL time.Duration `env:"REGISTER_RATE_LIMIT_TTL"`
	RegisterAttempts     int           `env:"REGISTER_ATTEMPTS"`
	RegisterAttemptsTTL  time.Duration `env:"REGISTER_ATTEMPTS_TTL"`

	RestoreAccountKey          string        `env:"RESTORE_ACCOUNT_KEY"`
	RestoreAccountTTL          time.Duration `env:"RESTORE_ACCOUNT_TTL"`
	RestoreAccountRateLimit    int           `env:"RESTORE_ACCOUNT_RATE_LIMIT"`
	RestoreAccountRateLimitTTL time.Duration `env:"RESTORE_ACCOUNT_RATE_LIMIT_TTL"`
	RestoreAccountAttempts     int           `env:"RESTORE_ACCOUNT_ATTEMPTS"`
	RestoreAccountAttemptsTTL  time.Duration `env:"RESTORE_ACCOUNT_ATTEMPTS_TTL"`
}

type SMTP struct {
	Host        string `env:"HOST"`
	Port        int    `env:"PORT"`
	Username    string `env:"USERNAME"`
	AppPassword string `env:"APP_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
