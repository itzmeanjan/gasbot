package gasz

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/data"
)

// CurrentGasPrice - Queries `gasz` service for current recommended
// gas price & returns repsonse back
//
// @note Caller is supposed to handle errors responsibily
func CurrentGasPrice() (*data.CurrentGasPrice, error) {

	// Dialing will spend at max 1 second
	dialer := &net.Dialer{
		Timeout: time.Duration(1) * time.Second,
	}

	// TLS handshake must happen with in 1 second
	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: time.Duration(1) * time.Second,
	}

	// Whole process must happen with in 3 seconds
	//
	// Otherwise it'll time out
	client := &http.Client{
		Timeout:   time.Duration(3) * time.Second,
		Transport: transport,
	}

	// Making HTTP GET request to remote for getting
	// latest gas price recommendation
	resp, err := client.Get(config.GetGaszQueryURL())
	if err != nil {
		return nil, err
	}

	// Scheduling closure of response body
	//
	// @note To be invoked when returning from this execution scope
	defer resp.Body.Close()

	// Reading whole body of response in byte array
	_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gasPrice data.CurrentGasPrice

	// Attempting to deserialize byte array into structured format
	// so that it can be easily interacted with
	if err := json.Unmarshal(_data, &gasPrice); err != nil {
		return nil, err
	}

	return &gasPrice, nil

}
