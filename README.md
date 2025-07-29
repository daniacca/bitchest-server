# Bitchest

<img src="doc/img/bitchest_logo.png" alt="Bitchest Logo" width="200">

**Bitchest** is a lightweight in-memory key-value database written in Go, inspired by the core ideas of Redis, but designed with simplicity, clarity, and educational value in mind.

It supports plain-text TCP connections and a minimal set of commands for managing string values. The project is modular, testable, and easy to extend.

**Features:**
- ✅ Full RESP protocol compliance
- ✅ Proper null response handling (`$-1\r\n`)
- ✅ Built-in CLI client
- ✅ Expiration support with TTL
- ✅ Conditional operations (NX/XX)
- ✅ Comprehensive test coverage
- ✅ Configurable server settings

---

## 🚀 Supported Commands

| Command                      | Description                                          |
|------------------------------|------------------------------------------------------|
| `SET key value [EX seconds] [NX\|XX]` | Sets a key with a string value (optional expiration and existence conditions) |
| `GET key`                    | Retrieves the value associated with a key            |
| `DEL key...`                 | Deletes one or more keys                             |
| `EXISTS key...`              | Checks if one or more keys exist                     |
| `KEYS`                       | Returns all current keys                             |
| `FLUSHALL`                   | Removes all keys from the database                   |
| `EXPIRE key seconds`         | Sets an expiration time for a key in seconds         |
| `TTL key`                    | Returns the time to live for a key in seconds        |

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

| Flag | Default | Description |
|------|---------|-------------|
| `-host` | `localhost` | Host to bind the server to |
| `-port` | `7463` | Port to bind the server to |

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
- Additional commands (`SETNX`, etc.)
- File-based persistence
- Full RESP protocol support

---

## 📄 License

Distributed under the [MIT](./LICENSE) license.
