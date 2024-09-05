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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a WebSocket URL as the first argument.")
		return
	}
	log.Println("Connecting to WebSocket server:", os.Args[1])
	// Parse the WebSocket URL
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("Invalid WebSocket URL:", err)
	}

	// Extract the username and password from the URL
	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	log.Println("Username:", username)
	log.Println("Password:", password)

	// Connect to the WebSocket server with Basic Auth credentials
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
	// strip the username and password from the URL
	u.User = nil
	log.Println("URL:", u.String())
	conn, _, err := dialer.Dial(u.String(), requestHeader)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer conn.Close()

	// Start a goroutine to handle incoming messages
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

	// Start a goroutine to handle user input
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

	// Wait for a termination signal to exit the program
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
