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

type SortedSetValue struct {
	// TODO: implement
}

func (z *SortedSetValue) Type() ValueType {
	return SortedSetType
}

func (z *SortedSetValue) IsExpired() bool {
	return false // TODO: implement expiration for sorted sets
}