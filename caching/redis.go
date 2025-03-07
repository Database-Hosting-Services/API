package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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

// Set sets a key-value pair in Redis with the given expiration time.
// to set the key to a json object read from a struct object value should be a pointer to the struct.
// Set uses the SetJson to do the marshal operation.
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	if t := reflect.ValueOf(value); t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return r.SetJson(key, value, expiration)
	}
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// obj is the struct object from which the value for the key will be generated from in the format of json
func (r *RedisClient) SetJson(key string, obj interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return r.Client.Set(ctx, key, jsonData, expiration).Err()
}

// Get retrieves the value associated with the given key.
// if the read value is a json object and you want to read it into a struct object dest should be a pointer to that struct.
// Note: that in this case the return values will be nil, error.
// if the read value is a premetive type dest should be set to nil and the value will be returned.
func (r *RedisClient) Get(key string, dest interface{}) (interface{}, error) {
	if t := reflect.ValueOf(dest); t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return nil, r.GetJson(key, dest)
	}
	ctx := context.Background()
	return r.Client.Get(ctx, key).Result()
}

// dest is a pointer to the struct object in which data will be read into
func (r *RedisClient) GetJson(key string, dest interface{}) error {
	ctx := context.Background()
	jsonData, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(jsonData), dest); err != nil {
		return err
	}
	return nil
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

func (r *RedisClient) Eval(ctx context.Context, script string, args ...interface{}) (interface{}, error) {
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
	return r.Client.Eval(ctx, script, keys).Result()
}

// Close closes the Redis client connection.
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
