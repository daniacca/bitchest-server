package commands

import (
	"strings"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/parser"
)

// Command is the interface that each command implements
type Command interface {
	Execute(args []string, db *db.InMemoryDB) (string, error)
	IsWrite() bool
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

// ExecuteLine tokenizes a single-line command and executes it against the DB
// returning the output and the command handler (to inspect IsWrite)
func ExecuteLine(line string, store *db.InMemoryDB) (string, Command, error) {
	parts, err := parser.Tokenize(line)
	if err != nil {
		return "", nil, err
	}
	if len(parts) == 0 {
		return "", nil, nil
	}
	name := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	cmd, ok := ExtractCommand(name)
	if !ok {
		return "", nil, nil
	}
	out, execErr := cmd.Execute(args, store)
	return out, cmd, execErr
}
