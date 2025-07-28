package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

type GetCommand struct{}

func (c *GetCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 1 {
		return "", errors.New("wrong number of arguments for 'GET'")
	}

	key := args[0]

	val, ok := store.Get(key)
	if !ok {
		return protocol.NullBulk(), nil
	}

	strVal, ok := val.(*db.StringValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return protocol.Bulk(strVal.Get()), nil
}

func init() {
	RegisterCommand("GET", &GetCommand{})
}
