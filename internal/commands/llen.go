package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LLenCommand implements the LLEN command
// LLEN key
// Returns the length of the list at the specified key.
type LLenCommand struct{}

func (c *LLenCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments for 'LLEN'")
	}

	key := args[0]

	val, ok := store.Get(key)
	if !ok {
		return protocol.Integer(0), nil
	}

	if list, ok := val.(*db.ListValue); ok {
		return protocol.Integer(list.Items.GetLength()), nil
	}

	return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
}

func init() {
	RegisterCommand("LLEN", &LLenCommand{})
}

func (c *LLenCommand) IsWrite() bool { return false }