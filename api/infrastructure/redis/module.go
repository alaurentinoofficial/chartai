package infrastructure_redis

import (
	infrastructure_config "github.com/alaurentinoofficial/chartai/infrastructure/configs"
	"github.com/go-redis/redis/v9"
	"go.uber.org/fx"
)

func NewRedisConnection(lc fx.Lifecycle, config *infrastructure_config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
	})

	return client
}
