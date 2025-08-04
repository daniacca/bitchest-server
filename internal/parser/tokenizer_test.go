package parser

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		expectedErr error
	}{
		{
			name:     "simple command",
			input:    "SET key value",
			expected: []string{"SET", "key", "value"},
			expectedErr: nil,
		},
		{
			name:     "quoted string with spaces",
			input:    `SET key "value with spaces"`,
			expected: []string{"SET", "key", "value with spaces"},
			expectedErr: nil,
		},
		{
			name:     "quoted string with options",
			input:    `SET key "my string with space" EX 60`,
			expected: []string{"SET", "key", "my string with space", "EX", "60"},
			expectedErr: nil,
		},
		{
			name:     "quoted string with NX option",
			input:    `SET key "quoted value" NX`,
			expected: []string{"SET", "key", "quoted value", "NX"},
			expectedErr: nil,
		},
		{
			name:     "multiple quoted strings",
			input:    `SET key1 "value one" key2 "value two"`,
			expected: []string{"SET", "key1", "value one", "key2", "value two"},
			expectedErr: nil,
		},
		{
			name:     "escaped quotes",
			input:    `SET key "value with \"quotes\""`,
			expected: []string{"SET", "key", "value with \"quotes\""},
			expectedErr: nil,
		},
		{
			name:     "empty quoted string",
			input:    `SET key ""`,
			expected: []string{"SET", "key"},
			expectedErr: nil,
		},
		{
			name:     "command with no arguments",
			input:    "PING",
			expected: []string{"PING"},
			expectedErr: nil,
		},
		{
			name:     "command with empty input",
			input:    "",
			expected: []string{},
			expectedErr: nil,
		},
		{
			name:     "command with only spaces",
			input:    "   ",
			expected: []string{},
			expectedErr: nil,
		},
		{
			name:     "quoted string with special characters",
			input:    `SET key "value with @#$%^&*()"`,
			expected: []string{"SET", "key", "value with @#$%^&*()"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Tokenize(tt.input)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("parseCommand(%q) = %v, want %v", tt.input, err, tt.expectedErr)
			} else if err == nil && tt.expectedErr != nil {
				t.Errorf("parseCommand(%q) = %v, want error", tt.input, tt.expectedErr)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseCommand(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseCommandEdgeCases(t *testing.T) {
	t.Run("unclosed quote", func(t *testing.T) {
		input := `SET key "unclosed quote`
		result, err := Tokenize(input)
		if err == nil {
			t.Errorf("parseCommand(%q) = %v, want error", input, err)
		}
		if result != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, result)
		}
	})

	t.Run("quotes in middle of word", func(t *testing.T) {
		input := `SET key"value"more`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key\"value\"more"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("multiple spaces", func(t *testing.T) {
		input := "SET   key    value"
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("escaped quotes", func(t *testing.T) {
		input := `SET key "value with \"quotes\""`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", "value with \"quotes\""}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("unterminated quote", func(t *testing.T) {
		input := `SET key "unterminated quote`
		result, err := Tokenize(input)
		if err == nil {
			t.Errorf("parseCommand(%q) = %v, want error", input, err)
		}
		if result != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, result)
		}
	})

	t.Run("set command with json input", func(t *testing.T) {
		input := `SET key '{"key": "value"}'`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", `{"key": "value"}`}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("set command with json input and escaped quotes", func(t *testing.T) {
		input := `SET key "{\"key\": \"value\"}"`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", `{"key": "value"}`}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("set command with json input and escaped quotes inside the json", func(t *testing.T) {
		input := `SET key '{"key": "value with \"quotes\""}' EX 60`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", `{"key": "value with "quotes""}`, "EX", "60"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})

	t.Run("set command with complex json input", func(t *testing.T) {
		input := `SET key '{"key": "value with \"quotes\"", "key2": 14, "key3": {"key4": "value4"}, "key5": [1, 2, 3] }'`
		result, err := Tokenize(input)
		if err != nil {
			t.Errorf("parseCommand(%q) = %v, want nil", input, err)
		}
		expected := []string{"SET", "key", `{"key": "value with "quotes"", "key2": 14, "key3": {"key4": "value4"}, "key5": [1, 2, 3] }`}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("parseCommand(%q) = %v, want %v", input, result, expected)
		}
	})
}