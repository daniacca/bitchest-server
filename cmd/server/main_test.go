package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	// Use dynamic port (":0" -> any available)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen on dynamic port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	// Parse the address to get host and port
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("Failed to split address: %v", err)
	}

	port, err := net.LookupPort("tcp", portStr)
	if err != nil {
		t.Fatalf("Failed to lookup port: %v", err)
	}

	config := &Config{
		Host: host,
		Port: port,
		Addr: addr,
	}

	go func() {
		_ = StartServer(config) // Ignore error (exits only in test)
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	writer.WriteString("SET ping pong\n")
	writer.Flush()
	resp, _ := reader.ReadString('\n')
	if !strings.HasPrefix(resp, "+OK") {
		t.Errorf("Expected +OK, got %q", resp)
	}
}
