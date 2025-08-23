package persistence

import (
	"context"
	"io"
	"testing"
)

type fakeStorageAdapter struct {
	manifest         *Manifest
	loadManifestErr  error
	listSegmentsErr  error
	openSegmentErr   error
	newSegmentErr    error
	saveManifestErr  error
	listSegments     []ObjectInfo
	openSegmentCalls int
	newSegmentCalls  int
	saveManifestCalls int
}

func (f *fakeStorageAdapter) LoadManifest(ctx context.Context) (*Manifest, error) {
	return f.manifest, f.loadManifestErr
}
func (f *fakeStorageAdapter) ListSegments(ctx context.Context, fromSeq uint64) ([]ObjectInfo, error) {
	return f.listSegments, f.listSegmentsErr
}
func (f *fakeStorageAdapter) OpenSegment(ctx context.Context, key string) (io.ReadCloser, error) {
	f.openSegmentCalls++
	return io.NopCloser(&fakeReader{}), f.openSegmentErr
}
func (f *fakeStorageAdapter) NewSegmentWriter(ctx context.Context, seq uint64) (SyncWriteCloser, string, error) {
	f.newSegmentCalls++
	return &fakeWriteCloser{}, "aof-key", f.newSegmentErr
}
func (f *fakeStorageAdapter) SaveManifest(ctx context.Context, m *Manifest) error {
	f.saveManifestCalls++
	return f.saveManifestErr
}
func (f *fakeStorageAdapter) OpenLatestSnapshot(ctx context.Context) (io.ReadCloser, *ObjectInfo, error) {
	return nil, nil, nil
}
func (f *fakeStorageAdapter) DeleteSegment(ctx context.Context, key string) error { return nil }
func (f *fakeStorageAdapter) DeleteSnapshot(ctx context.Context, key string) error { return nil }
func (f *fakeStorageAdapter) ListSnapshots(ctx context.Context) ([]ObjectInfo, error) { return nil, nil }
func (f *fakeStorageAdapter) NewSnapshotWriter(ctx context.Context, name string) (SyncWriteCloser, error) { return nil, nil }

type fakeReader struct{}
func (f *fakeReader) Read(p []byte) (int, error) { return 0, io.EOF }

type fakeWriteCloser struct{}
func (f *fakeWriteCloser) Write(p []byte) (int, error) { return len(p), nil }
func (f *fakeWriteCloser) Close() error                { return nil }
func (f *fakeWriteCloser) Sync() error                 { return nil }

func TestNewManager(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  60,
	}
	st := &fakeStorageAdapter{}

	mgr := NewManager(cfg, st)

	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}
	if mgr.cfg != cfg {
		t.Errorf("Expected cfg to be set, got %+v", mgr.cfg)
	}
	if mgr.st != st {
		t.Errorf("Expected st to be set")
	}
	if mgr.aof != nil {
		t.Errorf("Expected aof to be nil")
	}
	if mgr.snap != nil {
		t.Errorf("Expected snap to be nil")
	}
	if mgr.manifest != nil {
		t.Errorf("Expected manifest to be nil")
	}
}

func TestManager_Start_DisabledConfig(t *testing.T) {
	cfg := Config{
		Enabled:           false,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  60,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)

	err := mgr.Start(context.Background(), nil, nil)
	if err != nil {
		t.Errorf("Expected nil error when persistence is disabled, got %v", err)
	}
}

func TestManager_Start_ValidFlow(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		manifest: &Manifest{Version: 1, AOFSeq: 0},
		listSegments: []ObjectInfo{
			{Key: "seg1"},
			{Key: "seg2"},
		},
	}
	mgr := NewManager(cfg, fakeSt)
	applyCalled := 0
	apply := func(b []byte) error {
		applyCalled++
		return nil
	}
	err := mgr.Start(context.Background(), nil, apply)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if mgr.manifest == nil {
		t.Errorf("Expected manifest to be set")
	}
	if mgr.aof == nil {
		t.Errorf("Expected aof to be set")
	}
	if fakeSt.openSegmentCalls != 2 {
		t.Errorf("Expected openSegmentCalls=2, got %d", fakeSt.openSegmentCalls)
	}
	if fakeSt.newSegmentCalls != 1 {
		t.Errorf("Expected newSegmentCalls=1, got %d", fakeSt.newSegmentCalls)
	}
	if fakeSt.saveManifestCalls == 0 {
		t.Errorf("Expected SaveManifest to be called")
	}
}

