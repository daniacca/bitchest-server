package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/daniacca/bitchest/internal/commands"
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/handler"
	"github.com/daniacca/bitchest/internal/persistence"
	"github.com/daniacca/bitchest/internal/persistence/fsadapter"
)

// Config holds server configuration
type Config struct {
	Host string
	Port int
	Addr string // Computed from Host:Port
	EnablePersistance bool
	StorageSupport persistence.StorageKind
	Directory string
	SnapshotInterval int // minutes
	AOFMaxSegmentMB int
}

// Default configuration
var defaultConfig = Config{
	Host: "localhost",
	Port: 7463,
	EnablePersistance: true,             // persistence enabled
	StorageSupport: persistence.StorageFS, // default to file system
	Directory: "./data",
	SnapshotInterval: 5, // minutes
	AOFMaxSegmentMB: 10,
}

// parseFlags parses command line flags and returns configuration
func parseFlags() *Config {
	config := defaultConfig

	// Define command line flags
	host := flag.String("host", config.Host, "Host to bind the server to")
	port := flag.Int("port", config.Port, "Port to bind the server to")
	enablePersistence := flag.Bool("enable-persistence", config.EnablePersistance, "Enable writing a copy of DB data on persistent storage")
	storageSupport := flag.String("support", string(config.StorageSupport), "Storage backend for persistence: fs | s3 | minio")
	dataDirectory := flag.String("out-dir", config.Directory, "Output directory for saving data (fs backend)")
	snapshotInterval := flag.Int("snapshot-interval", config.SnapshotInterval, "Snapshot interval, in minutes")
	aofMaxSegmentMB := flag.Int("aof-max-segment-mb", config.AOFMaxSegmentMB, "Max size (MB) per AOF segment")

	// Add help text
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Bitchest - A lightweight in-memory key-value database\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                      # Start on localhost:7463\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -port 6379                           # Start on localhost:6379\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host 0.0.0.0                        # Start on all interfaces\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host 0.0.0.0 -port 6379             # Start on all interfaces:6379\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -enable-persistence=false            # Disable persistence\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -support fs -out-dir ./data          # Use filesystem backend with data dir\n", os.Args[0])
	}

	flag.Parse()

	// Update config with flag values
	config.Host = *host
	config.Port = *port
	config.EnablePersistance = *enablePersistence
	if err := config.StorageSupport.Parse(*storageSupport); err != nil {
		log.Printf("Invalid storage support specified: %q. Valid values: 'fs' | 's3' | 'minio'. Using default '%s'.", *storageSupport, defaultConfig.StorageSupport)
		config.StorageSupport = defaultConfig.StorageSupport
	}
	config.Directory = *dataDirectory
	config.SnapshotInterval = *snapshotInterval
	config.AOFMaxSegmentMB = *aofMaxSegmentMB

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

	cfg := persistence.Config{
		Enabled:           config.EnablePersistance,
		Backend:           config.StorageSupport,
		FS:                persistence.FSConfig{ DataDir: config.Directory },
		AppendFsyncPolicy: "everysec",
		AOFMaxSegmentMB:   config.AOFMaxSegmentMB,
		SnapshotInterval:  time.Minute * time.Duration(config.SnapshotInterval),
	}

	// Create storage adapter
	var storageAdapter persistence.StorageAdapter
	switch cfg.Backend {
	case persistence.StorageFS:
		storageAdapter = fsadapter.New(cfg.FS.DataDir)
	default:
		log.Printf("Storage backend %q not implemented yet; falling back to filesystem", cfg.Backend)
		storageAdapter = fsadapter.New(cfg.FS.DataDir)
	}

	// Create persistence manager
	manager := persistence.NewManager(cfg, storageAdapter)

	// Start the manager
	ctx := context.Background()
	apply := func(line []byte) error {
		_, _, err := commands.ExecuteLine(string(line), store)
		return err
	}
	if err := manager.Start(ctx, store, apply); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind on %s: %w", config.Addr, err)
	}
	defer listener.Close()

	log.Printf("Bitchest is running on %s\n", config.Addr)
	log.Printf("Waiting for connections...\n")

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		
		clientAddr := connection.RemoteAddr().String()
		log.Printf("New client connected: %s", clientAddr)
		
		go func() {
			if cfg.Enabled {
				handler.HandleWithPersistence(connection, store, manager)
			} else {
				handler.Handle(connection, store)
			}
			log.Printf("Client disconnected: %s", clientAddr)
		}()
	}
}

func main() {
	config := parseFlags()
	
	if err := StartServer(config); err != nil {
		log.Fatal(err)
	}
}
