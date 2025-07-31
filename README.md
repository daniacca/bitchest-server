# Bitchest

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![CI](https://github.com/daniacca/bitchest/actions/workflows/release.yml/badge.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/kaelisra/bitchest)
![Coverage](https://img.shields.io/badge/coverage-62.9%25-brightgreen)

<img src="doc/img/bitchest_logo.png" alt="Bitchest Logo" width="200">

**Bitchest** is a lightweight in-memory key-value database written in Go, inspired by the core ideas of Redis, but designed with simplicity, clarity, and educational value in mind.

It supports plain-text TCP connections and a minimal set of commands for managing string values. The project is modular, testable, and easy to extend.

**Main Features:**

- ✅ Built-in CLI client
- ✅ SET and GET strings data type
- ✅ Expiration support with TTL
- ✅ Configurable server settings (host and port)

---

## 🚀 Full Features

- **RESP Protocol Compliance**: Full Redis Serialization Protocol support
- **Null Response Handling**: Proper `$-1\r\n` format for nil results
- **Interactive CLI Client**: Built-in command-line interface
- **Key Expiration**: TTL support with `EXPIRE` and `TTL` commands
- **Conditional Operations**: `NX` (set if Not eXists) and `XX` (set if eXists) options
- **Comprehensive Test Coverage**: Unit tests for all commands and database operations
- **Configurable Server Settings**: Command-line flags for host and port
- **Docker Support**: Containerized deployment with proper networking
- **Server Logging**: Detailed logging for client connections and command processing

---

## 🚀 Supported Commands

| Command                               | Description                                                                   |
| ------------------------------------- | ----------------------------------------------------------------------------- |
| `SET key value [EX seconds] [NX\|XX]` | Sets a key with a string value (optional expiration and existence conditions) |
| `GET key`                             | Retrieves the value associated with a key                                     |
| `DEL key...`                          | Deletes one or more keys                                                      |
| `EXISTS key...`                       | Checks if one or more keys exist                                              |
| `KEYS`                                | Returns all current keys                                                      |
| `FLUSHALL`                            | Removes all keys from the database                                            |
| `EXPIRE key seconds`                  | Sets an expiration time for a key in seconds                                  |
| `TTL key`                             | Returns the time to live for a key in seconds                                 |

### SET Command Options

- **`EX seconds`**: Set expiration time in seconds
- **`NX`**: Only set the key if it does not already exist
- **`XX`**: Only set the key if it already exists

---

## 💻 Local Development

```bash
make               # Starts the server locally on port 7463
make run-port      # Starts the server on port 6379
make run-host      # Starts the server on all interfaces (0.0.0.0)
make run-host-port # Starts the server on all interfaces:6379
make build         # Builds the server binary
make build-cli     # Builds the CLI client
make build-all     # Builds both server and CLI
make test          # Runs all unit tests
make help          # Shows server command-line options
```

### Git Hooks (Lefthook)

The project uses **Lefthook** for git hooks to enforce code quality:

```bash
# Install lefthook hooks
npm run prepare

# Run lefthook manually
npm run lefthook

# Check commit message format
npm run commitlint:check
```

**Hooks include:**

- **commit-msg**: Validates conventional commit format
- **pre-commit**: Runs `go fmt`, `go vet`, and tests on staged Go files
- **pre-push**: Runs full test suite and build verification

---

## ⚙️ Server Configuration

The Bitchest server supports command-line configuration:

```bash
# Default configuration (localhost:7463)
./out/bitchest

# Custom port
./out/bitchest -port 6379

# Bind to all interfaces
./out/bitchest -host 0.0.0.0

# Custom host and port
./out/bitchest -host 0.0.0.0 -port 6379

# Show help
./out/bitchest -h
```

### Configuration Options

| Flag    | Default     | Description                |
| ------- | ----------- | -------------------------- |
| `-host` | `localhost` | Host to bind the server to |
| `-port` | `7463`      | Port to bind the server to |

---

## 🖥️ CLI Client

Bitchest includes a built-in CLI client for easy interaction with the server:

```bash
# Build the CLI
make build-cli

# Connect to default server (localhost:7463)
./out/bitchest-cli

# Connect to custom host and port
./out/bitchest-cli localhost 6379
```

### CLI Features

- **Interactive Mode**: Type commands directly in the terminal
- **Built-in Help**: Type `help` to see available commands
- **Special Commands**:
  - `help` - Show command help
  - `clear` - Clear screen
  - `quit` or `exit` - Close connection
- **RESP Protocol Support**: Handles all Bitchest response types
- **Redis-like Output**: Displays `(nil)` for null responses, `(empty list or set)` for empty arrays

---

## 🐳 Docker

### Build

```bash
make docker-build
```

### Run

```bash
make docker-run
```

The server will be available at `localhost:7463`.

**Note**: The Docker container automatically binds to `0.0.0.0:7463` to allow external access from the host and other containers.

### Custom Docker Run

```bash
# Run with custom port mapping
docker run --rm -p 6379:7463 bitchest:latest

# Run with custom host and port
docker run --rm -p 0.0.0.0:6379:7463 bitchest:latest

# Run in background
docker run -d --name bitchest-server -p 7463:7463 bitchest:latest
```

---

## 📊 Server Logging

The server provides comprehensive logging for monitoring and debugging:

### Connection Logging

- **Client Connections**: Logs when new clients connect with their IP address
- **Client Disconnections**: Logs when clients disconnect
- **Connection Errors**: Logs any connection-related errors

### Command Logging

- **Command Reception**: Logs every command received from clients
- **Execution Timing**: Measures and logs command execution time
- **Success Logging**: Logs successful command completions with timing
- **Error Logging**: Logs command failures with detailed error messages

### Log Format

```
2025/07/29 11:12:27 New client connected: 127.0.0.1:44422
2025/07/29 11:12:27 [127.0.0.1:44422] Received command: SET testkey testvalue
2025/07/29 11:12:27 [127.0.0.1:44422] Command 'SET' completed successfully in 31.907µs
2025/07/29 11:12:27 [127.0.0.1:44422] Client disconnected
```

### Example Output

```bash
$ make run
Bitchest is running on localhost:7463
Waiting for connections...
2025/07/29 11:12:27 New client connected: 127.0.0.1:44422
2025/07/29 11:12:27 [127.0.0.1:44422] Received command: PING
2025/07/29 11:12:27 [127.0.0.1:44422] Command 'PING' completed successfully in 1.472µs
2025/07/29 11:12:27 [127.0.0.1:44422] Received command: SET key value
2025/07/29 11:12:27 [127.0.0.1:44422] Command 'SET' completed successfully in 25.123µs
2025/07/29 11:12:27 [127.0.0.1:44422] Client disconnected
```

---

## 🧪 Testing

All components are covered by unit tests:

- in-memory database
- communication protocol
- command implementations
- TCP handler
- server startup (`StartServer`)
- expiration functionality
- existence conditions (NX/XX)
- RESP protocol compliance

---

## 📦 Future Plans

- Advanced types (`LIST`, `ZSET`)
- Server status commands (i.e. `MEMORY STATS`)
- DB persistence (file, storage/bucket)
- Cluster configuration
- Scripting for advanced commands

---

## 📄 License

Distributed under the [MIT](./LICENSE) license.

## 📊 Code Coverage

Bitchest maintains comprehensive test coverage to ensure code quality and reliability.

### Current Coverage Status

**Overall Coverage: 62.9%**

| Package             | Coverage | Status               |
| ------------------- | -------- | -------------------- |
| `cmd/cli`           | 0.0%     | ❌ No coverage       |
| `cmd/server`        | 35.0%    | ⚠️ Needs improvement |
| `internal/commands` | 96.5%    | ✅ Excellent         |
| `internal/db`       | 86.7%    | ✅ Good              |
| `internal/handler`  | 84.4%    | ✅ Good              |
| `internal/protocol` | 100.0%   | ✅ Perfect           |

### Coverage Commands

```bash
# Generate coverage report
make coverage

# Show coverage summary
make coverage-summary

# Clean coverage files
make coverage-clean
```

### Coverage Reports

- **HTML Report**: `coverage.html` - Detailed visual coverage report
- **Function Report**: `coverage.out` - Raw coverage data
- **Summary**: Shows per-function and per-package coverage

### Coverage Threshold

The project maintains a **60% minimum coverage threshold** enforced by CI/CD pipelines.

### Coverage Workflow

GitHub Actions automatically:

- ✅ Runs tests with coverage on every PR
- ✅ Generates HTML and function reports
- ✅ Uploads coverage artifacts
- ✅ Comments PRs with coverage summary
- ✅ Enforces 60% coverage threshold

---

## ⚙️ Example Usage

### Using the CLI Client

```bash
$ ./out/bitchest-cli
Connected to Bitchest server at localhost:7463
Type 'quit' or 'exit' to close the connection
Type 'help' for available commands

bitchest> SET user:123 John
+OK
bitchest> GET user:123
John
bitchest> SET counter 1 NX
+OK
bitchest> SET counter 2 NX
(nil)
bitchest> SET counter 3 XX
+OK
bitchest> GET nonexistent
(nil)
bitchest> KEYS
*2
bitchest> quit
Goodbye!
```

### Using netcat

```bash
$ nc localhost 7463

SET hello world
+OK

GET hello
$5
world

SET session user123 EX 60
+OK

TTL session
:60

EXPIRE hello 30
:1

SET counter 1 NX
+OK

SET counter 2 NX
$-1

SET counter 3 XX
+OK

SET newkey value XX
$-1

SET tempkey temp EX 10 NX
+OK

KEYS
*4
$5
hello
$7
session
$7
counter
$7
tempkey
```

---
