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
