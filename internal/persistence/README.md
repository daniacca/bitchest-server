# Persistence Layer

This package provides a robust persistence layer for the Bitchest in-memory key-value store, supporting both snapshots and append-only file (AOF) logging.

## Architecture

The persistence layer is designed with extensibility in mind, supporting multiple storage backends:

- **Filesystem** (implemented) - Local file storage
- **S3** (planned) - Amazon S3 storage
- **MinIO** (planned) - MinIO/S3-compatible storage

## Components

### Core Types

- `Manager` - Main persistence manager that orchestrates snapshots and AOF
- `StorageAdapter` - Interface for different storage backends
- `Store` - Interface that the database must implement
- `Entry` - Represents a key-value entry with optional TTL

### Storage Adapters

- `fsadapter.FSAdapter` - Filesystem-based storage implementation

## Features

### Snapshots

- **Binary format** - Efficient binary serialization of database state
- **Atomic creation** - Snapshots are created atomically to ensure consistency
- **Background creation** - Snapshots are created in the background without blocking writes
- **Automatic rotation** - Configurable snapshot intervals

### Append-Only File (AOF)

- **RESP format** - Commands stored in Redis RESP format
- **Segmented** - AOF is split into segments for efficient management
- **Configurable fsync** - Support for different fsync policies:
  - `always` - Sync after every write
  - `everysec` - Sync every second
  - `no` - Let OS handle syncing
- **Automatic rotation** - Segments are rotated when they reach a configurable size

### Recovery

- **Snapshot + AOF** - Recovery combines latest snapshot with subsequent AOF segments
- **Best-effort loading** - Snapshot loading failures don't prevent startup
- **Sequential replay** - AOF segments are replayed in sequence order

## Configuration

```go
cfg := persistence.Config{
    Enabled:           true,
    Backend:           persistence.BackendFS,
    FS:                persistence.FSConfig{DataDir: "./data"},
    AppendFsyncPolicy: "everysec",
    AOFMaxSegmentMB:   10,
    SnapshotInterval:  time.Minute * 5,
}
```

## Usage

```go
// Create filesystem adapter
fsAdapter := fsadapter.New("./data")

// Create persistence manager
manager := persistence.NewManager(cfg, fsAdapter)

// Start the manager
ctx := context.Background()
if err := manager.Start(ctx, db); err != nil {
    log.Fatal(err)
}

// Record mutations
manager.OnMutationRESP([]byte("SET key value"))

// Stop the manager
manager.Stop(ctx)
```

## File Structure

```
data/
├── manifest.json          # Persistence metadata
├── aof/                   # AOF segments
│   ├── 000001.aof
│   ├── 000002.aof
│   └── ...
└── snapshots/             # Database snapshots
    ├── snapshot-1234567890.bin
    └── snapshot-1234567950.bin
```

## Binary Format

### Snapshot Format

```
Header: "BITCHEST_SNAPSHOT_V1" (19 bytes)
Timestamp: int64 (8 bytes, big-endian)
Entries: [Entry...]
Entry Count: uint64 (8 bytes, big-endian)

Entry Format:
- Key Length: uint32 (4 bytes, big-endian)
- Key: []byte
- Value Length: uint32 (4 bytes, big-endian)
- Value: []byte
- Has TTL: bool (1 byte)
- TTL (if Has TTL): int64 (8 bytes, big-endian)
```

### AOF Format

```
Command Length: uint32 (4 bytes, big-endian)
Command: []byte (RESP format)
```

## Extending

To add a new storage backend:

1. Implement the `StorageAdapter` interface
2. Add the backend type to `BackendKind` enum
3. Add configuration struct for the backend
4. Update the manager to handle the new backend

## Error Handling

- **Snapshot failures** are logged but don't prevent normal operation
- **AOF write failures** are propagated to the caller
- **Recovery failures** are logged but the system continues with an empty state
- **Configuration errors** prevent startup

## Performance Considerations

- **Memory usage** - Snapshots create a copy of the data structure
- **I/O performance** - AOF fsync policy affects write performance
- **Recovery time** - Larger datasets take longer to recover
- **Storage space** - Both snapshots and AOF consume disk space
