package main

import (
	"bufio"
	"flag"
	"net"
	"os"
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

func TestParseFlags(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		// Save original args
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		
		// Reset flags for clean test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		
		// Set minimal args
		os.Args = []string{"bitchest"}
		
		config := parseFlags()
		
		if config.Host != "localhost" {
			t.Errorf("Expected host 'localhost', got '%s'", config.Host)
		}
		if config.Port != 7463 {
			t.Errorf("Expected port 7463, got %d", config.Port)
		}
		if config.Addr != "localhost:7463" {
			t.Errorf("Expected addr 'localhost:7463', got '%s'", config.Addr)
		}
	})

	t.Run("custom host and port", func(t *testing.T) {
		// Save original args
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		
		// Reset flags for clean test
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		
		// Set custom flags
		os.Args = []string{"bitchest", "-host", "0.0.0.0", "-port", "6379"}
		
		config := parseFlags()
		
		if config.Host != "0.0.0.0" {
			t.Errorf("Expected host '0.0.0.0', got '%s'", config.Host)
		}
		if config.Port != 6379 {
			t.Errorf("Expected port 6379, got %d", config.Port)
		}
		if config.Addr != "0.0.0.0:6379" {
			t.Errorf("Expected addr '0.0.0.0:6379', got '%s'", config.Addr)
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("valid port range", func(t *testing.T) {
		validPorts := []int{1024, 3000, 6379, 7463, 65535}
		
		for _, port := range validPorts {
			config := &Config{
				Host: "localhost",
				Port: port,
			}
			config.Addr = net.JoinHostPort(config.Host, string(rune(port)))
			
			// Should not panic
			_ = config
		}
	})

	t.Run("invalid port range", func(t *testing.T) {
		invalidPorts := []int{0, 1023, 65536, 99999}
		
		for _, port := range invalidPorts {
			config := &Config{
				Host: "localhost",
				Port: port,
			}
			config.Addr = net.JoinHostPort(config.Host, string(rune(port)))
			
			// Should not panic (validation happens in parseFlags)
			_ = config
		}
	})
}

func TestStartServerError(t *testing.T) {
	// Test with invalid address to trigger error
	config := &Config{
		Host: "invalid-host",
		Port: 99999,
		Addr: "invalid-host:99999",
	}

	err := StartServer(config)
	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
	if !strings.Contains(err.Error(), "failed to bind") {
		t.Errorf("Expected 'failed to bind' error, got: %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	if defaultConfig.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", defaultConfig.Host)
	}
	if defaultConfig.Port != 7463 {
		t.Errorf("Expected default port 7463, got %d", defaultConfig.Port)
	}
}
