package persistence

// Common error types for the persistence layer
var (
	ErrInvalidConfig = &ConfigError{msg: "invalid configuration"}
)

type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string {
	return e.msg
} 