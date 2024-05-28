package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(addr, password string, db int) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisService{
		Client: client,
	}
}

func (r *RedisService) Start() error {
	ctx := context.Background()
	if _, err := r.Client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("Failed to ping Redis: %w", err)
	}
	log.Println("[Redis] Connection established")
	return nil
}

func (r *RedisService) Stop() error {
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("Failed to close Redis connection: %w", err)
	}
	log.Println("[Redis] Connection closed")
	return nil
}
