package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LPopCommand implements the LPOP command
// LPOP key [count]
// Removes and returns the first element of the list at the specified key.
// If the key doesn't exist, it returns an error.
// If the list is empty, it returns an error.
// If count is provided, it removes and returns the first count elements of the list.
// Returns a nil reply if the list is empty.
// Returns an array of elements if count is provided.
// Returns a bulk string with the first element if count is not provided.
type LPopCommand struct{}

func (c *LPopCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) > 2 || len(args) == 0 {
		return "", errors.New("wrong number of arguments for 'LPOP'")
	}

	key := args[0]
	count := 1
	optionalCount := false

	if len(args) == 2 {
		argsCount, err := strconv.Atoi(args[1])
		if err != nil {
			return "", errors.New("invalid count for 'LPOP'")
		}
		count = argsCount
		optionalCount = true
	}

	val, ok := store.Get(key)
	if !ok {
		return protocol.NullBulk(), nil
	}

	if list, ok := val.(*db.ListValue); ok {
		items := []string{}
		for i := 0; i < count; i++ {
			item, err := list.Items.Shift()
			if err != nil {
				break
			}
			items = append(items, item)
		}

		if len(items) == 0 {
			return protocol.NullBulk(), nil
		}

		store.Set(key, list)
		if optionalCount {
			return protocol.Array(items), nil
		}
		return protocol.Bulk(items[0]), nil
	}

	return "", errors.New("wrong type for 'LPOP'")
}

func init() {
	RegisterCommand("LPOP", &LPopCommand{})
}