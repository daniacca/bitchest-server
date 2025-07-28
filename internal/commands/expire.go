package commands

import (
	"errors"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// ExpireCommand sets an expiration time for a key
type ExpireCommand struct{}

func (c *ExpireCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 2 {
		return "", errors.New("wrong number of arguments for 'EXPIRE'")
	}

	key := args[0]
	secondsStr := args[1]

	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		return "", errors.New("invalid expiration time")
	}

	if seconds < 0 {
		return "", errors.New("expiration time must be non-negative")
	}

	success := store.SetExpiration(key, seconds)
	if success {
		return protocol.Integer(1), nil
	}
	return protocol.Integer(0), nil
}

func init() {
	RegisterCommand("EXPIRE", &ExpireCommand{})
} 