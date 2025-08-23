package initialization

import (
	"fmt"
	"log"
	"strings"

	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() (*setting.Config, error) {
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

	var cfg setting.Config
	// Configure structure
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode configuration %w", err)
	}
	return &cfg, nil
}
