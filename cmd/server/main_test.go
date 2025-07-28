package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	// Usa porta dinamica (":0" -> qualunque disponibile)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen on dynamic port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	go func() {
		_ = StartServer(addr) // Ignora l'errore (esce solo in test)
	}()

	// Attendi che il server sia partito
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
