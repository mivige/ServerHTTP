package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
	"io"
	"strconv"
	"bytes"
	"compress/gzip"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

// Mapping of status codes to their string representations
var statusCodeToString = map[int]string{
	200: "OK",
	201: "Created",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

// Directory to serve files from, set via command line argument
var directory string

func main() {
	// Parse command line arguments to get the directory
	for i, arg := range os.Args {
		if arg == "--directory" && i+1 < len(os.Args) {
			directory = os.Args[i+1]
			break
		}
	}

	// Print a log message for debugging
	fmt.Println("Logs from your program will appear here!")

	// Start listening on port 4221
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Accept connections in a loop
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		// Handle each connection in a new goroutine
		go handleConnection(conn)
	}
}

// Handle an individual connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the request
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}
	request := string(buffer[:n])
	lines := strings.Split(request, "\r\n")
	requestLine := strings.Split(lines[0], " ")
	method := requestLine[0]
	path := requestLine[1]

	// Check for Accept-Encoding header to see if gzip is supported
	acceptsGzip := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Accept-Encoding:") {
			encodings := strings.Split(strings.TrimPrefix(line, "Accept-Encoding:"), ",")
			for _, encoding := range encodings {
				if strings.TrimSpace(encoding) == "gzip" {
					acceptsGzip = true
					break
				}
			}
			break
		}
	}

	var response string
	if method == "GET" {
		if path == "/" {
			// Handle root path
			response = getStatus(200) + "\r\n\r\n"
		} else if strings.HasPrefix(path, "/echo/") {
			// Handle /echo/{message} path
			message := strings.TrimPrefix(path, "/echo/")
			contentType := "Content-Type: text/plain\r\n"
			
			if acceptsGzip {
				// Compress the message if gzip is accepted
				compressedBody := compressGzip([]byte(message))
				contentEncoding := "Content-Encoding: gzip\r\n"
				contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(compressedBody))
				response = fmt.Sprintf("%s\r\n%s%s%s\r\n", getStatus(200), contentType, contentEncoding, contentLength)
				response += string(compressedBody)
			} else {
				// Send uncompressed message
				contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(message))
				response = fmt.Sprintf("%s\r\n%s%s\r\n%s", getStatus(200), contentType, contentLength, message)
			}
		} else if path == "/user-agent" {
			// Handle /user-agent path
			userAgent := strings.Split(strings.Split(request, "User-Agent:")[1], "\r\n")[0]
			response = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200), len(userAgent), userAgent)
		} else if strings.HasPrefix(path, "/files/") {
			// Handle /files/{filename} path
			filename := strings.TrimPrefix(path, "/files/")
			response = handleFileRequest(filename)
		} else {
			// Handle unknown paths
			response = getStatus(404) + "\r\n\r\n"
		}
	} else if method == "POST" && strings.HasPrefix(path, "/files/") {
		// Handle POST /files/{filename} path
		filename := strings.TrimPrefix(path, "/files/")
		response = handleFileCreation(filename, lines, conn)
	}

	// Write the response to the connection
	conn.Write([]byte(response))
}

// Handle file creation for POST /files/{filename}
func handleFileCreation(filename string, headers []string, conn net.Conn) string {
	var contentLength int
	for _, header := range headers {
		if strings.HasPrefix(header, "Content-Length:") {
			length, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(header, "Content-Length:")))
			if err == nil {
				contentLength = length
			}
		}
	}

	filePath := filepath.Join(directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return getStatus(500) + "\r\n\r\n"
	}
	defer file.Close()

	// Read the request body
	body := make([]byte, contentLength)
	_, err = io.ReadFull(conn, body)
	if err != nil {
		return getStatus(400) + "\r\n\r\n"
	}

	// Write the body to the file
	_, err = file.Write(body)
	if err != nil {
		return getStatus(500) + "\r\n\r\n"
	}

	return getStatus(201) + "\r\n\r\n"
}

// Handle file requests for GET /files/{filename}
func handleFileRequest(filename string) string {
	filePath := filepath.Join(directory, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return getStatus(404) + "\r\n\r\n"
	}

	contentLength := len(content)
	response := fmt.Sprintf("%s\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", getStatus(200), contentLength)
	return response + string(content)
}

// Get the status line for a given status code and text
func getStatus(statusCode int) string {
	return fmt.Sprintf("HTTP/1.1 %d %s", statusCode, statusCodeToString[statusCode])
}

// Compress data using gzip
func compressGzip(data []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(data)
	zw.Close()
	return buf.Bytes()
}
