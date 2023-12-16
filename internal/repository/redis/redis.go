package redis

import (
	"fmt"
	"github.com/Verce11o/yata-comments/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) *redis.Client {
	redisHost := fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client
}
