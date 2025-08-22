package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"time"
)

// Mock SyncWriteCloser
type mockSyncWriteCloser struct {
	buf    *bytes.Buffer
	closed bool
	synced bool
}

func (m *mockSyncWriteCloser) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}
func (m *mockSyncWriteCloser) Close() error {
	m.closed = true
	return nil
}
func (m *mockSyncWriteCloser) Sync() error {
	m.synced = true
	return nil
}

// Mock Store
type mockStore struct {
	strings []struct {
		key       string
		value     string
		expiresAt *time.Time
	}
	lists []struct {
		key       string
		items     []string
		expiresAt *time.Time
	}
}

func (m *mockStore) RangeStrings(fn func(key, value string, expiresAt *time.Time) bool) {
	for _, s := range m.strings {
		if !fn(s.key, s.value, s.expiresAt) {
			break
		}
	}
}

func (m *mockStore) RangeLists(fn func(key string, items []string, expiresAt *time.Time) bool) {
	for _, l := range m.lists {
		if !fn(l.key, l.items, l.expiresAt) {
			break
		}
	}
}

func TestDumpSnapshot_EmptyDB(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := &mockSyncWriteCloser{buf: buf}
	db := &mockStore{}

	err := dumpSnapshot(writer, db)
	if err != nil {
		t.Fatalf("dumpSnapshot failed: %v", err)
	}
	if !writer.closed {
		t.Error("writer not closed")
	}
	if !writer.synced {
		t.Error("writer not synced")
	}

	// Validate header
	got := buf.Bytes()
	if len(got) < 20 {
		t.Fatal("snapshot too short")
	}
	if string(got[:20]) != "BITCHEST_SNAPSHOT_V1" {
		t.Errorf("header mismatch: %q", got[:20])
	}

	// Validate markers
	// After header (20) + timestamp (8), should be 'E', 'e', then two uint64 counts
	if got[28] != 'E' {
		t.Errorf("expected 'E' marker, got %q", got[28])
	}
	if got[29] != 'e' {
		t.Errorf("expected 'e' marker, got %q", got[29])
	}

	// Counts should be zero
	strCount := binary.BigEndian.Uint64(got[30:38])
	listCount := binary.BigEndian.Uint64(got[38:46])
	if strCount != 0 {
		t.Errorf("expected strCount 0, got %d", strCount)
	}
	if listCount != 0 {
		t.Errorf("expected listCount 0, got %d", listCount)
	}
}

