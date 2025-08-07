package commands

import (
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// FlushAllCommand implements the FLUSHALL command
// FLUSHALL
// Removes all keys from the database.
// Returns "OK".
type FlushAllCommand struct{}

func (c *FlushAllCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	store.FlushAll()
	return protocol.Simple("OK"), nil
}

func init() {
	RegisterCommand("FLUSHALL", &FlushAllCommand{})
}
