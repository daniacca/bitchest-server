package persistence

import (
	"context"
	"io"
	"time"
)

type SyncWriteCloser interface {
	io.WriteCloser
	Sync() error
}

// Minimal object info on saved entries (file/key)
type ObjectInfo struct {
	Key       string
	Size      int64
	ModTime   time.Time
	Sequence  uint64 // for identify AOF segment
}

// Interface that must be implemented for each adapter (FS, S3, MinIO, …)
type StorageAdapter interface {
	// Manifest
	LoadManifest(ctx context.Context) (*Manifest, error)
	SaveManifest(ctx context.Context, m *Manifest) error

	// Snapshot
	NewSnapshotWriter(ctx context.Context, name string) (SyncWriteCloser, error)
	OpenLatestSnapshot(ctx context.Context) (io.ReadCloser, *ObjectInfo, error)
	ListSnapshots(ctx context.Context) ([]ObjectInfo, error)
	DeleteSnapshot(ctx context.Context, key string) error

	// AOF segments (append-only by segment)
	NewSegmentWriter(ctx context.Context, seq uint64) (SyncWriteCloser, string /*key*/, error)
	ListSegments(ctx context.Context, fromSeq uint64) ([]ObjectInfo, error)
	OpenSegment(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteSegment(ctx context.Context, key string) error
}
