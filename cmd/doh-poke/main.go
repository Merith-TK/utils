// Package main implements a DNS-over-HTTPS (DoH) client that queries DNS servers
// using the DoH protocol with JSON response parsing.
//
// The doh-poke utility performs DNS queries over HTTPS using the RFC 8484 standard,
// providing a secure alternative to traditional DNS queries. It supports automatic
// HTTPS scheme detection and configurable DoH server endpoints.
//
// Features:
//   - DoH query support with JSON response parsing
//   - Automatic HTTPS scheme detection and URL reconstruction
//   - Configurable DoH server endpoints
//   - A record resolution with multiple answer support
//   - Debug output integration
//
// Usage:
//   doh-poke [doh-server] [domain]
//
// Parameters:
//   doh-server  DoH server URL (scheme optional, defaults to HTTPS)
//   domain      Domain name to resolve
//
// Examples:
//   doh-poke cloudflare-dns.com example.com
//   doh-poke https://1.1.1.1/dns-query google.com
//   doh-poke 8.8.8.8 github.com
//
// The tool automatically appends /dns-query to the path if not specified
// and handles URL parsing with proper scheme detection.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Merith-TK/utils/pkg/debug"
)

type DoHResponse struct {
	Status   int  `json:"Status"`
	Tc       bool `json:"TC"`
	Rd       bool `json:"RD"`
	Ra       bool `json:"RA"`
	Ad       bool `json:"AD"`
	Cd       bool `json:"CD"`
	Question []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"Question"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

func main() {
	flag.Parse()
	debug.Print("Starting DoH poke...")
	// Parse command line arguments
	dohServer := flag.Arg(0)
	domain := flag.Arg(1)

	debug.Print("dohServer:", dohServer)
	debug.Print("domain:", domain)
	if !strings.HasPrefix(dohServer, "http://") && !strings.HasPrefix(dohServer, "https://") {
		dohServer = "https://" + dohServer
		debug.Print("apply default scheme:", dohServer)
	}

	urlp, err := url.Parse(dohServer)
	if err != nil {
		log.Fatalln("Error parsing URL:", err)
	}
	debug.Print("urlp:", urlp)

	urlRebuild := []string{}
	if urlp.Scheme == "" {
		urlRebuild = append(urlRebuild, "https://")
	} else {
		urlRebuild = append(urlRebuild, urlp.Scheme+"://")
	}
	urlRebuild = append(urlRebuild, urlp.Host)
	if urlp.Path == "" || urlp.Path == "/" || urlp.Path == flag.Arg(0) {
		urlRebuild = append(urlRebuild, "/dns-query")
	} else {
		urlRebuild = append(urlRebuild, urlp.Path)
	}
	st := strings.Join(urlRebuild, "")
	debug.Print("st:", st)

	// Create the query URL
	queryURL := fmt.Sprintf("%s?name=%s&type=A", st, domain)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Set("Accept", "application/dns-json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform DNS query: %v", err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed to perform DNS query: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the response body into a DoHResponse struct
	var dohResponse DoHResponse
	err = json.Unmarshal(body, &dohResponse)
	if err != nil {
		log.Println("Failed to parse JSON response:", err)
		log.Println("Response body:", string(body))
	}

	debug.Print("Value:", string(body))

	// Print the answer
	for _, answer := range dohResponse.Answer {
		fmt.Println("Answer:", answer.Data)
	}
}
