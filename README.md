# Bitchest

<img src="doc/img/bitchest_logo.png" alt="Bitchest Logo" width="200">

**Bitchest** is a lightweight in-memory key-value database written in Go, inspired by the core ideas of Redis, but designed with simplicity, clarity, and educational value in mind.

It supports plain-text TCP connections and a minimal set of commands for managing string values. The project is modular, testable, and easy to extend.

---

## 🚀 Supported Commands

| Command                                 | Description                                                                   |
|-----------------------------------------|-------------------------------------------------------------------------------|
| `SET key value [EX seconds] [NX \| XX]` | Sets a key with a string value (optional expiration and existence conditions) |
| `GET key`                               | Retrieves the value associated with a key                                     |
| `DEL key...`                            | Deletes one or more keys                                                      |
| `EXISTS key...`                         | Checks if one or more keys exist                                              |
| `KEYS`                                  | Returns all current keys                                                      |
| `FLUSHALL`                              | Removes all keys from the database                                            |
| `EXPIRE key seconds`                    | Sets an expiration time for a key in seconds                                  |
| `TTL key`                               | Returns the time to live for a key in seconds                                 |

### SET Command Options

- **`EX seconds`**: Set expiration time in seconds
- **`NX`**: Only set the key if it does not already exist
- **`XX`**: Only set the key if it already exists

---

## 💻 Local Development

```bash
make            # Starts the server locally on port 7463
make build      # Builds the local binary
make test       # Runs all unit tests
```

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

From terminal:

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
(nil)

SET counter 3 XX
+OK

SET newkey value XX
(nil)

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

---

## 📦 Future Plans

- Advanced types (`LIST`, `ZSET`)
- Additional commands (`SETNX`, etc.)
- File-based persistence
- Built-in CLI client
- Full RESP protocol support

---

## 📄 License

Distributed under the [MIT](./LICENSE) license.