func TestDumpSnapshot_WithStringsAndLists(t *testing.T) {
	now := time.Now()
	exp := now.Add(10 * time.Second)
	db := &mockStore{
		strings: []struct {
			key       string
			value     string
			expiresAt *time.Time
		}{
			{"foo", "bar", nil},
			{"baz", "qux", &exp},
		},
		lists: []struct {
			key       string
			items     []string
			expiresAt *time.Time
		}{
			{"mylist", []string{"a", "b"}, nil},
			{"expirelist", []string{"x"}, &exp},
		},
	}
	buf := &bytes.Buffer{}
	writer := &mockSyncWriteCloser{buf: buf}

	err := dumpSnapshot(writer, db)
	if err != nil {
		t.Fatalf("dumpSnapshot failed: %v", err)
	}
	got := buf.Bytes()

	// Validate header
	if string(got[:20]) != "BITCHEST_SNAPSHOT_V1" {
		t.Errorf("header mismatch: %q", got[:20])
	}

	// Check String section
	offset := 20 + 8 // header (20 chars) + timestamp (8 bytes)
	for i := range 2 {
		if got[offset] != 'S' {
			t.Errorf("expected 'S' marker at %d, got %q", offset, got[offset])
		}
		offset++
		keyLen := binary.BigEndian.Uint32(got[offset : offset+4])
		offset += 4
		key := string(got[offset : offset+int(keyLen)])
		offset += int(keyLen)
		valLen := binary.BigEndian.Uint32(got[offset : offset+4])
		offset += 4
		val := string(got[offset : offset+int(valLen)])
		offset += int(valLen)
		hasTTL := got[offset]
		offset++
		if hasTTL != 0 {
			_ = binary.BigEndian.Uint64(got[offset : offset+8])
			offset += 8
		}
		// Check key/value
		if i == 0 && (key != "foo" || val != "bar") {
			t.Errorf("unexpected string entry: %q %q", key, val)
		}
		if i == 1 && (key != "baz" || val != "qux") {
			t.Errorf("unexpected string entry: %q %q", key, val)
		}
	}
	if got[offset] != 'E' {
		t.Errorf("expected 'E' marker after strings, got %q", got[offset])
	}
	offset++

	// List section
	for i := range 2 {
		if got[offset] != 'L' {
			t.Errorf("expected 'L' marker at %d, got %q", offset, got[offset])
		}
		offset++
		keyLen := binary.BigEndian.Uint32(got[offset : offset+4])
		offset += 4
		key := string(got[offset : offset+int(keyLen)])
		offset += int(keyLen)
		numItems := binary.BigEndian.Uint32(got[offset : offset+4])
		offset += 4
		items := []string{}
		for j := uint32(0); j < numItems; j++ {
			ilen := binary.BigEndian.Uint32(got[offset : offset+4])
			offset += 4
			item := string(got[offset : offset+int(ilen)])
			offset += int(ilen)
			items = append(items, item)
		}
		hasTTL := got[offset]
		offset++
		if hasTTL != 0 {
			_ = binary.BigEndian.Uint64(got[offset : offset+8])
			offset += 8
		}
		// Check key/items
		if i == 0 && (key != "mylist" || len(items) != 2 || items[0] != "a" || items[1] != "b") {
			t.Errorf("unexpected list entry: %q %v", key, items)
		}
		if i == 1 && (key != "expirelist" || len(items) != 1 || items[0] != "x") {
			t.Errorf("unexpected list entry: %q %v", key, items)
		}
	}
	if got[offset] != 'e' {
		t.Errorf("expected 'e' marker after lists, got %q", got[offset])
	}
	offset++

	// Counts
	strCount := binary.BigEndian.Uint64(got[offset : offset+8])
	listCount := binary.BigEndian.Uint64(got[offset+8 : offset+16])
	if strCount != 2 {
		t.Errorf("expected strCount 2, got %d", strCount)
	}
	if listCount != 2 {
		t.Errorf("expected listCount 2, got %d", listCount)
	}
}

// Writer that fails after writing header
type mockFailWriter struct {
	buf    *bytes.Buffer
}

func (m *mockFailWriter) Write(p []byte) (int, error) {
	if bytes.Equal(p, []byte("BITCHEST_SNAPSHOT_V1")) {
		return m.buf.Write(p)
	}
	return 0, io.ErrClosedPipe
}

func (m *mockFailWriter) Close() error { return nil }

func (m *mockFailWriter) Sync() error {	return nil }

func TestDumpSnapshot_ErrorPropagation(t *testing.T) {
	// Writer that always fails
	writer := &mockFailWriter{
		buf: &bytes.Buffer{},
	}
	db := &mockStore{
		strings: []struct {
			key       string
			value     string
			expiresAt *time.Time
		}{
			{"fail", "fail", nil},
		},
	}
	err := dumpSnapshot(writer, db)
	if err == nil {
		t.Error("expected error from writer, got nil")
	}
}

func TestRestoreSnapshot_EmptyDB(t *testing.T) {
	// Prepare a snapshot with no strings/lists
	buf := &bytes.Buffer{}
	buf.Write([]byte("BITCHEST_SNAPSHOT_V1")) // header (20 bytes)
	binary.Write(buf, binary.BigEndian, time.Now().Unix()) // timestamp (8 bytes)
	buf.WriteByte('E') // end of strings
	buf.WriteByte('e') // end of lists
	binary.Write(buf, binary.BigEndian, uint64(0)) // strCount
	binary.Write(buf, binary.BigEndian, uint64(0)) // listCount

	var applied [][]byte
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := restoreSnapshot(reader, func(cmd []byte) error {
		applied = append(applied, cmd)
		return nil
	})
	if err != nil {
		t.Fatalf("restoreSnapshot failed: %v", err)
	}
	if len(applied) != 0 {
		t.Errorf("expected no commands applied, got %d", len(applied))
	}
}

