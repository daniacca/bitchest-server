package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// ExistsCommand implements the EXISTS command
// EXISTS key [key ...]
// Returns 1 if the key exists, otherwise returns 0.
// If multiple keys are provided, returns 1 if at least one key exists, otherwise returns 0.
type ExistsCommand struct{}

func (c *ExistsCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) == 0 {
		return "", errors.New("wrong number of arguments for 'EXISTS'")
	}

	count := 0
	for _, key := range args {
		if _, ok := store.Get(key); ok {
			count++
		}
	}

	return protocol.Integer(count), nil
}

func init() {
	RegisterCommand("EXISTS", &ExistsCommand{})
}
