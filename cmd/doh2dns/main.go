package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Merith-TK/utils/debug"
	"github.com/miekg/dns"
)

// DOC: https://developers.google.com/speed/public-dns/docs/doh/json

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

var dnsServer string   // Global variable to hold the DNS server address
var hostAddress string // Global variable to hold the host address
func main() {
	// Define a command-line flag for the DNS server
	flag.StringVar(&hostAddress, "host", "", "The hosted address for this DNS server")
	flag.StringVar(&dnsServer, "dns", "1.1.1.1", "DNS server address (e.g. 1.1.1.1:53). Leave empty to use the system resolver.")
	flag.Parse()

	if debug.GetDebug() {
		startDoHServer()
		return
	}
	hostAddress = strings.TrimPrefix(hostAddress, "http://")
	hostAddress = strings.Join([]string{"https://", hostAddress}, "")

	if !strings.Contains(dnsServer, ":") {
		dnsServer = strings.Join([]string{dnsServer, "53"}, ":")
	}

	log.Println("Host address:", hostAddress)
	randStr, err := generateRandomString(16)
	if err != nil {
		log.Fatal("Failed to generate random string:", err)
	}

	srv := &http.Server{Addr: ":8080"}
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(randStr))
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// start server and wait to kill it
	go func() {
		log.Println("Starting Verification server on :8080")
		log.Println(srv.ListenAndServe())
	}()
	// request the random string from the server
	resp, err := http.Get("https://" + hostAddress + ":8080")
	if err != nil {
		srv.Shutdown(ctx)
		log.Fatal("Failed to get random string:", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		srv.Shutdown(ctx)
		log.Fatal("Failed to read response body:", err)
	}
	err = resp.Body.Close()
	if err != nil {
		srv.Shutdown(ctx)
		log.Fatal("Failed to close response body:", err)
	}

	if string(body) != randStr {
		log.Println("Random String Mismatch!")
		log.Println("Expected:", randStr)
		log.Println("Received:", string(body))
	} else {
		log.Println("Random String Match!")
		srv.Shutdown(ctx)

		startDoHServer()
	}

}

func startDoHServer() {
	fmt.Println("Starting DoH to DNS server...")
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
	qtype = strings.ToUpper(qtype)
	var qtypeInt16 uint16

	// If no type provided, set default to A (IPv4 address)
	if qtype == "" {
		qtypeInt16 = dns.TypeA
	} else {
		// Try converting string to int
		if parsedInt, err := strconv.Atoi(qtype); err == nil {
			qtypeInt16 = uint16(parsedInt) // Successfully parsed an integer
		} else if val, found := dns.StringToType[qtype]; found {
			// Check if the string matches a known DNS type
			qtypeInt16 = val
		} else {
			// Default fallback if the input is unrecognized
			qtypeInt16 = dns.TypeNone // or handle this as an error, if needed
		}
	}

	debug.Print("Received DNS query from", r.RemoteAddr, "for", domain, "type", qtype, "\n\t Parsed Type:", qtypeInt16)
	if domain == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}
	if !strings.HasSuffix(domain, ".") {
		domain += "." // Ensure domain ends with a dot
	}

	// make a DNS request to get the ip of example.com, and get as much information as possible
	dnsClient := new(dns.Client)

	dnsClient.Dialer = &net.Dialer{
		Timeout: 200 * time.Millisecond,
	}
	dnsQuery := new(dns.Msg)
	dnsQuery.SetQuestion(domain, qtypeInt16)

	in, _, err := dnsClient.Exchange(dnsQuery, "1.1.1.1:53")
	if err != nil {
		log.Fatalf("Failed to exchange: %v", err)
	}

	// Convert DNS response to DoH response
	dohResp := DoHResponse{
		Status:   in.Rcode,
		Tc:       in.Truncated,
		Rd:       in.RecursionDesired,
		Ra:       in.RecursionAvailable,
		Ad:       in.AuthenticatedData,
		Cd:       in.CheckingDisabled,
		Question: []DNSQuestion{},
		Answer:   []DNSAnswer{},
	}

	for _, q := range in.Question {
		q.Name = strings.TrimSuffix(q.Name, ".")

		dnsQuestion := DNSQuestion{
			Name: q.Name,
			Type: int(q.Qtype),
		}
		dohResp.Question = append(dohResp.Question, dnsQuestion)
	}

	for _, ans := range in.Answer {
		dnsAnswer := DNSAnswer{
			Name: strings.TrimSuffix(ans.Header().Name, "."),
			Type: int(ans.Header().Rrtype),
			TTL:  int(ans.Header().Ttl),
			Data: ans.String(),
		}
		dohResp.Answer = append(dohResp.Answer, dnsAnswer)
	}

	debug.Print("DoH Response:", dohResp)

	// Respond with the DNS response as JSON
	w.Header().Set("Content-Type", "application/dns-json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dohResp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
