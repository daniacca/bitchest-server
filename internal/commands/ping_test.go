package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestPingCommand(t *testing.T) {
	store := db.NewDB()
	cmd := &PingCommand{}

	out, err := cmd.Execute(nil, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if out != "+PONG\r\n" {
		t.Errorf("Expected PONG, got %q", out)
	}
}