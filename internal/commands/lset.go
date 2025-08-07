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
    
	key, idxStr, value := args[0], args[1], args[2]
    
	idx, err := strconv.Atoi(idxStr)
    if err != nil {
        return "", errors.New("invalid index for 'LSET'")
    }
    
	val, ok := store.Get(key)
    if !ok {
        return protocol.NullBulk(), nil
    }
    
	list, ok := val.(*db.ListValue)
    if !ok {
        return "", errors.New("wrong type for 'LSET'")
    }
    
	if idx < 0 {
		idx = list.Items.GetLength() + idx
	}

	if err := list.Items.Set(idx, value); err != nil {
        return "", errors.New("index out of range")
    }
    
	store.Set(key, list)
    return protocol.Bulk("OK"), nil
}

func init() {
    RegisterCommand("LSET", &LSetCommand{})
}