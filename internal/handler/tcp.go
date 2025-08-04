package handler

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/daniacca/bitchest/internal/commands"
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/parser"
	"github.com/daniacca/bitchest/internal/protocol"
)

func Handle(conn net.Conn, store *db.InMemoryDB) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	
	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("[%s] Client disconnected", clientAddr)
			} else {
				log.Printf("[%s] Read error: %v", clientAddr, err)
			}
			return
		}

		input := strings.TrimSpace(raw)
		if input == "" {
			continue
		}

		// Log incoming command
		log.Printf("[%s] Received command: %s", clientAddr, input)

		// Tokenize the input
		parts, err := parser.Tokenize(input)
		if err != nil {
			conn.Write([]byte(protocol.Error(err.Error())))
			log.Printf("[%s] Command error: %s", clientAddr, err)
			continue
		}

		if len(parts) == 0 {
			continue
		}
		
		cmdName := strings.ToUpper(parts[0])
		args := parts[1:]

		// Start timing the command execution
		startTime := time.Now()

		cmd, found := commands.ExtractCommand(cmdName)
		if !found {
			errorMsg := "unknown command '" + cmdName + "'"
			conn.Write([]byte(protocol.Error(errorMsg)))
			log.Printf("[%s] Command error: %s", clientAddr, errorMsg)
			continue
		}

		output, err := cmd.Execute(args, store)
		executionTime := time.Since(startTime)

		if err != nil {
			conn.Write([]byte(protocol.Error(err.Error())))
			log.Printf("[%s] Command '%s' failed after %v: %v", clientAddr, cmdName, executionTime, err)
			continue
		}

		conn.Write([]byte(output))
		log.Printf("[%s] Command '%s' completed successfully in %v", clientAddr, cmdName, executionTime)
	}
}
