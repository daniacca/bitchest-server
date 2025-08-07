package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// KeysCommand implements the KEYS command
// KEYS
// Returns all keys in the database.
type KeysCommand struct{}

func (c *KeysCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 0 {
		return "", errors.New("wrong number of arguments for 'KEYS'")
	}

	keys := store.Keys()
	return protocol.Array(keys), nil
}

func init() {
	RegisterCommand("KEYS", &KeysCommand{})
}
