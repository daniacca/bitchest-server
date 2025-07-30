package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestDelCommand(t *testing.T) {
	store := db.NewDB()
	store.Set("a", &db.StringValue{Val: "1"})
	store.Set("b", &db.StringValue{Val: "2"})

	cmd := &DelCommand{}
	out, err := cmd.Execute([]string{"a", "b", "c"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, ":2") {
		t.Errorf("Expected 2 keys deleted, got %q", out)
	}
}

func TestDelWithNoInput(t *testing.T) {
	store := db.NewDB()
	cmd := &DelCommand{}
	out, err := cmd.Execute([]string{}, store)
	if err == nil {
		t.Errorf("Expected error, got %q", out)
	}
	if out != "" {
		t.Errorf("Expected empty response, got %q", out)
	}
}

func TestDelWithNonExistingKey(t *testing.T) {
	store := db.NewDB()
	store.Set("a", &db.StringValue{Val: "1"})

	cmd := &DelCommand{}
	out, err := cmd.Execute([]string{"b"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, ":0") {
		t.Errorf("Expected 0 keys deleted, got %q", out)
	}
}