func TestRestoreSnapshot_WithStringsAndLists(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("BITCHEST_SNAPSHOT_V1"))
	binary.Write(buf, binary.BigEndian, int64(1234567890)) // timestamp

	// String entry: foo -> bar (no TTL)
	buf.WriteByte('S')
	binary.Write(buf, binary.BigEndian, uint32(3)) // keyLen
	buf.Write([]byte("foo"))
	binary.Write(buf, binary.BigEndian, uint32(3)) // valLen
	buf.Write([]byte("bar"))
	binary.Write(buf, binary.BigEndian, false) // hasTTL

	// String entry: baz -> qux (with TTL)
	buf.WriteByte('S')
	binary.Write(buf, binary.BigEndian, uint32(3))
	buf.Write([]byte("baz"))
	binary.Write(buf, binary.BigEndian, uint32(3))
	buf.Write([]byte("qux"))
	binary.Write(buf, binary.BigEndian, true)
	binary.Write(buf, binary.BigEndian, int64(9876543210)) // expiresUnix

	buf.WriteByte('E') // end of strings

	// List entry: mylist -> ["a", "b"] (no TTL)
	buf.WriteByte('L')
	binary.Write(buf, binary.BigEndian, uint32(6))
	buf.Write([]byte("mylist"))
	binary.Write(buf, binary.BigEndian, uint32(2)) // numItems
	binary.Write(buf, binary.BigEndian, uint32(1))
	buf.Write([]byte("a"))
	binary.Write(buf, binary.BigEndian, uint32(1))
	buf.Write([]byte("b"))
	binary.Write(buf, binary.BigEndian, false)

	// List entry: expirelist -> ["x"] (with TTL)
	buf.WriteByte('L')
	binary.Write(buf, binary.BigEndian, uint32(10))
	buf.Write([]byte("expirelist"))
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, uint32(1))
	buf.Write([]byte("x"))
	binary.Write(buf, binary.BigEndian, true)
	binary.Write(buf, binary.BigEndian, int64(5555555555))

	buf.WriteByte('e') // end of lists
	binary.Write(buf, binary.BigEndian, uint64(2)) // strCount
	binary.Write(buf, binary.BigEndian, uint64(2)) // listCount

	var applied [][]byte
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := restoreSnapshot(reader, func(cmd []byte) error {
		applied = append(applied, cmd)
		return nil
	})
	if err != nil {
		t.Fatalf("restoreSnapshot failed: %v", err)
	}

	want := [][]byte{
		[]byte("SET foo bar"),
		[]byte("SET baz qux"),
		[]byte("DEL mylist"),
		[]byte("RPUSH mylist a"),
		[]byte("RPUSH mylist b"),
		[]byte("DEL expirelist"),
		[]byte("RPUSH expirelist x"),
	}
	if len(applied) != len(want) {
		t.Fatalf("expected %d commands, got %d", len(want), len(applied))
	}
	for i := range want {
		if string(applied[i]) != string(want[i]) {
			t.Errorf("cmd %d: want %q, got %q", i, want[i], applied[i])
		}
	}
}

func TestRestoreSnapshot_InvalidHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("INVALID_SNAPSHOT_HDR"))
	binary.Write(buf, binary.BigEndian, int64(0))
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := restoreSnapshot(reader, func(cmd []byte) error { return nil })
	if err != ErrInvalidSnapshotFormat {
		t.Errorf("expected ErrInvalidSnapshotFormat, got %v", err)
	}
}

func TestRestoreSnapshot_InvalidMarker(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("BITCHEST_SNAPSHOT_V1"))
	binary.Write(buf, binary.BigEndian, int64(0))
	buf.WriteByte('X') // invalid marker
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := restoreSnapshot(reader, func(cmd []byte) error { return nil })
	if err != ErrInvalidSnapshotFormat {
		t.Errorf("expected ErrInvalidSnapshotFormat, got %v", err)
	}
}

