package caching_test

import (
	"DBHS/caching"
	"errors"
	"github.com/redis/go-redis/v9"
	"reflect"
	"testing"
	"time"
)

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRedisClient(t *testing.T) {
	client, err := caching.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	t.Run("Set and Get Struct", func(t *testing.T) {
		key := "struct:" + t.Name()
		user := &testUser{Name: "Alice", Age: 30}

		// Test Set with struct
		if err := client.Set(key, user, time.Minute); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// Test Get with struct
		var result testUser
		_, err := client.Get(key, &result)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if !reflect.DeepEqual(*user, result) {
			t.Errorf("Expected %+v, got %+v", *user, result)
		}

		// Cleanup
		client.Delete(key)
	})

	t.Run("Set and Get Primitive Types", func(t *testing.T) {
		tests := []struct {
			key    string
			value  interface{}
			expect string
		}{
			{"string:" + t.Name(), "hello", "hello"},
			{"int:" + t.Name(), 42, "42"},
			{"bool:" + t.Name(), true, "1"},
			{"float:" + t.Name(), 3.14, "3.14"},
		}

		for _, tt := range tests {
			if err := client.Set(tt.key, tt.value, time.Minute); err != nil {
				t.Fatalf("Set failed for %s: %v", tt.key, err)
			}

			val, err := client.Get(tt.key, nil)
			if err != nil {
				t.Fatalf("Get failed for %s: %v", tt.key, err)
			}

			if val != tt.expect {
				t.Errorf("For %s expected %s, got %s", tt.key, tt.expect, val)
			}

			client.Delete(tt.key)
		}
	})

	t.Run("JSON Marshaling Direct", func(t *testing.T) {
		key := "json:" + t.Name()
		user := testUser{Name: "Bob", Age: 25}

		// Test SetJson directly
		if err := client.SetJson(key, user, time.Minute); err != nil {
			t.Fatalf("SetJson failed: %v", err)
		}

		// Test GetJson directly
		var result testUser
		if err := client.GetJson(key, &result); err != nil {
			t.Fatalf("GetJson failed: %v", err)
		}

		if !reflect.DeepEqual(user, result) {
			t.Errorf("Expected %+v, got %+v", user, result)
		}

		client.Delete(key)
	})

	t.Run("Exists and Delete", func(t *testing.T) {
		key := "exists:" + t.Name()

		// Test initial existence
		exists, err := client.Exists(key)
		if err != nil {
			t.Fatalf("Exists check failed: %v", err)
		}
		if exists {
			t.Error("Key should not exist initially")
		}

		// Set key and check existence
		if err := client.Set(key, "test", time.Minute); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		exists, err = client.Exists(key)
		if err != nil {
			t.Fatalf("Exists check failed: %v", err)
		}
		if !exists {
			t.Error("Key should exist after Set")
		}

		// Test Delete
		if err := client.Delete(key); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err = client.Exists(key)
		if err != nil {
			t.Fatalf("Exists check failed: %v", err)
		}
		if exists {
			t.Error("Key should not exist after Delete")
		}
	})

	t.Run("Connection Close", func(t *testing.T) {
		localClient, _ := caching.NewRedisClient("localhost:6379", "", 0)
		if err := localClient.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}

		// Test operation after close
		err := localClient.Set("closed", "test", 0)
		if err == nil || !errors.Is(err, redis.ErrClosed) {
			t.Errorf("Expected closed connection error, got: %v", err)
		}
	})

	t.Run("Edge Cases", func(t *testing.T) {
		// Test nil value
		key := "nil:" + t.Name()
		if err := client.Set(key, nil, time.Minute); err != nil {
			t.Errorf("Setting nil value failed: %v", err)
		}

		// Test non-existent key
		_, err := client.Get("non-existent", nil)
		if !errors.Is(err, redis.Nil) {
			t.Errorf("Expected redis.Nil error, got: %v", err)
		}
	})
}
