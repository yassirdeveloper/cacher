# Cacher

Cacher is a simple TCP server application written in Go that is designed to store and retrieve cached data efficiently. This project serves as a learning tool for working with Go's networking capabilities, handling TCP connections, and managing data in memory.

## Features

- **TCP Server**: Listens for incoming TCP connections on a specified port.
- **Data Caching**: Allows clients to store and retrieve key-value pairs.
- **Concurrency**: Handles multiple client connections simultaneously using goroutines.
- **Environment Configuration**: Uses environment variables for configuration settings.

## Getting Started

### Prerequisites

- Go 1.13 or later
- A compatible operating system (Linux, macOS, or Windows)

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/cacher.git
2. Set Up Environment Variables:
You can set the required environment variables in your terminal or through a .env file. The main variable is:
CACHER_PORT: The port on which the server will listen (default is 6969).

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
