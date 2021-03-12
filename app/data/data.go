package data

import "fmt"

// CurrentGasPrice - When `gasz` service is queried, it'll send
// response of this form back
type CurrentGasPrice struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
}

func (c *CurrentGasPrice) String() string {

	return fmt.Sprintf("fastest : %.2f Gwei | fast : %.2f Gwei | average : %.2f Gwei | safeLow : %.2f Gwei", c.Fastest, c.Fast, c.Average, c.SafeLow)

}

// Sendable - Send response to Telegram, when asked for latest gas price
func (c *CurrentGasPrice) Sendable() string {

	return fmt.Sprintf("Fastest : %.2f Gwei 🚀\nFast : %.2f Gwei\nAverage : %.2f Gwei\nSafeLow : %.2f Gwei 🐢", c.Fastest, c.Fast, c.Average, c.SafeLow)

}

// Payload - Subscribe to latest gas price feed of `gasz`, over
// websocket transport, by sending payload of this form
type Payload struct {
	Type      string  `json:"type"`
	Field     string  `json:"field"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
}

// Response - Subscription/ unsubscription confirmation messages
// to be received in this form
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Resources - ...
type Resources struct {
	Latest *CurrentGasPrice
}
