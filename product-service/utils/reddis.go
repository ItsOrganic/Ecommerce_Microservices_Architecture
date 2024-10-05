package utils

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis: ", err)
	}
	log.Println("Connected to Redis")
}

// SetKey sets a key-value pair in Redis with an optional expiration time
func SetKey(key string, value interface{}, expiration time.Duration) error {
	err := RDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting key %s: %v", key, err)
		return err
	}
	log.Printf("Key %s set successfully", key)
	return nil
}

// GetKey retrieves a value from Redis by key
func GetKey(key string) (string, error) {
	val, err := RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key %s does not exist", key)
		return "", nil
	} else if err != nil {
		log.Printf("Error getting key %s: %v", key, err)
		return "", err
	}
	log.Printf("Key %s retrieved successfully", key)
	return val, nil
}
