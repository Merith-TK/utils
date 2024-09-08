package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Merith-TK/utils/debug"
	"github.com/gen2brain/beeep"
	"github.com/gorilla/websocket"
)

var config Config

type Config struct {
	Notify struct {
		Email struct {
			Enabled bool   `json:"enabled,omitempty"`
			Server  string `json:"server,omitempty"`
			User    string `json:"user,omitempty"`
			Pass    string `json:"pass,omitempty"`
			From    string `json:"from,omitempty"`
			To      string `json:"to,omitempty"`
		} `json:"email,omitempty"`
		Sms struct {
			Enabled bool   `json:"enabled,omitempty"`
			Phone   string `json:"phone,omitempty"`
			Header  string `json:"header,omitempty"`
		} `json:"sms,omitempty"`
		Beep struct {
			Enabled  bool    `json:"enabled,omitempty"`
			Freq     float64 `json:"freq,omitempty"`
			Duration int     `json:"duration,omitempty"`
		} `json:"beep,omitempty"`
		Notification struct {
			Enabled bool `json:"enabled,omitempty"`
		} `json:"notification,omitempty"`
	} `json:"notify,omitempty"`
	Servers []string `json:"servers,omitempty"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "servers.json", "path to the config file")
}

func readConfigFromFile(filename string) error {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a Config struct
	var servers Config
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return err
	}
	emailconfig := servers.Notify.Email
	smsconfig := servers.Notify.Sms
	beepconfig := servers.Notify.Beep
	noticonfig := servers.Notify.Notification

	// populate the default values for
	if beepconfig.Freq == 0 {
		beepconfig.Freq = beeep.DefaultFreq
	}
	if beepconfig.Duration == 0 {
		beepconfig.Duration = beeep.DefaultDuration
	}

	if emailconfig.From == "" {
		emailconfig.From = emailconfig.User
	}

	servers.Notify.Email = emailconfig
	servers.Notify.Sms = smsconfig
	servers.Notify.Beep = beepconfig
	servers.Notify.Notification = noticonfig
	config = servers
	return nil
}

func main() {
	flag.Parse()
	debug.Print("Downtime client started")

	// Read the servers from the JSON file
	err := readConfigFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Config:")
	fmt.Println("\tNotify:")
	fmt.Println("\t\tEmail: ", config.Notify.Email.Enabled)
	if config.Notify.Email.Enabled {
		fmt.Println("\t\t\tServer:\t", config.Notify.Email.Server)
		fmt.Println("\t\t\tUser:\t", config.Notify.Email.User)
		fmt.Println("\t\t\tFrom:\t", config.Notify.Email.From)
		fmt.Println("\t\t\tTo:\t", config.Notify.Email.To)
	}
	fmt.Println("\t\tSms:   ", config.Notify.Sms.Enabled)
	if config.Notify.Sms.Enabled {
		fmt.Println("\t\t\tPhone:\t", config.Notify.Sms.Phone)
		fmt.Println("\t\t\tHeader:\t", config.Notify.Sms.Header)
	}
	fmt.Println("\t\tBeep:  ", config.Notify.Beep.Enabled)
	if config.Notify.Beep.Enabled {
		fmt.Println("\t\t\tFreq:\t ", strconv.FormatFloat(config.Notify.Beep.Freq, 'f', -1, 64)+"Hz")
		fmt.Println("\t\t\tDuration:", strconv.FormatInt(int64(config.Notify.Beep.Duration), 10)+"ms")
	}
	fmt.Println("\t\tNotify:", config.Notify.Notification.Enabled)
	fmt.Println("\tServers:")
	for _, server := range config.Servers {
		fmt.Println("\t\t", server)
	}

	// Start connection attempts for each server
	for _, server := range config.Servers {
		go maintainConnection(server)
	}

	// Keep the main goroutine alive
	select {}
}
func maintainConnection(server string) {
	var mainConn, heartbeatConn *websocket.Conn
	var err error

	for {
		// Check the main connection
		if mainConn == nil || !isConnectionAlive(mainConn) {
			// Attempt to connect to the main WebSocket server ("/")
			mainConn, err = connectToServer(server)
			if err != nil {
				debug.Print(fmt.Sprintf("Failed to connect to main server %s: %v", server, err))
				time.Sleep(10 * time.Second) // Retry after 10 seconds
				continue
			}
			debug.Print(fmt.Sprintf("Connected to main server %s", server))

			// Start listening for messages from the main server
			go handleServerMessages(server, mainConn)
		}

		// Check the heartbeat connection
		if heartbeatConn == nil || !isConnectionAlive(heartbeatConn) {
			// Attempt to connect to the heartbeat WebSocket server ("/heartbeat")
			debug.Print(fmt.Sprintf("Connecting to heartbeat server %s", server+"/heartbeat"))
			heartbeatConn, err = connectToServer(server + "/heartbeat")
			if err != nil {
				debug.Print(fmt.Sprintf("Failed to connect to heartbeat server %s: %v", server, err))
				time.Sleep(10 * time.Second) // Retry after 10 seconds
				continue
			}
			debug.Print(fmt.Sprintf("Connected to heartbeat server %s", server+"/heartbeat"))

			// Start listening for heartbeat messages
			go handleHeartbeatMessages(server, heartbeatConn)
		}

		// Sleep for a while before checking the connection status again
		time.Sleep(10 * time.Second)
	}
}

func handleServerMessages(url string, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			debug.Print(fmt.Sprintf("Connection to server %s ended: %v", url, err))
			return
		}

		// Handle the message
		debug.Print(fmt.Sprintf("Received message from server %s: %s", url, string(message)))
	}
}

// Handle messages from the heartbeat WebSocket server
func handleHeartbeatMessages(url string, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			debug.Print(fmt.Sprintf("Heartbeat connection to server %s ended: %v", url, err))
			handleDisconnection(url)
			return
		}

		// Handle the heartbeat message
		debug.Print(fmt.Sprintf("Received heartbeat from server %s: %s", url, string(message)))
	}
}

func isConnectionAlive(conn *websocket.Conn) bool {
	debug.Print("Checking connection status of:", conn.LocalAddr())

	// Send a ping message to the server
	err := conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		return false
	}

	// Wait for the pong message
	_, _, err = conn.ReadMessage()
	return err == nil
}

func connectToServer(urlString string) (*websocket.Conn, error) {
	debug.Print("Connecting to server:", urlString)
	// Create a WebSocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout:  10 * time.Second,
		EnableCompression: false,
	}
	if strings.HasPrefix(urlString, "http://") || strings.HasPrefix(urlString, "https://") {
		urlString = strings.Replace(urlString, "http", "ws", 1)
	}
	if !strings.HasPrefix(urlString, "ws://") && !strings.HasPrefix(urlString, "wss://") {
		urlString = "ws://" + urlString
	}

	newURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	// Connect to the WebSocket server
	debug.Print("Connecting to parsed:", newURL)
	conn, _, err := dialer.Dial(newURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func handleDisconnection(url string) {
	debug.SetTitle("Downtime Alert")
	debug.Print(fmt.Sprintf("Connection to server %s ended", url))
	if config.Notify.Beep.Enabled {
		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}
		debug.Print("Beeped")
	}
	if config.Notify.Notification.Enabled {
		err := beeep.Notify("Downtime Alert", fmt.Sprintf("Connection to server %s ended", url), "")
		if err != nil {
			panic(err)
		}
		debug.Print("Notified")
	}

	// TODO: Implement Email Notification
	// TODO: Implement SMS Notification

	debug.ResetTitle()
}
