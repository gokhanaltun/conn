# CONN

A simple command-line tool for establishing TCP connections and listening for incoming connections, with optional SOCKS5 proxy support.

## Features

* Connect to a TCP server with or without a SOCKS5 proxy
* Start a TCP server to listen for incoming connections

## Prerequisites
- Ensure that `Golang` is installed on your system.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/gokhanaltun/conn
cd conn
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the project:

```bash
go build -o conn
```

## Usage

### Connect to a TCP Address

Connect to a TCP server at the specified address. Optionally, use a SOCKS5 proxy:

```bash
./conn c <address> -s5 <proxy-address>
```

* `<address>` - Target server address (e.g., `localhost:8080`)
* `-s5` - Optional SOCKS5 proxy address (e.g., `127.0.0.1:1080`)

**Examples:**

* Connect without proxy:

  ```bash
  ./conn c localhost:8080
  ```

* Connect using a SOCKS5 proxy:

  ```bash
  ./conn c localhost:8080 -s5 127.0.0.1:1080
  ```

### Start a TCP Server

Listen for incoming connections on a specified port:

```bash
./conn l <port>
```

* `<port>` - Port to listen on (e.g., `8080`)

**Example:**

```bash
./conn l 8080
```
