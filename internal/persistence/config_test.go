package persistence

import "testing"

func TestPersistenceConfig(t *testing.T) {
	t.Run("StorageKind should convert to string", func(t *testing.T) {
		storage := StorageFS
		if storage.String() != "fs" {
			t.Errorf("Expected 'fs', got '%s'", storage.String())
		}

		storage = StorageS3
		if storage.String() != "s3" {
			t.Errorf("Expected 's3', got '%s'", storage.String())
		}

		storage = StorageMinIO
		if storage.String() != "minio" {
			t.Errorf("Expected 'minio', got '%s'", storage.String())
		}
	})

	t.Run("StorageKind should parse valid strings", func(t *testing.T) {
		var storage StorageKind
		err := storage.Parse("fs")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if storage != StorageFS {
			t.Errorf("Expected 'fs', got '%s'", storage.String())
		}

		err = storage.Parse("s3")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if storage != StorageS3 {
			t.Errorf("Expected 's3', got '%s'", storage.String())
		}

		err = storage.Parse("minio")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if storage != StorageMinIO {
			t.Errorf("Expected 'minio', got '%s'", storage.String())
		}
	})

	t.Run("StorageKind should return error for invalid strings", func(t *testing.T) {
		var storage StorageKind
		err := storage.Parse("invalid")
		if err == nil {
			t.Error("Expected error for invalid storage kind, got nil")
		}
		if err.Error() != "cannot parse Storage string" {
			t.Errorf("Expected 'invalid storage kind: invalid', got '%s'", err.Error())
		}
	})
}