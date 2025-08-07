package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// DelCommand implements the DEL command
// DEL key [key ...]
// Removes the specified keys. A key is ignored if it does not exist.
// Returns the number of keys that were removed.
type DelCommand struct{}

func (c *DelCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) == 0 {
		return "", errors.New("wrong number of arguments for 'DEL'")
	}

	deleted := 0
	for _, key := range args {
		if store.Delete(key) {
			deleted++
		}
	}

	return protocol.Integer(deleted), nil
}

func init() {
	RegisterCommand("DEL", &DelCommand{})
}
