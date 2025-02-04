# Cacher

Cacher is a simple TCP server application written in Go that is designed to store and retrieve cached data efficiently. This project serves as a learning tool for working with Go's networking capabilities, handling TCP connections, and managing data in memory.

## Features

- **TCP Server**: Listens for incoming TCP connections on a specified port.
- **Data Caching**: Allows clients to store and retrieve key-value pairs.
- **Concurrency**: Handles multiple client connections simultaneously using goroutines.
- **Environment Configuration**: Uses environment variables for configuration settings.

## Project Structure Overview
### cache.go :
- Defines a Cache interface with methods like Get, Set, Delete, and Clear.
- Implements two cache types:
1. cache[K, V]: A standard map-based cache with an RWMutex for concurrency safety.
2. syncCache[K, V]: A sync.Map-based cache for high-frequency access.
- Provides a CacheManager to manage both cache types and select the appropriate one based on the frequentAccess flag.
### command.go :
- Defines Command and ExecutableCommand interfaces.
- Implements command structs (setCommand, getCommand, etc.) that encapsulate parsing and execution logic.
- Includes utility structs like commandArgument and commandOption for defining arguments and options.
### connection.go :
- Defines a Connection interface for handling client connections.
- Implements TCPConnection to wrap net.Conn and provide logging, reading, sending, and closing functionality.
### executor.go :
- Defines the Executor interface to handle command execution.
- Implements executor[K, V] to delegate execution to commands while managing the cache context.
### server.go :
- Defines a Server interface with methods for starting, shutting down, and handling connections.
- Implements server[K, V] to coordinate between the listener, worker pool, CommandManager, CacheManager, and Executor.
### main.go :
- Configures the server by reading environment variables for port and number of workers.
- Initializes the logger, CommandManager, CacheManager, and server.
Starts the server with a graceful shutdown timeout.


## Getting Started

### Prerequisites

- Go 1.13 or later
- A compatible operating system (Linux, macOS, or Windows)

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/cacher.git
2. Set Up Environment Variables:
You can set the required environment variables in your terminal or through a .env file. The required variables are:
    - CACHER_PORT: The port on which the server will listen
    - CACHER_NBR_WORKERS: the number of workers/goroutines to run

3. Install Dependencies (if any):
    ```bash
    go mod tidy
### Running the Application
1. To run the Cacher server, execute:
    ```bash
    go run .
The server will start and listen for incoming connections on the specified port.

### Sending Data to the Server
You can send data to the Cacher server using various methods:
1. Using telnet:
    ```bash
    telnet localhost <CACHER_PORT>
2. Using netcat (nc):
    ```bash
    echo "set mykey myvalue" | nc localhost <CACHER_PORT>
