package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Snapshotter interface {
    Start(ctx context.Context, dumpFn func(SyncWriteCloser) error, onSuccess func(string))
    Stop()
}

type snapshotter struct {
	storage  StorageAdapter
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func newSnapshotter(storage StorageAdapter, interval time.Duration) *snapshotter {
	return &snapshotter{
		storage:  storage,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (s *snapshotter) Start(ctx context.Context, dumpFn func(SyncWriteCloser) error, onSuccess func(string)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.run(ctx, dumpFn, onSuccess)
	}()
}

func (s *snapshotter) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *snapshotter) run(ctx context.Context, dumpFn func(SyncWriteCloser) error, onSuccess func(string)) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.createSnapshot(ctx, dumpFn, onSuccess)
		}
	}
}

func (s *snapshotter) createSnapshot(ctx context.Context, dumpFn func(SyncWriteCloser) error, onSuccess func(string)) {
	// Generate snapshot name with timestamp
	name := fmt.Sprintf("snapshot-%d.bin", time.Now().Unix())
	
	writer, err := s.storage.NewSnapshotWriter(ctx, name)
	if err != nil {
		// Log error but don't fail - snapshots are best effort
		return
	}

	if err := dumpFn(writer); err != nil {
		_ = writer.Close()
		// Log error but don't fail
		return
	}

	if err := writer.Close(); err != nil {
		// Log error but don't fail
		return
	}

	// Notify success
	onSuccess(name)
}
