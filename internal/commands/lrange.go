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
		return "", errors.New("value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return "", errors.New("value is not an integer or out of range")
	}

	val, ok := store.Get(key)
	if !ok {
		return protocol.Array([]string{}), nil
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	items := list.Items.GetItems()
	n := len(items)
	if start < 0 { start = n + start }
	if stop < 0 { stop = n + stop }
	if start < 0 { start = 0 }
	if stop >= n { stop = n - 1 }
	if start > stop || start >= n { return protocol.Array([]string{}), nil }
	return protocol.Array(items[start : stop+1]), nil
}

func init() {
	RegisterCommand("LRANGE", &LRangeCommand{})
}

func (c *LRangeCommand) IsWrite() bool { return false }