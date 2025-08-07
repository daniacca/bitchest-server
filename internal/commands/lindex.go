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
		return "", errors.New("invalid index for 'LINDEX'")
	}

	val, ok := store.Get(key)
	if !ok {
		return protocol.NullBulk(), nil
	}

	if list, ok := val.(*db.ListValue); ok {
		if index < 0 {
			index = list.Items.GetLength() + index
		}
		
		item, err := list.Items.Index(index)
		if err != nil {
			return protocol.NullBulk(), nil
		}

		return protocol.Bulk(item), nil
	}

	return "", errors.New("wrong type for 'LINDEX'")
}

func init() {
	RegisterCommand("LINDEX", &LIndexCommand{})
}