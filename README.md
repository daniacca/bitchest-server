# Bitchest

<img src="doc/img/bitchest_logo.png" alt="Bitchest Logo" width="200">

**Bitchest** is a lightweight in-memory key-value database written in Go, inspired by the core ideas of Redis, but designed with simplicity, clarity, and educational value in mind.

It supports plain-text TCP connections and a minimal set of commands for managing string values. The project is modular, testable, and easy to extend.

---

## 🚀 Supported Commands

| Command                      | Description                                          |
|------------------------------|------------------------------------------------------|
| `SET key value [EX seconds]` | Sets a key with a string value (optional expiration) |
| `GET key`                    | Retrieves the value associated with a key            |
| `DEL key...`                 | Deletes one or more keys                             |
| `EXISTS key...`              | Checks if one or more keys exist                     |
| `KEYS`                       | Returns all current keys                             |
| `FLUSHALL`                   | Removes all keys from the database                   |
| `EXPIRE key seconds`         | Sets an expiration time for a key in seconds         |
| `TTL key`                    | Returns the time to live for a key in seconds        |

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

KEYS
*2
$5
hello
$7
session
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
