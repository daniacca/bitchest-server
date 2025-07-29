package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestKeysCommand(t *testing.T) {
	store := db.NewDB()
	cmd := &KeysCommand{}

	// Test with no arguments (valid)
	result, err := cmd.Execute([]string{}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "*0\r\n" {
		t.Errorf("Expected empty array for empty database, got %s", result)
	}

	// Test with arguments (invalid)
	_, err = cmd.Execute([]string{"pattern"}, store)
	if err == nil {
		t.Error("Expected error for wrong number of arguments, got nil")
	}
	if err.Error() != "wrong number of arguments for 'KEYS'" {
		t.Errorf("Expected specific error message, got %v", err)
	}

	// Test with keys in database
	store.Set("key1", &db.StringValue{Val: "value1"})
	store.Set("key2", &db.StringValue{Val: "value2"})
	store.Set("key3", &db.StringValue{Val: "value3"})

	result, err = cmd.Execute([]string{}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check that we get an array with 3 keys (order doesn't matter)
	if !strings.HasPrefix(result, "*3\r\n") {
		t.Errorf("Expected array with 3 keys, got %s", result)
	}
	
	// Check that all expected keys are present
	expectedKeys := []string{"key1", "key2", "key3"}
	for _, expectedKey := range expectedKeys {
		if !strings.Contains(result, expectedKey) {
			t.Errorf("Expected key '%s' not found in result: %s", expectedKey, result)
		}
	}
} 