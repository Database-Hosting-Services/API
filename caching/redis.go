package caching

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the go-redis client.
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient initializes and returns a new RedisClient instance.
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: rdb}, nil
}

// Set sets a key-value pair in Redis with the given expiration.
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves the value associated with the given key.
func (r *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	return r.Client.Get(ctx, key).Result()
}

// Delete deletes the value associated with the given key
func (r *RedisClient) Delete(key string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis.
func (r *RedisClient) Exists(key string) (bool, error) {
	ctx := context.Background()
	exists, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisClient) Eval(ctx context.Context, script string, args ...interface{}) (interface{}, error){
	keys := make([]string, len(args))
	for i, v := range args {
		switch t := v.(type) {
		case string:
			keys[i] = v.(string)
		case int:
			keys[i] = strconv.Itoa(v.(int))
		default:
			return nil, fmt.Errorf("Unknown type: %T", t)
		}
	}
	return  r.Client.Eval(ctx, script, keys).Result()
}


// Close closes the Redis client connection.
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
