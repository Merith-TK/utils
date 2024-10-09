/*
Package main provides a simple WebSocket client that connects to a WebSocket server, optionally authenticating using basic authentication if credentials are embedded in the URL. The client supports sending messages from the command line and receiving messages from the server concurrently.

Usage:

	go run main.go <WebSocket_URL>

The URL can include username and password for basic authentication, for example:

	ws://user:password@localhost:8080/ws

The program also listens for system signals (e.g., SIGINT, SIGTERM) to gracefully shut down the WebSocket connection.

Functions:
*/
package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

/*
main is the entry point of the application. It checks for the presence of a WebSocket URL passed as a command-line argument, connects to the server, and starts two goroutines:
  - One for handling incoming WebSocket messages from the server.
  - One for sending messages to the WebSocket server from the user's input.

The program terminates when it receives an interrupt or termination signal (SIGINT, SIGTERM).

If a username and password are included in the WebSocket URL, they are extracted and used for basic authentication when connecting.
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a WebSocket URL as the first argument.")
		return
	}
	log.Println("Connecting to WebSocket server:", os.Args[1])

	// Parse the WebSocket URL.
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("Invalid WebSocket URL:", err)
	}

	// Extract username and password from the URL, if present.
	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	log.Println("Username:", username)
	log.Println("Password:", password)

	// Prepare to connect with optional Basic Auth credentials.
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 0,
	}
	requestHeader := http.Header{}
	if username != "" && password != "" {
		auth := username + ":" + password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		requestHeader.Set("Authorization", basicAuth)
	}

	// Strip the user info from the URL.
	u.User = nil
	log.Println("URL:", u.String())

	// Connect to the WebSocket server.
	conn, _, err := dialer.Dial(u.String(), requestHeader)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer conn.Close()

	// Goroutine for handling incoming messages from the server.
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Println("Received:", string(message))
		}
	}()

	// Goroutine for handling user input and sending it to the server.
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			err := conn.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
		}
	}()

	// Wait for a termination signal (SIGINT or SIGTERM).
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
