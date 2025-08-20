package fsadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/daniacca/bitchest/internal/persistence"
)

type FSAdapter struct {
	base string // data dir
}

func New(base string) *FSAdapter { return &FSAdapter{base: base} }

func (fs *FSAdapter) manifestPath() string  { return filepath.Join(fs.base, "manifest.json") }
func (fs *FSAdapter) aofDir() string        { return filepath.Join(fs.base, "aof") }
func (fs *FSAdapter) snapsDir() string      { return filepath.Join(fs.base, "snapshots") }

// Utility private method that ensure the presence of all the necessary directory
func (fs *FSAdapter) ensureDirs() error {
	for _, d := range []string{fs.base, fs.aofDir(), fs.snapsDir()} {
		if err := os.MkdirAll(d, 0o755); err != nil { return err }
	}
	return nil
}

// Load the Manifest file from FS
func (fs *FSAdapter) LoadManifest(ctx context.Context) (*persistence.Manifest, error) {
	_ = fs.ensureDirs()
	b, err := os.ReadFile(fs.manifestPath())
	if err != nil { return nil, nil } // first run - no manifest file created yet
	var m persistence.Manifest
	if parseError := json.Unmarshal(b, &m); parseError != nil { return nil, parseError } // error parsing JSON data
	return &m, nil // return the parsed manifest
}

// Persist the Manifest file from FS. Save into a tmp file, and if everything goes
// well it will be moved to the Manifest right path.
func (fs *FSAdapter) SaveManifest(ctx context.Context, m *persistence.Manifest) error {
	_ = fs.ensureDirs()
	tmp := fs.manifestPath() + ".tmp"
	b, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(tmp, b, 0o644); err != nil { return err }
	return os.Rename(tmp, fs.manifestPath())
}

// Create a new Snapshot Writer
func (fs *FSAdapter) NewSnapshotWriter(ctx context.Context, name string) (persistence.SyncWriteCloser, error) {
	_ = fs.ensureDirs()
	p := filepath.Join(fs.snapsDir(), name)
	f, err := os.Create(p)
	if err != nil { return nil, err }
	return fileSyncWriteCloser{File: f}, nil
}

// Open Lastest Snapshot persisted
// This method use os package (ReadDir) to list all entries inside the snapshot directory
// it then loop over all the file and take the file with the "newer" ModTime.
// Be carefull, the execution time will scale linearly with the number of the snapshot persisted.
func (fs *FSAdapter) OpenLatestSnapshot(ctx context.Context) (io.ReadCloser, *persistence.ObjectInfo, error) {
	entries, err := os.ReadDir(fs.snapsDir())
	if err != nil { return nil, nil, err }
	var latest string
	var latestInfo os.FileInfo
	for _, e := range entries {
		if e.IsDir() { continue }
		fi, _ := e.Info()
		if latestInfo == nil || fi.ModTime().After(latestInfo.ModTime()) {
			latest = filepath.Join(fs.snapsDir(), e.Name())
			latestInfo = fi
		}
	}
	if latest == "" { return nil, nil, nil }
	f, err := os.Open(latest)
	if err != nil { return nil, nil, err }
	return f, &persistence.ObjectInfo{
		Key: latest, Size: latestInfo.Size(), ModTime: latestInfo.ModTime(),
	}, nil
}

// Return the List of all the snapshots actually saved on the snap directory
// The output is returned in the same order the file are retrieved from FS
func (fs *FSAdapter) ListSnapshots(ctx context.Context) ([]persistence.ObjectInfo, error) {
	entries, err := os.ReadDir(fs.snapsDir())
	if err != nil { return nil, err }
	var out []persistence.ObjectInfo
	for _, e := range entries {
		if e.IsDir() { continue }
		fi, _ := e.Info()
		out = append(out, persistence.ObjectInfo{
			Key: filepath.Join(fs.snapsDir(), e.Name()),
			Size: fi.Size(),
			ModTime: fi.ModTime(),
		})
	}
	return out, nil
}

// Delete the snapshot specified with the key string
func (fs *FSAdapter) DeleteSnapshot(ctx context.Context, key string) error {
	return os.Remove(key)
}

// Create a new Segment Writer for AOF files
func (fs *FSAdapter) NewSegmentWriter(ctx context.Context, seq uint64) (persistence.SyncWriteCloser, string, error) {
	_ = fs.ensureDirs()
	name := filepath.Join(fs.aofDir(), fmt.Sprintf("%06d.aof", seq))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil { return nil, "", err }
	return fileSyncWriteCloser{File: f}, name, nil
}

// This method return a list of AOF segements that are already saved on file system.
// It work as same as ListSnapshot, but for AOF segments
// The output is ordered in ascending order. 
func (fs *FSAdapter) ListSegments(ctx context.Context, fromSeq uint64) ([]persistence.ObjectInfo, error) {
	entries, err := os.ReadDir(fs.aofDir())
	if err != nil { return nil, err }
	var out []persistence.ObjectInfo
	for _, e := range entries {
		if e.IsDir() { continue }
		fi, _ := e.Info()
		seq := parseSeq(e.Name()) // es. 000123.aof
		if seq >= fromSeq {
			out = append(out, persistence.ObjectInfo{
				Key: filepath.Join(fs.aofDir(), e.Name()),
				Size: fi.Size(),
				ModTime: fi.ModTime(),
				Sequence: seq,
			})
		}
	}

	// order by ascending sequence
	sort.Slice(out, func(i,j int) bool { return out[i].Sequence < out[j].Sequence })
	return out, nil
}

func (fs *FSAdapter) OpenSegment(ctx context.Context, key string) (io.ReadCloser, error) {
	return os.Open(key)
}

func (fs *FSAdapter) DeleteSegment(ctx context.Context, key string) error {
	return os.Remove(key)
}

type fileSyncWriteCloser struct{ *os.File }
func (f fileSyncWriteCloser) Sync() error { return f.File.Sync() }

// parseSeq extracts sequence number from AOF filename (e.g., "000123.aof" -> 123)
func parseSeq(filename string) uint64 {
	if len(filename) < 7 || !strings.HasSuffix(filename, ".aof") {
		return 0
	}
	seqStr := filename[:len(filename)-4] // remove ".aof"
	seq, _ := strconv.ParseUint(seqStr, 10, 64)
	return seq
}
