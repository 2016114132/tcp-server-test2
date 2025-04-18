package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Define a command line flag named "port" with a default value of "4000"
	port := flag.String("port", "4000", "Port to listen on")
	flag.Parse()

	// Construct the address string in the form ":port"
	address := ":" + *port

	// Start listening for TCP connections on the specified address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", *port)

	// Our program runs an infinite loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Get client's address and log it
	clientAddress := conn.RemoteAddr().String()
	log.Printf("Client Connected: %s at %s", clientAddress, time.Now().Format(time.RFC3339))

	// Extract only the IP (no port) to name the log file
	host, _, _ := net.SplitHostPort(clientAddress)
	logFileName := filepath.Join("logs", host+".log")

	// Ensure the logs/ folder exists
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Printf("Could not create logs directory: %v", err)
		return
	}

	// Open the log file for writing messages
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file for %s: %v", host, err)
		return
	}
	defer logFile.Close()

	// Read data line by line
	reader := bufio.NewReader(conn)

	for {
		// Set a deadline for how long the server will wait for input from the client
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		message, err := reader.ReadString('\n')
		if err != nil {
			netErr, ok := err.(net.Error)

			if ok && netErr.Timeout() {
				log.Printf("Client %s disconnected due to inactivity at %s", clientAddress, time.Now().Format(time.RFC3339))
			} else if errors.Is(err, io.EOF) {
				log.Printf("Client Disconnected: %s at %s", clientAddress, time.Now().Format(time.RFC3339))
			} else {
				log.Printf("Error from %s: %v", clientAddress, err)
			}
			return
		}

		trimmed := strings.TrimSpace(message)

		// Define the max length of the message allowed
		maxLength := 1024
		// Check if the trimmed message exceeds the maximum allowed size
		if len(trimmed) > maxLength {
			// Tell the client that the message was too long
			_, err = conn.Write([]byte("Error: message too long\n"))
			if err != nil {
				log.Printf("Error sending overflow warning to %s: %v", clientAddress, err)
			}
			// Print that message was rejected
			log.Printf("Client %s sent an oversized message. Rejected.", clientAddress)
			continue
		}

		// Personality mode response for exact messages
		switch trimmed {
		case "":
			// If client sent an empty message
			_, err = conn.Write([]byte("Say something...\n"))
			if err != nil {
				log.Printf("Error writing to client %s: %v", clientAddress, err)
			}
			continue

		case "hello":
			// If client sent hello, we respond with "Hi there!"
			_, err = conn.Write([]byte("Hi there!\n"))
			if err != nil {
				log.Printf("Error writing to client %s: %v", clientAddress, err)
			}
			continue

		case "bye":
			// If client sent bye, we respond "Goodby!" and close the connection
			_, err = conn.Write([]byte("Goodbye!\n"))
			if err != nil {
				log.Printf("Error writing to client %s: %v", clientAddress, err)
			}
			log.Printf("Client %s said 'bye' — closing connection", clientAddress)
			return
		}

		// Handle commands starting with a slash (e.g., /time, /quit, /echo msg)
		if strings.HasPrefix(trimmed, "/") {
			fields := strings.Fields(trimmed)

			switch fields[0] {
			case "/time":
				currentTime := time.Now().Format("Mon Jan 2 15:04:05 2006")
				conn.Write([]byte(currentTime + "\n"))
				continue

			case "/quit":
				conn.Write([]byte("Closing connection...\n"))
				log.Printf("Client %s issued /quit — closing connection", clientAddress)
				return

			case "/echo":
				// Check if there's a message after /echo
				if len(fields) > 1 {
					messagePart := strings.Join(fields[1:], " ")
					conn.Write([]byte(messagePart + "\n"))
				} else {
					conn.Write([]byte("Error: Missing message for /echo\n"))
				}
				continue

			default:
				// Unrecognized command
				conn.Write([]byte("Error: Unknown command\n"))
				continue
			}
		}

		// Echo the message back to the client
		_, err = conn.Write([]byte(trimmed + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}

		// Log to client's file
		logLine := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), trimmed)
		logFile.WriteString(logLine)
	}
}
