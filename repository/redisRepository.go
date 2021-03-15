package repository

import (
	"github.com/go-redis/redis/v7"
	"github.com/braianpablodiaz/meli-proxy/environment"
	"fmt"
)

func InitialConnection() *redis.Client{

	client := redis.NewClient(&redis.Options{
		Addr:     environment.GetEnv("REDIS_URL", "localhost:6379"),
		Password: environment.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Printf(err.Error())
	}

	return client
}