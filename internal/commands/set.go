package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// SetCommand imposta una chiave a un valore stringa
type SetCommand struct{}

func (c *SetCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 2 {
		return "", errors.New("wrong number of arguments for 'SET'")
	}

	key := args[0]
	value := args[1]

	store.Set(key, &db.StringValue{Val: value})
	return protocol.Simple("OK"), nil
}

func init() {
	RegisterCommand("SET", &SetCommand{})
}
