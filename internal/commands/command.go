package commands

import (
	"strings"

	"github.com/daniacca/bitchest/internal/db"
)

// Command is the interface that each command implements
type Command interface {
	Execute(args []string, db *db.InMemoryDB) (string, error)
}

// Registry maps the command names (e.g. "SET") to their respective handlers
var Registry = map[string]Command{}

// RegisterCommand registers a new command
func RegisterCommand(name string, cmd Command) {
	Registry[strings.ToUpper(name)] = cmd
}

func ExtractCommand(name string) (Command, bool) {
	cmd, ok := Registry[strings.ToUpper(name)]
	return cmd, ok
}
