package protocol

import "testing"

func TestSimple(t *testing.T) {
	res := Simple("OK")
	if res != "+OK\r\n" {
		t.Errorf("Expected '+OK\\r\\n', got %q", res)
	}
}

func TestBulk(t *testing.T) {
	res := Bulk("hello")
	expected := "$5\r\nhello\r\n"
	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}

func TestNullBulk(t *testing.T) {
	if NullBulk() != "$-1\r\n" {
		t.Error("Expected $-1\\r\\n")
	}
}

func TestInteger(t *testing.T) {
	if Integer(42) != ":42\r\n" {
		t.Error("Expected :42\\r\\n")
	}
}

func TestArray(t *testing.T) {
	res := Array([]string{"a", "b"})
	expected := "*2\r\n$1\r\na\r\n$1\r\nb\r\n"
	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}

func TestError(t *testing.T) {
	res := Error("oops")
	if res != "-ERR oops\r\n" {
		t.Errorf("Expected '-ERR oops\\r\\n', got %q", res)
	}
}
