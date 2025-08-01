package db

import (
	"testing"
	"time"
)

func TestStringValue_Size(t *testing.T) {
	value := &StringValue{
		Val: "test", // 4 bytes
	}

	if value.Size() != 12 { // 4 bytes for the value + 8 bytes for the expiration time
		t.Errorf("Expected size to be 12, got %d", value.Size())
	}
}

func TestStringValue_IsExpired(t *testing.T) {
	value := &StringValue{
		Val: "test",
	}

	// No expiration time, so not expired
	if value.IsExpired() {
		t.Errorf("Expected value to not be expired, got %t", value.IsExpired())
	}

	expireAt := time.Now().Add(time.Second * -2) // 2 seconds in the past
	value.ExpireAt = &expireAt

	if !value.IsExpired() {
		t.Errorf("Expected value to be expired, got %t", value.IsExpired())
	}
}

func TestStringValue_Type(t *testing.T) {
	value := &StringValue{
		Val: "test",
	}

	if value.Type() != StringType {
		t.Errorf("Expected type to be %s, got %s", StringType, value.Type())
	}
}

func TestListValue_Type(t *testing.T) {
	value := &ListValue{
		Items: []string{"test"},
	}

	if value.Type() != ListType {
		t.Errorf("Expected type to be %s, got %s", ListType, value.Type())
	}
}

func TestListValue_IsExpired(t *testing.T) {
	value := &ListValue{
		Items: []string{"test"},
	}

	if value.IsExpired() {
		t.Errorf("Expected value to not be expired, got %t", value.IsExpired())
	}
}

func TestSortedSetValue_Type(t *testing.T) {
	value := &SortedSetValue{}

	if value.Type() != SortedSetType {
		t.Errorf("Expected type to be %s, got %s", SortedSetType, value.Type())
	}
}