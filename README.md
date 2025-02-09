# Cacher

Cacher is a high-performance TCP server application written in Go that provides efficient in-memory caching capabilities. This project demonstrates how to work with Go's networking features, handle concurrent TCP connections, and implement modular designs for data management and expiration tracking.

## Features

- **TCP Server**: Listens for incoming TCP connections on a specified port.
- **Data Caching**: Supports storing and retrieving key-value pairs with optional TTL (time-to-live) and frequent-access optimizations.
- **Concurrency**: Handles multiple client connections simultaneously using goroutines and Go's built-in concurrency primitives.
- **Environment Configuration**: Configures the server using environment variables (`CACHER_PORT`, `CACHER_NBR_WORKERS`).
- **Logging**: Provides detailed logs for all server, cache, command, and connection events.
- **Error Handling**: Implements robust error handling with custom error types for invalid commands, unexpected issues, and graceful shutdowns.
- **Modular Design**: Separates concerns into distinct modules (`cache`, `command`, `connection`, `executor`, `server`) for maintainability and extensibility.

## Project Structure Overview

### `cache.go`
- Defines the `Cache` interface with methods like `Get`, `Set`, `Delete`, and `Clear`.
- Implements two cache types:
  1. `cache[K, V]`: A standard map-based cache with an `RWMutex` for thread-safe operations.
  2. `syncCache[K, V]`: A `sync.Map`-based cache optimized for high-frequency access.
- Provides a `CacheManager` to manage both cache types and dynamically select the appropriate one based on the `frequentAccess` flag.

### `records.go`
- Defines the `Records` interface for managing expiration times of cached entries.
- Implements `records[K]` to group keys by their expiration timestamps (with configurable precision).
- Provides efficient batch deletion of expired keys using callbacks for integration with the cache.

### `command.go`
- Defines `Command` and `ExecutableCommand` interfaces for defining and executing commands.
- Implements command structs (e.g., `setCommand`, `getCommand`) that encapsulate parsing and execution logic.
- Includes utility structs like `commandArgument` and `commandOption` for defining arguments and options.

### `connection.go`
- Defines the `Connection` interface for handling client connections.
- Implements `TCPConnection` to wrap `net.Conn` and provide logging, reading, sending, and closing functionality.
- Logs all incoming and outgoing messages with remote address information.

### `executor.go`
- Defines the `Executor` interface for delegating command execution to the appropriate cache.
- Implements `executor[K, V]` to manage the execution context (e.g., selecting the correct cache type) and invoking the command's `Run` method.

### `janitor.go`
- Defines the `Janitor` interface for periodic cleanup of expired cache entries.
- Implements `janitor[K, V]` to manage cleanup intervals and invoke the cache's `ClearExpired` method.

### `server.go`
- Defines the `Server` interface with methods for starting, shutting down, and handling connections.
- Implements `server[K, V]` to coordinate between the listener, worker pool, `CommandManager`, `CacheManager`, and `Executor`.
- Ensures graceful shutdown with a timeout mechanism.

### `main.go`
- Configures the server by reading environment variables for port and number of workers.
- Initializes the logger, `CommandManager`, `CacheManager`, and `Server`.
- Starts the server with a graceful shutdown timeout.

## Getting Started

### Prerequisites
- **Go 1.13 or later**: Ensure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).
- **Compatible Operating System**: Linux, macOS, or Windows.

### Installation
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/cacher.git
   cd cacher
   ```

2. **Set Up Environment Variables**:
   Configure the server using the following environment variables:
   - `CACHER_PORT`: The port on which the server will listen (default: `8080`).
   - `CACHER_NBR_WORKERS`: The number of worker goroutines to handle connections (default: `10`).

   Example:
   ```bash
   export CACHER_PORT=9090
   export CACHER_NBR_WORKERS=20
   ```

3. **Install Dependencies**:
   If there are external dependencies, install them using:
   ```bash
   go mod tidy
   ```

4. **Build the Application**:
   Compile the server binary:
   ```bash
   go build .
   ```

5. **Run the Application**:
   Start the server:
   ```bash
   ./cacher
   ```

   Alternatively, run directly:
   ```bash
   go run .
   ```

   The server will start and listen for incoming connections on the specified port.

---

## Sending Data to the Server

You can interact with the Cacher server using various tools:

1. **Using Telnet**:
   ```bash
   telnet localhost <CACHER_PORT>
   ```
   Example:
   ```bash
   telnet localhost 9090
   set mykey myvalue
   get mykey
   ```

2. **Using Netcat (`nc`)**:
   ```bash
   echo "set mykey myvalue" | nc localhost <CACHER_PORT>
   ```

3. **Using Custom Clients**:
   Write a Go or Python script to send TCP commands to the server.

---

## Logging

The server uses a centralized logging system to track all events. Logs are written to a file (`server.log`) with the following levels:
- **Info**: General operational messages.
- **Warning**: Non-critical issues or potential problems.
- **Error**: Critical failures or unexpected conditions.

Example log entry:
```
[2023-10-01 12:00:00] [INFO] New connection from 127.0.0.1:34567
[2023-10-01 12:00:01] [INFO] [127.0.0.1:34567] > set mykey myvalue
[2023-10-01 12:00:01] [INFO] [127.0.0.1:34567] < OK
```

---

## Commands

The server supports the following commands:

1. **SET**:
   - Syntax: `SET key value [TTLInSeconds]`
   - Example: `SET mykey myvalue 60` (sets `mykey` with a TTL of 60 seconds).

2. **GET**:
   - Syntax: `GET key`
   - Example: `GET mykey`.

3. **DEL**:
   - Syntax: `DEL key`
   - Example: `DEL mykey`.

4. **FLUSH**:
   - Syntax: `FLUSH`
   - Clears all cached data.

---

## Error Handling

The server implements a robust error-handling system with custom error types:
- **InvalidCommandError**: Raised when an unrecognized command is received.
- **InvalidCommandUsageError**: Raised when a command is used incorrectly.
- **UnexpectedError**: Captures unexpected issues during connection handling or cache operations.
- **CommandNotExecutableError**: Raised when a command doesn't implement the `ExecutableCommand` interface.

Errors are logged and sent back to the client as plain-text responses.

---

## Modular Design

The Cacher project is designed with modularity in mind:
- **Separation of Concerns**: Each module focuses on a specific responsibility (e.g., caching, command parsing, connection handling).
- **Dynamic Configuration**: Cache precision and janitor intervals can be adjusted dynamically.
- **Extensibility**: Adding new commands or cache types requires minimal changes to existing code.

---

## Contributing

To contribute to the project:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with clear documentation and tests.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Future Enhancements

- **Metrics and Monitoring**: Add support for metrics (e.g., hit rate, request count) and integrate with monitoring tools.
- **Persistent Storage**: Extend the cache to support persistence for long-term storage.
- **Advanced Expiration Policies**: Implement additional eviction strategies (e.g., LRU, MRU).

---

This README provides a comprehensive overview of the Cacher project, its features, and how to use it. It also highlights the modular design and future enhancements, making it easier for contributors and users to understand and extend the application.