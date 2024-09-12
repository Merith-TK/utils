package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
)

// Structs for DoH response
type DNSQuestion struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type DNSAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

type DoHResponse struct {
	Status   int           `json:"Status"`
	Tc       bool          `json:"TC"`
	Rd       bool          `json:"RD"`
	Ra       bool          `json:"RA"`
	Ad       bool          `json:"AD"`
	Cd       bool          `json:"CD"`
	Question []DNSQuestion `json:"Question"`
	Answer   []DNSAnswer   `json:"Answer"`
}

func main() {
	http.HandleFunc("/dns-query", handleDNSRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDNSRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/dns-query" {
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	// Parse query parameters
	q := r.URL.Query()
	domain := q.Get("name")
	qtype := q.Get("type")

	if domain == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}

	if qtype == "" {
		qtype = "A" // Default to 'A' record if not specified
	}

	queryType := getType(qtype)
	if queryType == 0 {
		http.Error(w, "Unsupported query type", http.StatusBadRequest)
		return
	}

	// Perform DNS lookup
	ips, err := net.LookupHost(domain)
	if err != nil {
		http.Error(w, "DNS lookup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Construct DoH response
	dohResp := DoHResponse{
		Status: http.StatusOK,
		Tc:     false,
		Rd:     true,
		Ra:     true,
		Ad:     false,
		Cd:     false,
		Question: []DNSQuestion{
			{Name: domain, Type: queryType},
		},
	}

	// Add answers to DoH response
	for _, ip := range ips {
		answer := DNSAnswer{
			Name: domain,
			Type: queryType,
			TTL:  300, // Static TTL value; you can make this dynamic if needed
			Data: ip,
		}
		dohResp.Answer = append(dohResp.Answer, answer)
	}

	// Respond with the DNS response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dohResp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// getType translates a query type (e.g., "A", "AAAA") to its DNS code
func getType(qtype string) int {
	switch qtype {
	case "A":
		return 1
	case "AAAA":
		return 28
	case "CNAME":
		return 5
	case "MX":
		return 15
	case "NS":
		return 2
	case "PTR":
		return 12
	case "SOA":
		return 6
	case "TXT":
		return 16
	default:
		return 0 // Unsupported type
	}
}
