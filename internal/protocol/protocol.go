package protocol

import (
	"fmt"
)

func Simple(msg string) string {
	return "+" + msg + "\r\n"
}

func Bulk(msg string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
}

func Integer(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}

func Error(msg string) string {
	return "-ERR " + msg + "\r\n"
}

func Array(items []string) string {
	out := fmt.Sprintf("*%d\r\n", len(items))
	for _, i := range items {
		out += Bulk(i)
	}
	return out
}

func NullBulk() string {
	return "$-1\r\n"
}