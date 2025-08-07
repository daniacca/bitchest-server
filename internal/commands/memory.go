package commands

import (
	"errors"
	"fmt"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

// MemoryStatsCommand implements the MEMORY STATS command
// MEMORY STATS
// Returns memory statistics for the database.
// Returns a map with the following keys:
// - keys: number of keys in the database
// - memory_usage: total memory usage of the database
// - memory_per_key: average memory usage per key
// - peak_memory_usage: peak memory usage of the database
// - number_of_expired_keys: number of expired keys
// - data_size: total data size of the database
type MemoryStatsCommand struct{}

func (c *MemoryStatsCommand) Execute(args []string, store *db.InMemoryDB) (string, error) {
	if len(args) != 1 || args[0] != "STATS" {
		return "", errors.New("wrong number of arguments for 'MEMORY STATS'")
	}

	stats := store.GetStats()

	return protocol.Array([]string{
		fmt.Sprintf("keys=%d", stats.Keys),
		fmt.Sprintf("memory_usage=%d", stats.MemoryUsage),
		fmt.Sprintf("memory_per_key=%d", stats.MemoryPerKey),
		fmt.Sprintf("peak_memory_usage=%d", stats.PeakMemoryUsage),
		fmt.Sprintf("number_of_expired_keys=%d", stats.NumberOfExpiredKeys),
		fmt.Sprintf("data_size=%d", stats.DataSize),
	}), nil
}

func init() {
	RegisterCommand("MEMORY", &MemoryStatsCommand{})
}