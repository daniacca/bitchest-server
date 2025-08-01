package commands

import (
	"strings"
	"testing"
	"time"

	"github.com/daniacca/bitchest/internal/db"
)

func TestMemoryStatsCommand(t *testing.T) {
	t.Run("empty database", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		result, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify initial stats
		if !strings.Contains(result, "keys=0") {
			t.Errorf("Expected keys=0, got: %s", result)
		}
		if !strings.Contains(result, "memory_usage=0") {
			t.Errorf("Expected memory_usage=0, got: %s", result)
		}
		if !strings.Contains(result, "data_size=0") {
			t.Errorf("Expected data_size=0, got: %s", result)
		}
	})

	t.Run("with arguments should fail", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		_, err := cmd.Execute([]string{"extra", "STATS"}, store)
		if err == nil {
			t.Error("Expected error for extra arguments, got nil")
		}
		if !strings.Contains(err.Error(), "wrong number of arguments") {
			t.Errorf("Expected 'wrong number of arguments' error, got: %v", err)
		}
	})

	t.Run("after adding keys", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		// Add some keys
		store.Set("key1", &db.StringValue{Val: "value1"})
		store.Set("key2", &db.StringValue{Val: "value2"})

		result, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify stats after adding keys
		if !strings.Contains(result, "keys=2") {
			t.Errorf("Expected keys=2, got: %s", result)
		}
		if !strings.Contains(result, "data_size=2") {
			t.Errorf("Expected data_size=2, got: %s", result)
		}
		// Memory usage should be > 0 (6 chars + 8 bytes per string = 14 bytes each = 28 total)
		if !strings.Contains(result, "memory_usage=") {
			t.Errorf("Expected memory_usage in result, got: %s", result)
		}
	})

	t.Run("after deleting keys", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		// Add and then delete a key
		store.Set("key1", &db.StringValue{Val: "value1"})
		store.Delete("key1")

		result, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify stats after deletion
		if !strings.Contains(result, "keys=0") {
			t.Errorf("Expected keys=0 after deletion, got: %s", result)
		}
		if !strings.Contains(result, "memory_usage=0") {
			t.Errorf("Expected memory_usage=0 after deletion, got: %s", result)
		}
	})

	t.Run("with expired keys", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		// Add a key with expiration in the past
		expiredTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
		store.Set("expired", &db.StringValue{Val: "value", ExpireAt: &expiredTime})

		_, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Try to get the expired key to trigger cleanup
		store.Get("expired")

		result2, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should show expired keys count
		if !strings.Contains(result2, "number_of_expired_keys=") {
			t.Errorf("Expected number_of_expired_keys in result, got: %s", result2)
		}
	})

	t.Run("peak memory tracking", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		// Add a key
		store.Set("key1", &db.StringValue{Val: "value1"})
		cmd.Execute([]string{"STATS"}, store)

		// Add another key
		store.Set("key2", &db.StringValue{Val: "value2"})
		result2, _ := cmd.Execute([]string{"STATS"}, store)

		// Delete a key
		store.Delete("key1")
		result3, _ := cmd.Execute([]string{"STATS"}, store)

		// Peak memory should be highest in result2
		// Extract peak memory values (simplified check)
		if !strings.Contains(result2, "peak_memory_usage=") {
			t.Errorf("Expected peak_memory_usage in result, got: %s", result2)
		}
		if !strings.Contains(result3, "peak_memory_usage=") {
			t.Errorf("Expected peak_memory_usage in result, got: %s", result3)
		}
	})

	t.Run("memory per key calculation", func(t *testing.T) {
		store := db.NewDB()
		cmd := &MemoryStatsCommand{}

		// Add keys with different sizes
		store.Set("key1", &db.StringValue{Val: "short"})
		store.Set("key2", &db.StringValue{Val: "longer_value"})

		result, err := cmd.Execute([]string{"STATS"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should show memory per key
		if !strings.Contains(result, "memory_per_key=") {
			t.Errorf("Expected memory_per_key in result, got: %s", result)
		}
	})
}