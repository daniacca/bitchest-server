// persistence/config.go
package persistence

import "time"

type BackendKind string

const (
	BackendFS    BackendKind = "fs"
	BackendS3    BackendKind = "s3"   	// To DO
	BackendMinIO BackendKind = "minio" 	// To DO
)

type Config struct {
	Enabled           bool
	Backend           BackendKind
	FS                FSConfig     
	S3                S3Config     // To DO
	AppendFsyncPolicy string       // "always" | "everysec" | "no"
	AOFMaxSegmentMB   int
	SnapshotInterval  time.Duration
}

type FSConfig struct {
	DataDir string // root data dir
}

type S3Config struct {
	Bucket        string
	Prefix        string // es. "bitchest/"
	Region        string
	Endpoint      string // MinIO compatibile S3
	AccessKey     string
	SecretKey     string
	ForcePathStyle bool
	UseTLS         bool
}
