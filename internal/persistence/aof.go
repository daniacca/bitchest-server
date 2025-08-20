package persistence

import (
	"bufio"
	"sync"
	"time"
)

type aofWriter struct {
	writer   SyncWriteCloser
	buffer   *bufio.Writer
	policy   FsyncPolicy
	size     int64
	mu       sync.Mutex
	stopChan chan struct{}
	ticker   *time.Ticker
}

func newAOFWriter(writer SyncWriteCloser, policy string) *aofWriter {
	aof := &aofWriter{
		writer:   writer,
		buffer:   bufio.NewWriter(writer),
		policy:   FsyncPolicy(policy),
		stopChan: make(chan struct{}),
	}

	// Start periodic sync if needed
	if aof.policy == FsyncEverySec {
		aof.ticker = time.NewTicker(time.Second)
		go aof.periodicSync()
	}

	return aof
}

func (a *aofWriter) Append(cmd []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Write command length as 4-byte big-endian
	length := uint32(len(cmd))
	header := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}

	if _, err := a.buffer.Write(header); err != nil {
		return err
	}

	if _, err := a.buffer.Write(cmd); err != nil {
		return err
	}

	a.size += int64(len(header) + len(cmd))

	// Sync immediately if policy is "always"
	if a.policy == FsyncAlways {
		return a.buffer.Flush()
	}

	return nil
}

func (a *aofWriter) SizeMB() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.size / (1024 * 1024)
}

func (a *aofWriter) Sync() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.buffer.Flush(); err != nil {
		return err
	}

	return a.writer.Sync()
}

func (a *aofWriter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Stop periodic sync
	if a.ticker != nil {
		a.ticker.Stop()
		close(a.stopChan)
	}

	if err := a.buffer.Flush(); err != nil {
		return err
	}

	return a.writer.Close()
}

func (a *aofWriter) periodicSync() {
	for {
		select {
		case <-a.ticker.C:
			_ = a.Sync()
		case <-a.stopChan:
			return
		}
	}
} 