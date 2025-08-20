package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LINDEX key index
// Returns the element at the specified index in the list at the specified key.
// The index is zero-based and can be negative.
// If the index is out of range, a null bulk is returned.
type LIndexCommand struct{}


func (c *LIndexCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 2 {
		return "", errors.New("wrong number of arguments for 'LINDEX'")
	}

	key := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("value is not an integer or out of range")
	}

	val, ok := store.Get(key)
	if !ok {
		return protocol.NullBulk(), nil
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	item, err := list.Items.Index(index)
	if err != nil {
		return protocol.NullBulk(), nil
	}

	return protocol.Bulk(item), nil
}

func init() {
	RegisterCommand("LINDEX", &LIndexCommand{})
}

func (c *LIndexCommand) IsWrite() bool { return false }