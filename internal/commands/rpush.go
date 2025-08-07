package commands

import (
	"errors"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// RPushCommand implements the RPUSH command
// RPUSH key value1 [value2 ...]
// Adds the value to the tail of the list at the specified key.
// If the key doesn't exist, it creates a new list with the value.
// If the key exists, it adds the value to the tail of the list.
// It is possible to add multiple values at once. Elements are
// inserted one after the other to the tail of the list. So for
// instance, the command RPUSH mylist a b c will result into a
// list containing a as first element, b as second element and
// c as third element.
// Returns the number of items in the list after the operation.
type RPushCommand struct{}

func (c *RPushCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) < 2 {
		return "", errors.New("wrong number of arguments for 'RPUSH'")
	}

	key := args[0]
	values := args[1:]

	val, ok := store.Get(key)
	if !ok {
		// If the key doesn't exist, create a new list
		list := &db.ListValue{
			Items: db.Queue{},
		}

		for _, value := range values {
			list.Items.Push(value)
		}

		store.Set(key, list)
		return protocol.Integer(len(values)), nil
	}

	if list, ok := val.(*db.ListValue); ok {
		// If the key exists, add the value to the list
		for _, value := range values {
			list.Items.Push(value)
		}
		store.Set(key, list)
		return protocol.Integer(list.Items.GetLength()), nil
	}

	return "", errors.New("wrong type for 'RPUSH'")
}

func init() {
	RegisterCommand("RPUSH", &RPushCommand{})
}