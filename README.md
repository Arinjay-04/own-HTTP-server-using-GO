# Simple TCP Server in Go

This project implements a simple TCP server in Go that listens on port `4221`. It handles basic HTTP requests and supports file operations.

## Features

- **Root Endpoint (`/`)**: Responds with a welcome message.
- **Echo Endpoint (`/echo/{string}`)**: Echos back the string provided in the URL.
- **User-Agent Endpoint (`/user-agent`)**: Returns the User-Agent string of the request.
- **File Operations**:
  - **GET `/files/{filename}`**: Reads and returns the content of the specified text file.
  - **POST `/files/{filename}`**: Creates a new text file with the content provided in the request body.

## Requirements

- Go (version 1.16 or later)

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/tcp-server.git
   cd tcp-server
