package config

import (
	"context"
	"elderly-care-backend/global"
	"github.com/redis/go-redis/v9"
)

func initRedis() {
	redisConfig := Config.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			global.Logger.Info("redis connect success")
			return nil
		},
	})

	global.RedisClient = redisClient
}
