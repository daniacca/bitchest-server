package commands

import (
	"strings"

	"github.com/daniacca/bitchest/internal/db"
)

// Command is the interface that each command implements
type Command interface {
	Execute(args []string, db *db.InMemoryDB) (string, error)
}

// CommandsRegistry maps the command names (e.g. "SET") to their respective handlers
var CommandsRegistry = map[string]Command{}

// RegisterCommand registers a new command
func RegisterCommand(name string, cmd Command) {
	CommandsRegistry[strings.ToUpper(name)] = cmd
}

// ExtractCommand search and return the command handler for the given name
func ExtractCommand(name string) (Command, bool) {
	cmd, ok := CommandsRegistry[strings.ToUpper(name)]
	return cmd, ok
}
