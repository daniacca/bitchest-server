package commands

import (
	"errors"
	"strconv"
	"time"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// SetCommand sets a key to a string value with optional expiration and existence conditions
type SetCommand struct{}

type setOptions int

const (
	NoOption setOptions = iota
	NX
	XX
	EX
)

func (c *SetCommand) checkExists(key string, store *db.InMemoryDB) bool {
	_, exists := store.Get(key)
	return exists
}

func (c *SetCommand) parseOptions(args []string, startIndex int, store *db.InMemoryDB, key string) (*db.StringValue, error) {
	if startIndex >= len(args) {
		return &db.StringValue{Val: args[1]}, nil
	}

	stringVal := &db.StringValue{Val: args[1]}
	nxOrXxFound := false
	index := startIndex
	nxOrXxFailed := false

	for index < len(args) {
		if index >= len(args) {
			break
		}

		option := args[index]
		switch option {
		case "NX":
			if nxOrXxFound {
				return nil, errors.New("multiple options NX or XX found")
			}
			nxOrXxFound = true
			exists := c.checkExists(key, store)
			if exists {
				nxOrXxFailed = true
			}
			index++
		case "XX":
			if nxOrXxFound {
				return nil, errors.New("multiple options NX or XX found")
			}
			nxOrXxFound = true
			exists := c.checkExists(key, store)
			if !exists {
				nxOrXxFailed = true
			}
			index++
		case "EX":
			if index+1 >= len(args) {
				return nil, errors.New("missing expiration time")
			}
			secondsStr := args[index+1]
			seconds, err := strconv.Atoi(secondsStr)
			if err != nil {
				return nil, errors.New("invalid expiration time")
			}
			if seconds < 0 {
				return nil, errors.New("expiration time must be non-negative")
			}

			expireAt := time.Now().Add(time.Duration(seconds) * time.Second)
			stringVal.ExpireAt = &expireAt
			index += 2
		default:
			return nil, errors.New("invalid option: " + option)
		}
	}

	// If NX/XX condition failed, return nil
	if nxOrXxFailed {
		return nil, nil
	}

	return stringVal, nil
}

func (c *SetCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) < 2 {
		return "", errors.New("wrong number of arguments for 'SET'")
	}

	key := args[0]
	
	// Parse options and create string value
	stringVal, err := c.parseOptions(args, 2, store, key)
	if err != nil {
		return "", err
	}
	
	if stringVal == nil {
		return protocol.NullBulk(), nil
	}

	store.Set(key, stringVal)
	return protocol.Simple("OK"), nil
}

func init() {
	RegisterCommand("SET", &SetCommand{})
}
