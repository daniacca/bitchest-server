package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

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
