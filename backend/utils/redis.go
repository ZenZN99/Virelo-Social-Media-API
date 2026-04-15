package utils

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

var Redis *RedisService

func InitRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("Invalid Redis URL: " + err.Error())
	}

	client := redis.NewClient(opt)

	// test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	Redis = &RedisService{
		client: client,
	}
}

func GetRedis() *redis.Client {
	return Redis.client
}
