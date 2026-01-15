package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

var cacheInstance *CacheService

func NewCacheService() *CacheService {
	if cacheInstance != nil {
		return cacheInstance
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "anime-redis:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		// We don't return nil, so the app can still run without cache
		// But ideally we should handle this gracefully
	} else {
		log.Println("Connected to Redis successfully")
	}

	cacheInstance = &CacheService{
		client: client,
		ctx:    ctx,
	}
	return cacheInstance
}

func (s *CacheService) Get(key string, target interface{}) bool {
	if s.client == nil {
		return false
	}

	val, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Printf("[Cache] ERROR Get: %v for key %s", err, key)
		return false
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		log.Printf("[Cache] ERROR Unmarshal: %v", err)
		return false
	}

	return true
}

func (s *CacheService) Set(key string, value interface{}, expiration time.Duration) {
	if s.client == nil {
		return
	}

	jsonVal, err := json.Marshal(value)
	if err != nil {
		log.Printf("Redis Marshal Error: %v", err)
		return
	}

	err = s.client.Set(s.ctx, key, jsonVal, expiration).Err()
	if err != nil {
		log.Printf("[Cache] ERROR Set: %v for key %s", err, key)
	}
}
