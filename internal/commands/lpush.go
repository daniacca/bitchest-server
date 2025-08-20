package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LPushCommand implements the LPUSH command
// LPUSH key value1 [value2 ...]
// Adds the value to the head of the list at the specified key.
// If the key doesn't exist, it creates a new list with the value.
// If the key exists, it adds the value to the head of the list.
// It is possible to add multiple values at once. Elements are
// inserted one after the other to the head of the list, from
// the leftmost element to the rightmost element. So for instance,
// the command LPUSH mylist a b c will result into a list containing
// c as first element, b as second element and a as third element.
// Returns the number of items in the list after the operation.
type LPushCommand struct{}

func (c *LPushCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) < 2 {
		return "", errors.New("wrong number of arguments for 'LPUSH'")
	}

	key := args[0]
	values := args[1:]

	val, ok := store.Get(key)
	if !ok {
		// Create new list
		l := &db.ListValue{Items: *db.NewQueue()}
		for _, v := range values {
			l.Items.Unshift(v)
		}
		store.Set(key, l)
		return protocol.Integer(l.Items.GetLength()), nil
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	for _, v := range values {
		list.Items.Unshift(v)
	}
	return protocol.Integer(list.Items.GetLength()), nil
}

func init() { RegisterCommand("LPUSH", &LPushCommand{}) }

func (c *LPushCommand) IsWrite() bool { return true }