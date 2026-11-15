package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	log.Println("listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	clientAddress := conn.RemoteAddr()
	log.Println("client connected:", clientAddress)

	reader := bufio.NewReader(conn)

	// 1. Request line
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("error:", clientAddress, err)
		return
	}

	firstLine = strings.TrimSuffix(firstLine, "\r\n")
	parts := strings.Split(firstLine, " ")
	if len(parts) != 3 {
		writeResponse(conn, 400, "Bad Request", "malformed request line\n")
		return
	}
	if !strings.HasPrefix(parts[2], "HTTP/") {
		writeResponse(conn, 400, "Bad Request", "malformed request line\n")
		return
	}

	method, path := parts[0], parts[1]

	// 2. Header
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error:", clientAddress, err)
			return
		}

		if line == "\r\n" {
			break
		}

		line = strings.TrimSuffix(line, "\r\n")
		name, value, found := strings.Cut(line, ":")
		if !found {
			writeResponse(conn, 400, "Bad Request", "header issue\n")
			return
		}

		headers[strings.ToLower(name)] = strings.TrimSpace(value)
	}

	// 3. Logging
	log.Printf("method: %s, path: %s, # of headers: %d\n", method, path, len(headers))

	// 4. Response
	body := fmt.Sprintf("you requested %s %s\n", method, path)
	writeResponse(conn, 200, "OK", body)
}

func writeResponse(conn net.Conn, status int, reason, body string) {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", status, reason, len(body), body)
	if _, err := conn.Write([]byte(response)); err != nil {
		log.Printf("error writing response for %v: %v", conn.RemoteAddr(), err)
	}
}
