package commands

import (
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
	if len(args) != 1 {
		return protocol.Error("wrong number of arguments for 'LPOP'"), nil
	}
	key := args[0]

	val, ok := store.Get(key)
	if !ok {
		return protocol.NullBulk(), nil
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value"), nil
	}

	item, err := list.Items.Shift()
	if err != nil {
		return protocol.NullBulk(), nil
	}
	return protocol.Bulk(item), nil
}

func init() { RegisterCommand("LPOP", &LPopCommand{}) }

func (c *LPopCommand) IsWrite() bool { return true }