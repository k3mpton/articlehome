package rediss

import (
	"context"
	getenvfield "legi/newspapers/project/utils/GetEnvField"
	"log"

	redis "github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	log.Println(getenvfield.Getenv("REDIS_ADDR"), "==========================================================")
	client := redis.NewClient(&redis.Options{
		Addr:     getenvfield.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if cmd := client.Ping(ctx); cmd.Err() != nil {
		log.Fatalf("Failed connection redis: %v", cmd.Err())
	}

	return client
}
