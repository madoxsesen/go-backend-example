package redis

import "github.com/go-redis/redis/v8"

var Client *redis.Client

func SetupRedisClient() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
