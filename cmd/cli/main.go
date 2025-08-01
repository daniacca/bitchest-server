package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = "7463"
)

func readArrayElement(reader *bufio.Reader) (string, error) {
    // Read the element type line
    elementTypeLine, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    elementTypeLine = strings.TrimSpace(elementTypeLine)
    
    // Handle different element types within the array
    switch {
    case strings.HasPrefix(elementTypeLine, "$"): // Bulk string
        if elementTypeLine == "$-1" {
            return "(nil)", nil
        }
        
        // Read the actual string content
        content, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        }
        return strings.TrimSpace(content), nil
        
    case strings.HasPrefix(elementTypeLine, ":"): // Integer
        return elementTypeLine, nil
        
    case strings.HasPrefix(elementTypeLine, "+"): // Simple string
        return elementTypeLine, nil
        
    case strings.HasPrefix(elementTypeLine, "-"): // Error
        return elementTypeLine, nil
        
    default:
        return elementTypeLine, nil
    }
}

func formatArrayOutput(elements []string) string {
    if len(elements) == 0 {
        return "(empty list or set)"
    }
    
    // For MEMORY STATS, format as key-value pairs
    if len(elements) > 0 && strings.Contains(elements[0], "=") {
        var result strings.Builder
        for i, element := range elements {
            if i > 0 {
                result.WriteString("\n")
            }
            result.WriteString(element)
        }
        return result.String()
    }
    
    // For other arrays, format as numbered list
    var result strings.Builder
    for i, element := range elements {
        if i > 0 {
            result.WriteString("\n")
        }
        result.WriteString(fmt.Sprintf("%d) %s", i+1, element))
    }
    return result.String()
}

func main() {
	host := defaultHost
	port := defaultPort

	// Check for command line arguments
	if len(os.Args) > 1 {
		host = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	// Connect to server
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Printf("Error connecting to server %s:%s: %v\n", host, port, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to Bitchest server at %s:%s\n", host, port)
	fmt.Println("Type 'quit' or 'exit' to close the connection")
	fmt.Println("Type 'help' for available commands")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("bitchest> ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}

		// Handle special commands
		switch strings.ToLower(command) {
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			printHelp()
			continue
		case "clear":
			fmt.Print("\033[H\033[2J") // Clear screen
			continue
		}

		// Send command to server
		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Printf("Error sending command: %v\n", err)
			return
		}

		// Read response using a simple approach
		response, err := readSimpleResponse(conn)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		// Print response
		if response != "" {
			fmt.Println(response)
		}
	}
}

func readSimpleResponse(conn net.Conn) (string, error) {
	// Set a short timeout for reading
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	
	reader := bufio.NewReader(conn)
	
	// Try to read a line
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Timeout means no response (empty response for NX/XX failures)
			conn.SetReadDeadline(time.Time{}) // Reset deadline
			return "", nil
		}
		conn.SetReadDeadline(time.Time{}) // Reset deadline
		return "", err
	}
	
	// Reset the deadline
	conn.SetReadDeadline(time.Time{})
	
	firstLine = strings.TrimSpace(firstLine)
	
	// Handle different RESP types
	switch {
	case strings.HasPrefix(firstLine, "+"): // Simple string
		return firstLine, nil
		
	case strings.HasPrefix(firstLine, "-"): // Error
		return firstLine, nil
		
	case strings.HasPrefix(firstLine, ":"): // Integer
		return firstLine, nil
		
	case strings.HasPrefix(firstLine, "$"): // Bulk string
		if firstLine == "$-1" {
			return "(nil)", nil // Display null bulk string as (nil) like Redis
		}
		
		// Read the actual string content
		content, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(content), nil // Just return the content, not the RESP format
		
	case strings.HasPrefix(firstLine, "*"): // Array
		if firstLine == "*0" {
			return "(empty list or set)", nil // Display empty array like Redis
		}
		
		// Parse array length
		arrayLength := 0
		_, err := fmt.Sscanf(firstLine, "*%d", &arrayLength)
		if err != nil {
			return "", fmt.Errorf("invalid array format: %v", err)
		}
		
		// Read all array elements
		var elements []string
		for i := 0; i < arrayLength; i++ {
			element, err := readArrayElement(reader)
			if err != nil {
				return "", fmt.Errorf("error reading array element %d: %v", i+1, err)
			}
			elements = append(elements, element)
		}
		
		// Format the array output nicely
		return formatArrayOutput(elements), nil
		
	default:
		return firstLine, nil
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  SET key value [EX seconds] [NX|XX]  - Set a key with optional expiration and conditions")
	fmt.Println("  GET key                             - Get the value of a key")
	fmt.Println("  DEL key1 [key2 ...]                 - Delete one or more keys")
	fmt.Println("  EXISTS key1 [key2 ...]              - Check if keys exist")
	fmt.Println("  KEYS                                - List all keys")
	fmt.Println("  FLUSHALL                            - Remove all keys")
	fmt.Println("  EXPIRE key seconds                  - Set expiration for a key")
	fmt.Println("  TTL key                             - Get time to live for a key")
	fmt.Println("  PING                                - Test server connection")
	fmt.Println()
	fmt.Println("Special commands:")
	fmt.Println("  help                                - Show this help")
	fmt.Println("  clear                               - Clear screen")
	fmt.Println("  quit, exit                          - Close connection")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  SET user:123 name John")
	fmt.Println("  SET session:456 token abc123 EX 3600")
	fmt.Println("  SET counter 1 NX")
	fmt.Println("  GET user:123")
	fmt.Println("  TTL session:456")
} 