// redis_test.go
package caching

import (
	"testing"
	"time"
	"DBHS/caching"
)

// TestRedisClient tests basic Redis operations: Set, Get, Exists, and Delete.
func TestRedisClient(t *testing.T) {
	// Create a new Redis client.
	client := caching.NewRedisClient("localhost:6379", "", 0)
	defer client.Close()

	// Test Set operation.
	if err := client.Set("testKey", "testValue", 10 * time.Second); err != nil {
		t.Fatalf("Set() returned an unexpected error: %v", err)
	}

	// Test Get operation.
	val, err := client.Get("testKey")
	if err != nil {
		t.Fatalf("Get() returned an unexpected error: %v", err)
	}
	if val != "testValue" {
		t.Errorf("Get() = %v, want %v", val, "testValue")
	}

	// Test Exists operation.
	exists, err := client.Exists("testKey")
	if err != nil {
		t.Fatalf("Exists() returned an unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("Exists() = false, want true")
	}

	// Test Delete operation.
	if err := client.Delete("testKey"); err != nil {
		t.Fatalf("Delete() returned an unexpected error: %v", err)
	}

	// Verify the key no longer exists.
	exists, err = client.Exists("testKey")
	if err != nil {
		t.Fatalf("Exists() after Delete returned an unexpected error: %v", err)
	}
	if exists {
		t.Errorf("Exists() = true after Delete, want false")
	}
}
