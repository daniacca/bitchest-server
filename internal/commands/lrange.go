package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LRangeCommand implements the LRANGE command
// LRANGE key start stop
// Returns the elements in the list at the specified key between the start and stop indices.
// The indices are zero-based and can be negative.
// If start is negative, it is counted from the end of the list.
// If stop is negative, it is counted from the end of the list.
// If start is greater than the length of the list, an empty array is returned.
// If stop is greater than the length of the list, it is set to the length of the list.
// If start is greater than stop, an empty array is returned.
type LRangeCommand struct{}

func (c *LRangeCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 3 {
		return "", errors.New("wrong number of arguments for 'LRANGE'")
	}

	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("invalid start index for 'LRANGE'")
	}
	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return "", errors.New("invalid stop index for 'LRANGE'")
	}

	val, ok := store.Get(key)
	if !ok {
		return protocol.Array([]string{}), nil
	}

	if list, ok := val.(*db.ListValue); ok {
		items := list.Items.GetItems()
		if start < 0 {
			start = len(items) + start
		}
		
		if stop < 0 {
			stop = len(items) + stop
		}
		
		if start >= len(items) {
			return protocol.Array([]string{}), nil
		}
		
		if stop >= len(items) {
			stop = len(items) - 1
		}
		
		if start > stop {
			return protocol.Array([]string{}), nil
		}

		return protocol.Array(items[start:stop+1]), nil
	}

	return "", errors.New("wrong type for 'LRANGE'")
}

func init() {
	RegisterCommand("LRANGE", &LRangeCommand{})
}