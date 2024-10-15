# ServerHTTP

![Version](https://img.shields.io/badge/Version-0.1.0-brightgreen)
![Language](https://img.shields.io/badge/Language-go-blue)
![License](https://img.shields.io/badge/License-GPL_3.0-red)
---
## Project Description

HTTP is the protocol that powers the web. This is an HTTP server that's capable of handling simple GET/POST requests, serving files, and handling multiple concurrent connections.

Along the way, I learnt about:
- TCP connections
- HTTP headers
- HTTP verbs (GET, POST)
- Handling multiple connections
- Serving files from a directory
- Gzip compression

## Features

- **Handle GET and POST requests**: The server can handle basic GET and POST requests.
- **Serve files**: The server can serve files from a specified directory.
- **Concurrent connections**: The server can handle multiple connections simultaneously using goroutines.
- **Gzip compression**: The server supports gzip compression for responses when requested by the client.

## Getting Started

### Prerequisites

- Go (version 1.16 or later)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/mivige/ServerHTTP.git
    cd ServerHTTP
    ```

2. Build the project:
    ```sh
    go build -o server.go
    ```

### Usage

1. Run the server with the `--directory` flag to specify the directory to serve files from:
    ```sh
    ./your_program --directory /path/to/your/directory
    ```

2. The server will start listening on port `4221`.

### Examples

#### GET Request

To get a simple response:

```sh
curl -v http://localhost:4221/
```

To echo a message:

```sh
curl -v http://localhost:4221/echo/your-message
```

To get the user-agent:

```sh
curl -v http://localhost:4221/user-agent
```

To get a file:

```sh
curl -v http://localhost:4221/files/your-filename
```

#### POST Request

To create a file:

```sh
curl -v --data "file content" -H "Content-Type: application/octet-stream" http://localhost:4221/files/your-filename
```

### Compression

To request gzip compression:

```sh
curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/your-message
```

## Code Overview

### Main Function

The `main` function sets up the server, parses command-line arguments, and starts listening for connections.

### handleConnection Function

The `handleConnection` function handles individual connections, parses requests, and constructs appropriate responses.

### handleFileCreation Function

The `handleFileCreation` function handles file creation for POST requests.

### handleFileRequest Function

The `handleFileRequest` function handles file requests for GET requests.

### compressGzip Function

The `compressGzip` function compresses data using gzip.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

This project was made following the [CodeCrafters](https://app.codecrafters.io/catalog) challenge: **Build your own HTTP server**.
