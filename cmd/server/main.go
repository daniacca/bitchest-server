package main

import (
	"fmt"
	"log"
	"net"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/handler"
)

const defaultPort = ":7463"

func StartServer(addr string) error {
	store := db.NewDB()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind on %s: %w", addr, err)
	}
	defer ln.Close()

	fmt.Println("Bitchest is running on", addr)

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
	if err := StartServer(defaultPort); err != nil {
		log.Fatal(err)
	}
}
