/*
* Protocol
*
* This package implements the protocol for the Bitchest server.
*
* The protocol is a simple text-based protocol that is used to communicate with the server.
* The protocol is based on the RESP protocol, from Redis.
 */
package protocol

import (
	"fmt"
)

// Simple string response
// Encapsulates a simple string response with a + prefix and a \r\n suffix
// Used for simple and short responses, like PONG, OK, etc.
func Simple(msg string) string {
	return "+" + msg + "\r\n"
}

// Bulk string response
// Return two lines:
// 1. $<length>
// 2. <message>
func Bulk(msg string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
}

// Integer response
// Encapsulates an integer response with a : prefix and a \r\n suffix
// Used for integer responses, like the number of items in a list, etc.
func Integer(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}

// Error response
// Encapsulates an error response with a - prefix and a \r\n suffix
// Used for error responses, like ERR, followed by the error message.
func Error(msg string) string {
	return "-ERR " + msg + "\r\n"
}

// Array response
// Return a new line for each element in the list of items, each item is a bulk string.
// The first line is the number of items in the array.
func Array(items []string) string {
	out := fmt.Sprintf("*%d\r\n", len(items))
	for _, i := range items {
		out += Bulk(i)
	}
	return out
}

// Null bulk string response
// Used for commands that could result in a null value, like GET, etc.
func NullBulk() string {
	return "$-1\r\n"
}