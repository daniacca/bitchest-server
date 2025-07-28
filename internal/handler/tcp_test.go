package handler

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestHandle_SET_GET(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	store := db.NewDB()

	go Handle(serverConn, store)

	writer := bufio.NewWriter(clientConn)
	reader := bufio.NewReader(clientConn)

	writer.WriteString("SET testkey testvalue\n")
	writer.Flush()
	resp1, _ := reader.ReadString('\n')
	if !strings.HasPrefix(resp1, "+OK") {
		t.Errorf("Expected +OK, got %q", resp1)
	}

	writer.WriteString("GET testkey\n")
	writer.Flush()

	line1, _ := reader.ReadString('\n')
	line2, _ := reader.ReadString('\n')

	if !strings.HasPrefix(line1, "$9") || strings.TrimSpace(line2) != "testvalue" {
		t.Errorf("Unexpected GET response: %q %q", line1, line2)
	}
}

func TestHandle_UnknownCommand(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go Handle(serverConn, db.NewDB())

	writer := bufio.NewWriter(clientConn)
	reader := bufio.NewReader(clientConn)

	writer.WriteString("FOO something\n")
	writer.Flush()

	line, _ := reader.ReadString('\n')
	if !strings.HasPrefix(line, "-ERR unknown command") {
		t.Errorf("Expected unknown command error, got %q", line)
	}
}
