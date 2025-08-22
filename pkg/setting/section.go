package setting

import "time"

type Config struct {
	Server   ServerSetting   `mapstructure:"server"`
	Postgres PostgresSetting `mapstructure:"postgres"`
	Logger   LoggerSetting   `mapstructure:"logger"`
	JWT      JWTSetting      `mapstructure:"jwt"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type PostgresSetting struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Dbname          string        `mapstructure:"dbname"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

type JWTSetting struct {
	AccessTokenKey        string `mapstructure:"access_token_key"`
	AccessTokenExpiresIn  int    `mapstructure:"access_token_expires_in"`
	RefreshTokenKey       string `mapstructure:"refresh_token_key"`
	RefreshTokenExpiresIn int    `mapstructure:"refresh_token_expires_in"`
}

type LoggerSetting struct {
	Log_level     string `mapstructure:"log_level"`
	File_log_name string `mapstructure:"file_log_name"`
	Max_backups   int    `mapstructure:"max_backups"`
	Max_age       int    `mapstructure:"max_age"`
	Max_size      int    `mapstructure:"max_size"`
	Compress      bool   `mapstructure:"compress"`
}
