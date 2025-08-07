package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// RPopCommand implements the RPOP command
// RPOP key [count]
// Removes and returns the last element of the list at the specified key.
// If the key doesn't exist, it returns an error.
// If the list is empty, it returns an error.
// If count is provided, it removes and returns the last count elements of the list.
// Returns a nil reply if the list is empty.
// Returns an array of elements if count is provided.
// Returns a bulk string with the last element if count is not provided.
type RPopCommand struct{}

func (c *RPopCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) > 2 || len(args) == 0 {
		return "", errors.New("wrong number of arguments for 'RPOP'")
	}

	key := args[0]
	count := 1
	optionalCount := false

	if len(args) == 2 {
		argsCount, err := strconv.Atoi(args[1])
		if err != nil {
			return "", errors.New("invalid count for 'RPOP'")
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
			item, err := list.Items.Pop()
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

	return "", errors.New("wrong type for 'RPOP'")
}

func init() {
	RegisterCommand("RPOP", &RPopCommand{})
}