package db

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	db := NewDB()

	db.Set("foo", &StringValue{Val: "bar"})

	val, ok := db.Get("foo")
	if !ok {
		t.Fatal("Expected key to exist")
	}

	strVal, ok := val.(*StringValue)
	if !ok {
		t.Fatal("Expected value to be of type StringValue")
	}

	if strVal.Get() != "bar" {
		t.Fatalf("Expected 'bar', got '%s'", strVal.Get())
	}
}

func TestDelete(t *testing.T) {
	db := NewDB()

	db.Set("foo", &StringValue{Val: "bar"})
	deleted := db.Delete("foo")
	if !deleted {
		t.Fatal("Expected key to be deleted")
	}

	_, ok := db.Get("foo")
	if ok {
		t.Fatal("Expected key to be gone")
	}
}

func TestFlushAll(t *testing.T) {
	db := NewDB()

	db.Set("a", &StringValue{Val: "1"})
	db.Set("b", &StringValue{Val: "2"})
	db.FlushAll()

	keys := db.Keys()
	if len(keys) != 0 {
		t.Fatalf("Expected 0 keys, got %d", len(keys))
	}
}
