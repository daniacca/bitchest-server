package commands

import (
	"errors"
	"strconv"
	"time"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// SetCommand sets a key to a string value with optional expiration
type SetCommand struct{}

func (c *SetCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) < 2 {
		return "", errors.New("wrong number of arguments for 'SET'")
	}

	key := args[0]
	value := args[1]
	
	// Create string value
	stringVal := &db.StringValue{Val: value}
	
	// Check for optional expiration
	if len(args) >= 4 && args[2] == "EX" {
		secondsStr := args[3]
		seconds, err := strconv.Atoi(secondsStr)
		if err != nil {
			return "", errors.New("invalid expiration time")
		}
		if seconds < 0 {
			return "", errors.New("expiration time must be non-negative")
		}
		
		// Set expiration time
		expireAt := time.Now().Add(time.Duration(seconds) * time.Second)
		stringVal.ExpireAt = &expireAt
	}

	store.Set(key, stringVal)
	return protocol.Simple("OK"), nil
}

func init() {
	RegisterCommand("SET", &SetCommand{})
}
