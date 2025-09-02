package initialization

import (
	"context"
	"fmt"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedis(rCfg *config.Redis, logger logger.Interface) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", rCfg.Host, rCfg.Port),
		Password: rCfg.Password, // no password set
		DB:       rCfg.Database, // use default DB
		PoolSize: 10,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Redis initialization error", zap.Error(err))
	}

	return rdb
}
