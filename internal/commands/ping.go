package commands

import (
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// PingCommand implements the PING command
// PING
// Returns "PONG".
type PingCommand struct{}

func (c *PingCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	return protocol.Simple("PONG"), nil
}

func init() {
	RegisterCommand("PING", &PingCommand{})
}