package commands

import (
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// TTLCommand returns the time to live for a key
type TTLCommand struct{}

func (c *TTLCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 1 {
		return protocol.Error("wrong number of arguments for 'TTL'"), nil
	}
	key := args[0]
	return protocol.Integer(store.GetTTL(key)), nil
}

func init() { RegisterCommand("TTL", &TTLCommand{}) }

func (c *TTLCommand) IsWrite() bool { return false } 