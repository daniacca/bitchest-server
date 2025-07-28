package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestFlushAllCommand(t *testing.T) {
	store := db.NewDB()
	store.Set("k", &db.StringValue{Val: "v"})

	cmd := &FlushAllCommand{}
	out, err := cmd.Execute(nil, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK, got %q", out)
	}
	if len(store.Keys()) != 0 {
		t.Errorf("Expected empty store, got %d keys", len(store.Keys()))
	}
}