func TestRestoreSnapshot_ApplyError(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("BITCHEST_SNAPSHOT_V1"))
	binary.Write(buf, binary.BigEndian, int64(0))
	// Add a string entry
	buf.WriteByte('S')
	binary.Write(buf, binary.BigEndian, uint32(3))
	buf.Write([]byte("foo"))
	binary.Write(buf, binary.BigEndian, uint32(3))
	buf.Write([]byte("bar"))
	binary.Write(buf, binary.BigEndian, false)
	buf.WriteByte('E')
	buf.WriteByte('e')
	binary.Write(buf, binary.BigEndian, uint64(1))
	binary.Write(buf, binary.BigEndian, uint64(0))

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := restoreSnapshot(reader, func(cmd []byte) error {
		return io.ErrClosedPipe
	})
	if err != io.ErrClosedPipe {
		t.Errorf("expected apply error, got %v", err)
	}
}

func TestReplayAOF_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	var applied [][]byte
	err := replayAOF(reader, func(cmd []byte) error {
		applied = append(applied, cmd)
		return nil
	})
	if err != nil {
		t.Fatalf("replayAOF failed: %v", err)
	}
	if len(applied) != 0 {
		t.Errorf("expected no commands applied, got %d", len(applied))
	}
}

func TestReplayAOF_SingleCommand(t *testing.T) {
	cmd := []byte("SET foo bar")
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, uint32(len(cmd)))
	buf.Write(cmd)
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	var applied [][]byte
	err := replayAOF(reader, func(cmd []byte) error {
		applied = append(applied, cmd)
		return nil
	})
	if err != nil {
		t.Fatalf("replayAOF failed: %v", err)
	}
	if len(applied) != 1 {
		t.Fatalf("expected 1 command, got %d", len(applied))
	}
	if string(applied[0]) != "SET foo bar" {
		t.Errorf("expected command %q, got %q", "SET foo bar", applied[0])
	}
}

func TestReplayAOF_MultipleCommands(t *testing.T) {
	cmds := [][]byte{
		[]byte("SET foo bar"),
		[]byte("DEL foo"),
		[]byte("RPUSH mylist a"),
	}
	buf := &bytes.Buffer{}
	for _, cmd := range cmds {
		binary.Write(buf, binary.BigEndian, uint32(len(cmd)))
		buf.Write(cmd)
	}
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	var applied [][]byte
	err := replayAOF(reader, func(cmd []byte) error {
		applied = append(applied, cmd)
		return nil
	})
	if err != nil {
		t.Fatalf("replayAOF failed: %v", err)
	}
	if len(applied) != len(cmds) {
		t.Fatalf("expected %d commands, got %d", len(cmds), len(applied))
	}
	for i := range cmds {
		if string(applied[i]) != string(cmds[i]) {
			t.Errorf("cmd %d: want %q, got %q", i, cmds[i], applied[i])
		}
	}
}

func TestReplayAOF_ApplyError(t *testing.T) {
	cmd := []byte("SET foo bar")
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, uint32(len(cmd)))
	buf.Write(cmd)
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := replayAOF(reader, func(cmd []byte) error {
		return io.ErrClosedPipe
	})
	if err != io.ErrClosedPipe {
		t.Errorf("expected apply error, got %v", err)
	}
}

func TestReplayAOF_InvalidLength(t *testing.T) {
	// Not enough bytes for length
	buf := &bytes.Buffer{}
	buf.Write([]byte{0x00, 0x00, 0x00}) // only 3 bytes, should be 4
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := replayAOF(reader, func(cmd []byte) error { return nil })
	if err == nil {
		t.Error("expected error due to invalid length, got nil")
	}
}

func TestReplayAOF_InvalidCommandData(t *testing.T) {
	// Length says 10, but only 5 bytes present
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, uint32(10))
	buf.Write([]byte("short"))
	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	err := replayAOF(reader, func(cmd []byte) error { return nil })
	if err == nil {
		t.Error("expected error due to short command data, got nil")
	}
}

func TestInvalidFormatError_Error(t *testing.T) {
	err := &InvalidFormatError{}
	want := "invalid snapshot format"
	got := err.Error()
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
