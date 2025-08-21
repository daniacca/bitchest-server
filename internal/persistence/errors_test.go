package persistence

import "testing"

func TestError(t *testing.T) {
	t.Run("ConfigError should return correct message", func(t *testing.T) {
		err := ErrInvalidConfig
		if err.Error() != "invalid configuration" {
			t.Errorf("Expected 'invalid configuration', got '%s'", err.Error())
		}
	})

	t.Run("ConfigError should implement error interface", func(t *testing.T) {
		var _ error = &ConfigError{}
	})
}