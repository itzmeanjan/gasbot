package gasz

import (
	"net"
	"net/http"
	"time"
)

// GetHTTPClient - HTTP client to be used for talking to
// Telegram servers & `gasz : Ethereum Gas Price Notifier` server
//
// When new message will be received from telegram over
// Webhook, it'll be responded back using this custom HTTP client
//
// When it's required to talk to `gasz`, for getting latest gas price feed
// this client to be used
func GetHTTPClient() *http.Client {

	// Dialing will spend at max 1 second
	dialer := &net.Dialer{
		Timeout: time.Duration(5) * time.Second,
	}

	// TLS handshake must happen with in 1 second
	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: time.Duration(5) * time.Second,
	}

	// Whole process must happen with in 3 seconds
	//
	// Otherwise it'll time out
	client := &http.Client{
		Timeout:   time.Duration(10) * time.Second,
		Transport: transport,
	}

	return client

}
