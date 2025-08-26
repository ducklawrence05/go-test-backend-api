package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTP     `mapstructure:"http"`
	Postgres Postgres `mapstructure:"postgres"`
	Redis    Redis    `mapstructure:"redis"`
	Logger   Logger   `mapstructure:"logger"`
	JWT      JWT      `mapstructure:"jwt"`
	OTP      OTP      `mapstructure:"otp"`
	SMTP     SMTP     `mapstructure:"smtp"`
}

type HTTP struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type Postgres struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Dbname          string        `mapstructure:"dbname"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type Logger struct {
	LogLevel    string `mapstructure:"log_level"`
	FileLogName string `mapstructure:"file_log_name"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxSize     int    `mapstructure:"max_size"`
	Compress    bool   `mapstructure:"compress"`
}

type JWT struct {
	AccessTokenKey       string        `mapstructure:"access_token_key"`
	AccessTokenExpiresIn time.Duration `mapstructure:"access_token_expires_in"`

	RefreshTokenKey       string        `mapstructure:"refresh_token_key"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`

	RegisterTokenKey       string        `mapstructure:"register_token_key"`
	RegisterTokenExpiresIn time.Duration `mapstructure:"register_token_expires_in"`

	RestoreAccountTokenKey       string        `mapstructure:"restore_account_token_key"`
	RestoreAccountTokenExpiresIn time.Duration `mapstructure:"restore_account_token_expires_in"`
}

type OTP struct {
	RegisterKey          string        `mapstructure:"register_key"`
	RegisterTTL          time.Duration `mapstructure:"register_ttl"`
	RegisterRateLimit    int           `mapstructure:"register_rate_limit"`
	RegisterRateLimitTTL time.Duration `mapstructure:"register_rate_limit_ttl"`
	RegisterAttempts     int           `mapstructure:"register_attempts"`
	RegisterAttemptsTTL  time.Duration `mapstructure:"register_attempts_ttl"`

	RestoreAccountKey          string        `mapstructure:"restore_account_key"`
	RestoreAccountTTL          time.Duration `mapstructure:"restore_account_ttl"`
	RestoreAccountRateLimit    int           `mapstructure:"restore_account_rate_limit"`
	RestoreAccountRateLimitTTL time.Duration `mapstructure:"restore_account_rate_limit_ttl"`
	RestoreAccountAttempts     int           `mapstructure:"restore_account_attempts"`
	RestoreAccountAttemptsTTL  time.Duration `mapstructure:"restore_account_attempts_ttl"`
}

type SMTP struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	AppPassword string `mapstructure:"app_password"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	// read env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := &Config{}
	// Configure structure
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}
