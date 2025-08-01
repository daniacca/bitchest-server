package commands

import (
	"errors"
	"fmt"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

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