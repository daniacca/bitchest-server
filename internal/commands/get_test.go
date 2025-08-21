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

func TestGetCommand_NoArgs(t *testing.T) {
	store := db.NewDB()
	cmd := &GetCommand{}

	out, err := cmd.Execute([]string{}, store)
	if err == nil {
		t.Errorf("Expected error for no arguments, got %q", out)
	}
	if out != "" {
		t.Errorf("Expected empty response, got %q", out)
	}
}

func TestGetCommand_TooManyArgs(t *testing.T) {
	store := db.NewDB()
	store.Set("key", &db.StringValue{Val: "value"})

	cmd := &GetCommand{}
	out, err := cmd.Execute([]string{"key", "extra"}, store)
	if err == nil {
		t.Errorf("Expected error for too many arguments, got %q", out)
	}
	if out != "" {
		t.Errorf("Expected empty response, got %q", out)
	}
}

func TestGetCommand_IsReading(t *testing.T) {
	cmd := &GetCommand{}

	isWrite := cmd.IsWrite()
	if isWrite {
		t.Errorf("Expected IsWrite to return false, got true")
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

func (m *mockValue) Size() int {
	return 0
}
