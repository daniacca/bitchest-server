package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestGetCommand_ExistingKey(t *testing.T) {
	store := db.NewDB()
	store.Set("key", &db.StringValue{Val: "value"})

	cmd := &GetCommand{}
	out, err := cmd.Execute([]string{"key"}, store)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, "value") {
		t.Errorf("Expected 'value' in response, got %q", out)
	}
}

func TestGetCommand_MissingKey(t *testing.T) {
	store := db.NewDB()
	cmd := &GetCommand{}

	out, err := cmd.Execute([]string{"not_exist"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if out != "$-1\r\n" {
		t.Errorf("Expected null bulk, got %q", out)
	}
}

func TestGetCommand_WrongType(t *testing.T) {
	store := db.NewDB()
	store.Set("foo", &mockValue{})

	cmd := &GetCommand{}
	_, err := cmd.Execute([]string{"foo"}, store)
	if err == nil {
		t.Fatal("Expected type error")
	}
}

// mock value to force type mismatch
type mockValue struct{}

func (m *mockValue) Type() db.ValueType {
	return "mock"
}

func (m *mockValue) IsExpired() bool {
	return false
}
