// persistence/config.go
package persistence

import (
	"fmt"
	"time"
)

type StorageKind string

const (
	StorageFS    StorageKind = "fs"
	StorageS3    StorageKind = "s3"   	// To DO
	StorageMinIO StorageKind = "minio" 	// To DO
)

func (k StorageKind) String() string {
	return string(k);
}

func (k *StorageKind) Parse(s string) (error) {
	validationMap := map[StorageKind]struct{}{
		StorageFS: {},
		StorageS3: {},
		StorageMinIO: {},
	}

	parsed := StorageKind(s)
	if _, ok := validationMap[parsed]; !ok {
		return fmt.Errorf("cannot parse Storage string");
	}

	*k = parsed
	return nil
}

type Config struct {
	Enabled           bool
	Backend           StorageKind
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
