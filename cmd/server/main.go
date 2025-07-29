package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/handler"
)

// Config holds server configuration
type Config struct {
	Host string
	Port int
	Addr string // Computed from Host:Port
}

// Default configuration
var defaultConfig = Config{
	Host: "localhost",
	Port: 7463,
}

// parseFlags parses command line flags and returns configuration
func parseFlags() *Config {
	config := defaultConfig

	// Define command line flags
	host := flag.String("host", config.Host, "Host to bind the server to")
	port := flag.Int("port", config.Port, "Port to bind the server to")
	
	// Add help text
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Bitchest - A lightweight in-memory key-value database\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Start on localhost:7463\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -port 6379         # Start on localhost:6379\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host 0.0.0.0      # Start on all interfaces\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host 0.0.0.0 -port 6379  # Start on all interfaces:6379\n", os.Args[0])
	}

	flag.Parse()

	// Update config with flag values
	config.Host = *host
	config.Port = *port

	// Validate port range
	if config.Port < 1024 || config.Port > 65535 {
		log.Fatalf("Invalid port number: %d. Port must be between 1024 and 65535", config.Port)
	}

	// Build address string
	config.Addr = net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	return &config
}

func StartServer(config *Config) error {
	store := db.NewDB()

	ln, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind on %s: %w", config.Addr, err)
	}
	defer ln.Close()

	fmt.Printf("Bitchest is running on %s\n", config.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept: %v", err)
			continue
		}
		go handler.Handle(conn, store)
	}
}

func main() {
	config := parseFlags()
	
	if err := StartServer(config); err != nil {
		log.Fatal(err)
	}
}
