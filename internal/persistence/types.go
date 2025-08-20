package persistence

import "time"

// Entry represents a key-value entry in the store
type Entry struct {
	Key       string
	Value     []byte
	ExpiresAt *time.Time // nil if no TTL
}

// FsyncPolicy defines when to sync data to disk
type FsyncPolicy string

const (
	FsyncAlways   FsyncPolicy = "always"   // Sync after every write
	FsyncEverySec FsyncPolicy = "everysec" // Sync every second
	FsyncNo       FsyncPolicy = "no"       // Let OS handle syncing
) 