package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LSET key index value
// Sets the element at the given index in the list at the specified key to the given value.
// If the index is out of range, a null bulk is returned.
type LSetCommand struct{}

func (c *LSetCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 3 {
		return "", errors.New("wrong number of arguments for 'LSET'")
	}
	key := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("value is not an integer or out of range")
	}
	value := args[2]

	val, ok := store.Get(key)
	if !ok {
		return "", errors.New("no such key")
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if err := list.Items.Set(index, value); err != nil {
		return "", err
	}

	return protocol.Simple("OK"), nil
}

func init() { RegisterCommand("LSET", &LSetCommand{}) }

func (c *LSetCommand) IsWrite() bool { return true }