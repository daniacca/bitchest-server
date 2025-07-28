package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// TTLCommand returns the time to live for a key
type TTLCommand struct{}

func (c *TTLCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments for 'TTL'")
	}

	key := args[0]
	ttl := store.GetTTL(key)
	return protocol.Integer(ttl), nil
}

func init() {
	RegisterCommand("TTL", &TTLCommand{})
} 