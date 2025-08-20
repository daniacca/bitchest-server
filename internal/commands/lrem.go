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
	key := args[0]
	count, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("value is not an integer or out of range")
	}
	value := args[2]

	val, ok := store.Get(key)
	if !ok {
		return protocol.Integer(0), nil
	}

	list, ok := val.(*db.ListValue)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	removed := list.Items.Remove(value, count)
	return protocol.Integer(removed), nil
}

func init() { RegisterCommand("LREM", &LRemCommand{}) }

func (c *LRemCommand) IsWrite() bool { return true }