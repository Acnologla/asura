package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

type redisConfig struct {
	URL      string `json:"REDIS_URL"`
	Password string `json:"REDIS_PASSWORD"`
}

func loadRedisEnv() redisConfig {
	var config = os.Getenv("REDIS_CONFIG")
	var redisConfig redisConfig
	json.Unmarshal([]byte(config), &redisConfig)
	return redisConfig
}

func Init() {
	config := loadRedisEnv()
	Client = redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
	})
	cmd := Client.Ping(context.Background())
	err := cmd.Err()
	if err != nil {
		log.Fatal(err)
	}
}
