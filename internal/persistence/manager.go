// persistence/manager.go
package persistence

import (
	"context"
	"time"
)

type Store interface {
	// For snapshot: only string values are exposed; persistence handles serialization
	RangeStrings(func(key string, value string, expiresAt *time.Time) bool)
	// For snapshot: list values
	RangeLists(func(key string, items []string, expiresAt *time.Time) bool)
}

type Manager struct {
	cfg      Config
	st       StorageAdapter
	aof      *aofWriter
	snap     *snapshotter
	manifest *Manifest
}

func NewManager(cfg Config, st StorageAdapter) *Manager {
	return &Manager{
		cfg: cfg, st: st,
	}
}

func (m *Manager) Start(ctx context.Context, db Store, apply func([]byte) error) error {
	if !m.cfg.Enabled {
		return nil
	}

	// Validate configuration
	if err := m.validateConfig(); err != nil {
		return err
	}

	// 1) Load manifest
	man, err := m.st.LoadManifest(ctx)
	if err != nil {
		return err
	}
	if man == nil {
		man = &Manifest{Version: 1}
	}
	m.manifest = man

	// 2) Load latest snapshot if available
	if man.LastSnapshot != "" {
		r, _, err := m.st.OpenLatestSnapshot(ctx)
		if err != nil {
			// Log error but continue - snapshot loading is best effort
			// You might want to add logging here
		} else if r != nil {
			if err := restoreSnapshot(r, apply); err != nil {
				r.Close()
				// Log error but continue - snapshot loading is best effort
			} else {
				r.Close()
			}
		}
	}

	// 3) Replay AOF segments
	segs, err := m.st.ListSegments(ctx, man.AOFSeq)
	if err != nil {
		return err
	}
	for _, s := range segs {
		r, err := m.st.OpenSegment(ctx, s.Key)
		if err != nil {
			return err
		}
		if err := replayAOF(r, apply); err != nil {
			r.Close()
			return err
		}
		r.Close()
	}

	// 4) Open new AOF writer for current segment
	aw, key, err := m.st.NewSegmentWriter(ctx, man.AOFSeq+1)
	if err != nil {
		return err
	}
	m.aof = newAOFWriter(aw, m.cfg.AppendFsyncPolicy)
	m.manifest.AOFSeq = man.AOFSeq + 1
	m.manifest.CurrentAOFKey = key
	if err := m.st.SaveManifest(ctx, m.manifest); err != nil {
		return err
	}

	// 5) Start snapshot ticker (background goroutine)
	if m.cfg.SnapshotInterval > 0 {
		m.snap = newSnapshotter(m.st, m.cfg.SnapshotInterval)
		m.snap.Start(ctx, func(w SyncWriteCloser) error {
			return dumpSnapshot(w, db)
		}, func(snapKey string) {
			// Update manifest on successful snapshot
			m.manifest.LastSnapshot = snapKey
			// Best effort - don't fail if manifest save fails
			_ = m.st.SaveManifest(ctx, m.manifest)
		})
	}

	return nil
}

// validateConfig validates the persistence configuration
func (m *Manager) validateConfig() error {
	if m.cfg.AOFMaxSegmentMB < 0 {
		return ErrInvalidConfig
	}
	
	switch m.cfg.AppendFsyncPolicy {
	case "always", "everysec", "no":
		// Valid policies
	default:
		return ErrInvalidConfig
	}
	
	return nil
}

func (m *Manager) OnMutationRESP(cmd []byte) {
	if m.aof == nil { return }
	_ = m.aof.Append(cmd) // gestisce buffer + fsync policy
	// rotazione dimensione
	if m.aof.SizeMB() >= int64(m.cfg.AOFMaxSegmentMB) && m.cfg.AOFMaxSegmentMB > 0 {
		_ = m.rotateSegment(context.Background())
	}
}

func (m *Manager) rotateSegment(ctx context.Context) error {
	_ = m.aof.Close()
	aw, key, err := m.st.NewSegmentWriter(ctx, m.manifest.AOFSeq+1)
	if err != nil { return err }
	m.aof = newAOFWriter(aw, m.cfg.AppendFsyncPolicy)
	m.manifest.AOFSeq++
	m.manifest.CurrentAOFKey = key
	return m.st.SaveManifest(ctx, m.manifest)
}

func (m *Manager) Stop(ctx context.Context) error {
	if m.snap != nil { m.snap.Stop() }
	if m.aof != nil { _ = m.aof.Sync(); _ = m.aof.Close() }
	return nil
}
