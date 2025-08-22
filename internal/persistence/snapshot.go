package persistence

import (
	"bufio"
	"encoding/binary"
	"io"
	"time"
)

// dumpSnapshot writes the entire database state to the writer
func dumpSnapshot(writer SyncWriteCloser, db Store) error {
	defer writer.Close()

	// Write snapshot header
	header := []byte("BITCHEST_SNAPSHOT_V1")
	if _, err := writer.Write(header); err != nil {
		return err
	}

	// Write timestamp
	timestamp := time.Now().Unix()
	if err := binary.Write(writer, binary.BigEndian, timestamp); err != nil {
		return err
	}

	// Iterate through all string entries and write them
	var strCount uint64
	db.RangeStrings(func(key string, value string, expiresAt *time.Time) bool {
		// type marker 'S'
		if _, err := writer.Write([]byte{'S'}); err != nil { return false }
		// key
		if err := binary.Write(writer, binary.BigEndian, uint32(len(key))); err != nil { return false }
		if _, err := writer.Write([]byte(key)); err != nil { return false }
		// value
		if err := binary.Write(writer, binary.BigEndian, uint32(len(value))); err != nil { return false }
		if _, err := writer.Write([]byte(value)); err != nil { return false }
		// ttl
		hasTTL := expiresAt != nil
		if err := binary.Write(writer, binary.BigEndian, hasTTL); err != nil { return false }
		if hasTTL {
			expires := expiresAt.Unix()
			if err := binary.Write(writer, binary.BigEndian, expires); err != nil { return false }
		}
		strCount++
		return true
	})
	// Write end-of-strings marker 'E'
	if _, err := writer.Write([]byte{'E'}); err != nil { return err }

	// Iterate through all list entries and write them
	var listCount uint64
	db.RangeLists(func(key string, items []string, expiresAt *time.Time) bool {
		// type marker 'L'
		if _, err := writer.Write([]byte{'L'}); err != nil { return false }
		// key
		if err := binary.Write(writer, binary.BigEndian, uint32(len(key))); err != nil { return false }
		if _, err := writer.Write([]byte(key)); err != nil { return false }
		// items length
		if err := binary.Write(writer, binary.BigEndian, uint32(len(items))); err != nil { return false }
		for _, it := range items {
			if err := binary.Write(writer, binary.BigEndian, uint32(len(it))); err != nil { return false }
			if _, err := writer.Write([]byte(it)); err != nil { return false }
		}
		// ttl
		hasTTL := expiresAt != nil
		if err := binary.Write(writer, binary.BigEndian, hasTTL); err != nil { return false }
		if hasTTL {
			expires := expiresAt.Unix()
			if err := binary.Write(writer, binary.BigEndian, expires); err != nil { return false }
		}
		listCount++
		return true
	})
	// Write end-of-lists marker 'e' (lowercase to differ from strings end if needed)
	if _, err := writer.Write([]byte{'e'}); err != nil { return err }

	// Write counts at end for validation
	if err := binary.Write(writer, binary.BigEndian, strCount); err != nil { return err }
	if err := binary.Write(writer, binary.BigEndian, listCount); err != nil { return err }

	return writer.Sync()
}

// restoreSnapshot reads and restores the database state from the reader
func restoreSnapshot(reader io.ReadCloser, apply func([]byte) error) error {
	defer reader.Close()

	// Read and validate header
	header := make([]byte, 20) // "BITCHEST_SNAPSHOT_V1"
	if _, err := io.ReadFull(reader, header); err != nil {
		return err
	}
	if string(header) != "BITCHEST_SNAPSHOT_V1" {
		return ErrInvalidSnapshotFormat
	}

	// Read timestamp
	var timestamp int64
	if err := binary.Read(reader, binary.BigEndian, &timestamp); err != nil {
		return err
	}

	// Read entries by markers
	for {
		marker := make([]byte, 1)
		if _, err := io.ReadFull(reader, marker); err != nil {
			if err == io.EOF { break }
			return err
		}
		switch marker[0] {
		case 'S': // string entry
			var keyLen uint32
			if err := binary.Read(reader, binary.BigEndian, &keyLen); err != nil { return err }
			keyBytes := make([]byte, keyLen)
			if _, err := io.ReadFull(reader, keyBytes); err != nil { return err }
			var valueLen uint32
			if err := binary.Read(reader, binary.BigEndian, &valueLen); err != nil { return err }
			value := make([]byte, valueLen)
			if _, err := io.ReadFull(reader, value); err != nil { return err }
			var hasTTL bool
			if err := binary.Read(reader, binary.BigEndian, &hasTTL); err != nil { return err }
			if hasTTL {
				var expiresUnix int64
				if err := binary.Read(reader, binary.BigEndian, &expiresUnix); err != nil { return err }
				// Optionally, could use EX here
			}
			if err := apply([]byte("SET "+string(keyBytes)+" "+string(value))); err != nil { return err }
		case 'E': // end of strings section
			// proceed to lists section
		case 'L': // list entry
			var keyLen uint32
			if err := binary.Read(reader, binary.BigEndian, &keyLen); err != nil { return err }
			keyBytes := make([]byte, keyLen)
			if _, err := io.ReadFull(reader, keyBytes); err != nil { return err }
			var numItems uint32
			if err := binary.Read(reader, binary.BigEndian, &numItems); err != nil { return err }
			// ensure the list is created empty
			if err := apply([]byte("DEL "+string(keyBytes))); err != nil { return err }
			for i := uint32(0); i < numItems; i++ {
				var ilen uint32
				if err := binary.Read(reader, binary.BigEndian, &ilen); err != nil { return err }
				item := make([]byte, ilen)
				if _, err := io.ReadFull(reader, item); err != nil { return err }
				if err := apply([]byte("RPUSH "+string(keyBytes)+" "+string(item))); err != nil { return err }
			}
			var hasTTL bool
			if err := binary.Read(reader, binary.BigEndian, &hasTTL); err != nil { return err }
			if hasTTL {
				var expiresUnix int64
				if err := binary.Read(reader, binary.BigEndian, &expiresUnix); err != nil { return err }
				// Could compute EX seconds relative, but snapshot uses absolute; skipped for now
			}
		case 'e': // end of lists section
			// done; there may be trailing counts, read and ignore
			var strCount, listCount uint64
			_ = binary.Read(reader, binary.BigEndian, &strCount)
			_ = binary.Read(reader, binary.BigEndian, &listCount)
			return nil
		default:
			return ErrInvalidSnapshotFormat
		}
	}
	return nil
}

// replayAOF reads and replays AOF commands from the reader
func replayAOF(reader io.ReadCloser, apply func([]byte) error) error {
	defer reader.Close()

	bufReader := bufio.NewReader(reader)

	for {
		// Read command length (4 bytes)
		lengthBytes := make([]byte, 4)
		if _, err := io.ReadFull(bufReader, lengthBytes); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		length := binary.BigEndian.Uint32(lengthBytes)

		// Read command
		cmd := make([]byte, length)
		if _, err := io.ReadFull(bufReader, cmd); err != nil {
			return err
		}

		// Apply command using provided executor
		if err := apply(cmd); err != nil {
			return err
		}
	}

	return nil
}

// ErrInvalidSnapshotFormat is returned when snapshot format is invalid
var ErrInvalidSnapshotFormat = &InvalidFormatError{}

type InvalidFormatError struct{}

func (e *InvalidFormatError) Error() string {
	return "invalid snapshot format"
} 