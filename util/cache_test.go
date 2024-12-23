package util

import (
	"testing"
	"time"
)

func TestGetAndSetCache(t *testing.T) {
	// Test case 1: Set and get a string value
	key1 := "testKey1"
	value1 := "testValue1"
	SetCache(key1, value1, 1*time.Minute)

	got, exists := GetCache(key1)
	if !exists {
		t.Errorf("GetCache(%q) = %v, want %v", key1, exists, true)
	}
	if got != value1 {
		t.Errorf("GetCache(%q) = %v, want %v", key1, got, value1)
	}

	// Test case 2: Set and get an integer value
	key2 := "testKey2"
	value2 := 42
	SetCache(key2, value2, 1*time.Minute)

	got, exists = GetCache(key2)
	if !exists {
		t.Errorf("GetCache(%q) = %v, want %v", key2, exists, true)
	}
	if got != value2 {
		t.Errorf("GetCache(%q) = %v, want %v", key2, got, value2)
	}

	// Test case 3: Get a non-existent key
	key3 := "nonExistentKey"
	_, exists = GetCache(key3)
	if exists {
		t.Errorf("GetCache(%q) exists = %v, want %v", key3, exists, false)
	}

	// Test case 4: Set with expiration
	key4 := "expiringKey"
	value4 := "expiringValue"
	SetCache(key4, value4, 100*time.Millisecond)

	// Wait for the cache item to expire
	time.Sleep(200 * time.Millisecond)

	_, exists = GetCache(key4)
	if exists {
		t.Errorf("GetCache(%q) exists = %v, want %v", key4, exists, false)
	}
}