func TestManager_Start_InvalidConfig(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   -1, // Invalid
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error for invalid config, got nil")
	}
}

func TestManager_Start_LoadManifestError(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		loadManifestErr: io.ErrUnexpectedEOF,
	}
	mgr := NewManager(cfg, fakeSt)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error from LoadManifest, got nil")
	}
}

func TestManager_Start_ListSegmentsError(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		manifest:        &Manifest{Version: 1, AOFSeq: 0},
		listSegmentsErr: io.ErrUnexpectedEOF,
	}
	mgr := NewManager(cfg, fakeSt)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error from ListSegments, got nil")
	}
}

func TestManager_Start_OpenSegmentError(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		manifest:     &Manifest{Version: 1, AOFSeq: 0},
		listSegments: []ObjectInfo{{Key: "seg1"}},
		openSegmentErr: io.ErrUnexpectedEOF,
	}
	mgr := NewManager(cfg, fakeSt)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error from OpenSegment, got nil")
	}
}

func TestManager_Start_NewSegmentWriterError(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		manifest:      &Manifest{Version: 1, AOFSeq: 0},
		newSegmentErr: io.ErrUnexpectedEOF,
	}
	mgr := NewManager(cfg, fakeSt)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error from NewSegmentWriter, got nil")
	}
}

func TestManager_Start_SaveManifestError(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	fakeSt := &fakeStorageAdapter{
		manifest:        &Manifest{Version: 1, AOFSeq: 0},
		saveManifestErr: io.ErrUnexpectedEOF,
	}
	
	mgr := NewManager(cfg, fakeSt)
	err := mgr.Start(context.Background(), nil, nil)
	if err == nil {
		t.Errorf("Expected error from SaveManifest, got nil")
	}
}

func TestManager_Stop_NoAOFOrSnap(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)
	
	// Both aof and snap are nil
	err := mgr.Stop(context.Background())
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

type fakeAOFWriter struct {
	syncCalled  bool
	closeCalled bool
}

func (f *fakeAOFWriter) Append([]byte) error { return nil }
func (f *fakeAOFWriter) Sync() error         { f.syncCalled = true; return nil }
func (f *fakeAOFWriter) Close() error        { f.closeCalled = true; return nil }
func (f *fakeAOFWriter) SizeMB() int64       { return 0 }

type fakeSnapshotter struct {
	stopCalled bool
}

func (f *fakeSnapshotter) Start(ctx context.Context, dumpFn func(SyncWriteCloser) error, onSuccess func(string)) {}

func (f *fakeSnapshotter) Stop() { f.stopCalled = true }

func TestManager_Stop_WithAOFAndSnap(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)
	fakeAOF := &fakeAOFWriter{}
	fakeSnap := &fakeSnapshotter{}
	mgr.aof = fakeAOF
	mgr.snap = fakeSnap

	err := mgr.Stop(context.Background())
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if !fakeAOF.syncCalled {
		t.Errorf("Expected Sync to be called on aof")
	}
	if !fakeAOF.closeCalled {
		t.Errorf("Expected Close to be called on aof")
	}
	if !fakeSnap.stopCalled {
		t.Errorf("Expected Stop to be called on snap")
	}
}

func TestManager_OnMutationRESP_NoAOF(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)
	mgr.aof = nil // aof is nil

	// Should not panic or call anything
	mgr.OnMutationRESP([]byte("SET foo bar"))
	// No assertion needed, just ensure no panic
}

type aofWriterMock struct {
	appendCalled bool
	sizeMB       int64
	rotateCalled bool
}

func (a *aofWriterMock) Append(cmd []byte) error {
	a.appendCalled = true
	return nil
}
func (a *aofWriterMock) Sync() error  { return nil }
func (a *aofWriterMock) Close() error { return nil }
func (a *aofWriterMock) SizeMB() int64 {
	return a.sizeMB
}

func TestManager_OnMutationRESP_Appends(t *testing.T) {
	cfg := Config{
		Enabled:           true,
		AOFMaxSegmentMB:   10,
		AppendFsyncPolicy: "always",
		SnapshotInterval:  0,
	}
	st := &fakeStorageAdapter{}
	mgr := NewManager(cfg, st)
	mockAOF := &aofWriterMock{sizeMB: 0}
	mgr.aof = mockAOF

	mgr.OnMutationRESP([]byte("SET foo bar"))

	if !mockAOF.appendCalled {
		t.Errorf("Expected Append to be called")
	}
}
