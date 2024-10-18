package main

import (
    "bufio"
    "fmt"
    "net"
    "net/http"
    "os"
    "strings"
    "io/ioutil"
)

func main() {
    fmt.Println("Logs from your program will appear here!")

    // Listen on TCP port 4221
    l, err := net.Listen("tcp", "0.0.0.0:4221")
    if err != nil {
        fmt.Println("Failed to bind to port 4221:", err)
        os.Exit(1)
    }
    defer l.Close()

    // Accept connections in a loop
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err.Error())
            continue
        }
        go handleConnection(conn) // Handle each connection in a new goroutine
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    req, err := http.ReadRequest(reader)
    if err != nil {
        fmt.Fprintln(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
        return
    }

    // Allow only GET method
    // if req.Method != "GET" {
    //     fmt.Fprintln(conn, "HTTP/1.1 405 Method Not Allowed\r\nAllow: GET\r\n\r\n")
    //     return
    // }

    // Check if request is for root "/"
    if req.URL.Path == "/" {
        body := "Welcome to the root!"
        fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
        return
    }

    // Handle the /echo/ endpoint
    if strings.HasPrefix(req.URL.Path, "/echo/") {
        echoStr := strings.TrimPrefix(req.URL.Path, "/echo/")
        fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)
        return
    }

    // Handle the /user-agent endpoint
    if req.URL.Path == "/user-agent" {
        userAgent := req.UserAgent()
        fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
        return
    }

    // Handle the /files/{filename} endpoint
    if strings.HasPrefix(req.URL.Path, "/files/") {
		file := strings.TrimPrefix(req.URL.Path, "/files/")
		file = strings.TrimSpace(file) + ".txt"
	
		if file == "" {
			fmt.Fprintln(conn, "HTTP/1.1 400 Bad Request\r\n\r\nFilename is required")
			return
		}
	
		if req.Method == "POST" {
			// Read the body of the POST request
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Fprintln(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\nError reading request body")
				return
			}
	
			// Create a new file
			newfile, err := os.Create(file)
			if err != nil {
				fmt.Fprintln(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\nError creating file")
				return
			}
			defer newfile.Close() // Ensure the file is closed after writing
	
			// Write the body content to the new file
			length, err := newfile.Write(body)
			if err != nil {
				fmt.Fprintln(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\nError writing to file")
				return
			}
	
			// Log the length of the written content
			fmt.Println("Length is:", length)
	
			// Respond to the client indicating success
			fmt.Fprintf(conn, "HTTP/1.1 201 Created\r\nContent-Length: 0\r\n\r\n")
			return
		} else {
			// Handle GET request to read the file
			content, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintln(conn, "HTTP/1.1 404 Not Found\r\n\r\nFile not found")
				return
			}
	
			// Respond with the file content
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)
			return
		}
	}
	



    // If no endpoint matches, return 404
    fmt.Fprintln(conn, "HTTP/1.1 404 Not Found\r\n\r\nPage not found")
}
