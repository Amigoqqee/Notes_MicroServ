package caching

import (
	"fmt"
	"notes/internal/config"

	"github.com/go-redis/redis"
)

func NewCaching(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})

	if err := client.Ping().Err(); err != nil {
		fmt.Printf("Ошибка подключения к Redis: %v\n", err)
		return nil, fmt.Errorf("не удалось подключиться к Redis: %w", err)
	}

	fmt.Println("Успешно подключено к Redis")

	return client, nil
}
