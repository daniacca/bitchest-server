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
