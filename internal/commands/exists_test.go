package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestExistsCommand(t *testing.T) {
	store := db.NewDB()
	store.Set("x", &db.StringValue{Val: "yes"})

	cmd := &ExistsCommand{}
	out, err := cmd.Execute([]string{"x", "y"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, ":1") {
		t.Errorf("Expected 1 key to exist, got %q", out)
	}
}

func TestExistsCommandWithNoInput(t *testing.T) {
	store := db.NewDB()
	cmd := &ExistsCommand{}
	out, err := cmd.Execute([]string{}, store)
	if err == nil {
		t.Errorf("Expected error, got %q", out)
	}
	if out != "" {
		t.Errorf("Expected empty response, got %q", out)
	}
}

func TestExistsCommandWithNonExistingKey(t *testing.T) {
	store := db.NewDB()
	store.Set("a", &db.StringValue{Val: "1"})

	cmd := &ExistsCommand{}
	out, err := cmd.Execute([]string{"b"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, ":0") {
		t.Errorf("Expected 0 keys to exist, got %q", out)
	}
}

func TestExistsIsReading(t *testing.T) {
	store := db.NewDB()
	store.Set("a", &db.StringValue{Val: "1"})

	cmd := &ExistsCommand{}
	isWrite := cmd.IsWrite()
	if isWrite {
		t.Errorf("Expected IsWrite to return false, got true")
	}
}