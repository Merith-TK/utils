package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Merith-TK/utils/pkg/debug"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: false, // Disable compression
	}
	startTime       time.Time
	activeConns     = make(map[*websocket.Conn]time.Time) // Track active connections
	connsMutex      sync.Mutex                            // Mutex for safe concurrent access
	connCheckPeriod = 10 * time.Second                    // Period to check for broken connections
)

func main() {
	flag.Parse()
	startTime = time.Now()
	log.Println("Downtime Server started")

	// Define two different WebSocket endpoints
	http.HandleFunc("/", handleWebSocket)          // For ping/pong and misc commands
	http.HandleFunc("/heartbeat", handleHeartbeat) // For heartbeat

	// Start the HTTP server
	go func() {
		for {
			time.Sleep(connCheckPeriod)
			checkConnections()
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle the main WebSocket connection at "/"
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	debug.Print("New connection from:", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	debug.Print("Connection upgraded successfully at '/'")

	connsMutex.Lock()
	activeConns[conn] = time.Now()
	connsMutex.Unlock()

	defer func() {
		conn.Close()
		connsMutex.Lock()
		delete(activeConns, conn)
		connsMutex.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		// Handle "ping" message with "pong" response
		if string(message) == "ping" {
			err = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			if err != nil {
				log.Println("Failed to send pong:", err)
				break
			}
		}

		// Handle uptime request
		if string(message) == "uptime" {
			uptime := time.Since(startTime).String()
			err = conn.WriteMessage(websocket.TextMessage, []byte(uptime))
			if err != nil {
				log.Println("Failed to send uptime:", err)
				break
			}
		}
	}
}

// Handle the heartbeat WebSocket connection at "/heartbeat"
func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	debug.Print("New heartbeat connection from:", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade heartbeat connection:", err)
		return
	}
	debug.Print("Connection upgraded successfully at '/heartbeat'")

	connsMutex.Lock()
	activeConns[conn] = time.Now()
	connsMutex.Unlock()

	defer func() {
		conn.Close()
		connsMutex.Lock()
		delete(activeConns, conn)
		connsMutex.Unlock()
	}()

	// Send heartbeat message every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte("Heartbeat"))
			if err != nil {
				log.Println("Failed to send heartbeat:", err)
				break
			}
		}
	}
}

// checkConnections periodically checks the status of connections
func checkConnections() {
	connsMutex.Lock()
	defer connsMutex.Unlock()

	for conn, lastActivity := range activeConns {
		// Check if the connection is inactive for too long
		if time.Since(lastActivity) > connCheckPeriod {
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("Connection to %s is broken, removing\n", conn.RemoteAddr())
				conn.Close()
				delete(activeConns, conn)
			} else {
				activeConns[conn] = time.Now() // Update last activity time
			}
		}
	}
}
