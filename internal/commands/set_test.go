package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestSetCommand_Valid(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	out, err := cmd.Execute([]string{"foo", "bar"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK response, got %q", out)
	}

	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key not set in store")
	}
	if str, ok := val.(*db.StringValue); !ok || str.Get() != "bar" {
		t.Errorf("Expected 'bar', got %v", val)
	}
}

func TestSetCommand_NoValueProvided(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	_, err := cmd.Execute([]string{"foo"}, store)
	if err == nil {
		t.Fatal("Expected error on wrong number of args")
	}
}

func TestSetCommand_TooManyArgs(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	_, err := cmd.Execute([]string{"foo", "bar", "baz"}, store)
	if err == nil {
		t.Fatal("Expected error on wrong number of args")
	}
}