package db

type ValueType string

const (
	StringType    ValueType = "string"
	ListType      ValueType = "list"
	SortedSetType ValueType = "zset"
)

type Value interface {
	Type() ValueType
}

type StringValue struct {
	Val string
}

func (s *StringValue) Type() ValueType {
	return StringType
}

func (s *StringValue) Get() string {
	return s.Val
}

// TODO: implement
type ListValue struct {
	Items []string
}

func (l *ListValue) Type() ValueType {
	return ListType
}

type SortedSetValue struct {
	// TODO: implement
}

func (z *SortedSetValue) Type() ValueType {
	return SortedSetType
}