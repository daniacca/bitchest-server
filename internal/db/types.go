package db

import "time"

type ValueType string

const (
	StringType    ValueType = "string"
	ListType      ValueType = "list"
	SortedSetType ValueType = "zset"
)

type Value interface {
	Type() ValueType
	IsExpired() bool
	Size() int
}

type StringValue struct {
	Val      string
	ExpireAt *time.Time // nil means no expiration
}

func (s *StringValue) Type() ValueType {
	return StringType
}

func (s *StringValue) Get() string {
	return s.Val
}

func (s *StringValue) IsExpired() bool {
	if s.ExpireAt == nil {
		return false
	}
	return time.Now().After(*s.ExpireAt)
}

func (s *StringValue) Size() int {
	return len(s.Val) + 8 // 8 bytes for the expiration time
}

// TODO: implement
type ListValue struct {
	Items []string
}

func (l *ListValue) Type() ValueType {
	return ListType
}

func (l *ListValue) IsExpired() bool {
	return false // TODO: implement expiration for lists
}

func (l *ListValue) Size() int {
	return 0 // TODO: implement size for lists
}

type SortedSetValue struct {
	// TODO: implement
}

func (z *SortedSetValue) Type() ValueType {
	return SortedSetType
}

func (z *SortedSetValue) IsExpired() bool {
	return false // TODO: implement expiration for sorted sets
}

func (z *SortedSetValue) Size() int {
	return 0 // TODO: implement size for sorted sets
}
