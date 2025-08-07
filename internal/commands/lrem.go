package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// LRemCommand implements the LREM command, which removes elements equal to a given value from a list stored at a key.
// The count argument influences the number and direction of removals:
//   - count > 0: Remove up to count occurrences from head to tail.
//   - count < 0: Remove up to |count| occurrences from tail to head.
//   - count == 0: Remove all occurrences.
// Returns the number of removed elements as an integer reply.
// If the key does not exist, returns 0.
// Returns an error if the key exists but is not a list, or if the count is not a valid integer.
type LRemCommand struct{}

func (c *LRemCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
    if len(args) != 3 {
        return "", errors.New("wrong number of arguments for 'LREM'")
    }
    
	key, countStr, value := args[0], args[1], args[2]
    count, err := strconv.Atoi(countStr)
    if err != nil {
        return "", errors.New("invalid count for 'LREM'")
    }
    
	val, ok := store.Get(key)
    if !ok {
        return protocol.Integer(0), nil
    }
    
	list, ok := val.(*db.ListValue)
    if !ok {
        return "", errors.New("wrong type for 'LREM'")
    }
	
    removed := list.Items.Remove(value, count)
    store.Set(key, list)
    return protocol.Integer(removed), nil
}

func init() {
    RegisterCommand("LREM", &LRemCommand{})
}