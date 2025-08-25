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
	AccessTokenKey  string `mapstructure:"access_token_key"`
	RefreshTokenKey string `mapstructure:"refresh_token_key"`

	AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`
}

type OTP struct {
	EmailVerifyKey       string        `mapstructure:"email_verify_key"`
	EmailVerifyExpiresIn time.Duration `mapstructure:"email_verify_expires_in"`
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
