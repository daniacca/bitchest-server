package handler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/daniacca/bitchest/internal/commands"
	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func Handle(conn net.Conn, store *db.InMemoryDB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Read error: %v\n", err)
			}
			return
		}

		input := strings.TrimSpace(raw)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmdName := strings.ToUpper(parts[0])
		args := parts[1:]

		cmd, found := commands.ExtractCommand(cmdName)
		if !found {
			conn.Write([]byte(protocol.Error("unknown command '" + cmdName + "'")))
			continue
		}

		output, err := cmd.Execute(args, store)
		if err != nil {
			conn.Write([]byte(protocol.Error(err.Error())))
			continue
		}

		conn.Write([]byte(output))
	}
}
