package data

import (
	"fmt"
	"sync"

	"gopkg.in/tucnak/telebot.v2"
)

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

// SatisfiedBy - Performs a check whether current gas price
// recommendation is satisfying the criteria some user has provided
//
// If yes, we can take next steps i.e. sending them notification
func (p *Payload) SatisfiedBy(gasPrice *CurrentGasPrice) bool {

	checkThreshold := func(price float64) bool {

		var yes bool

		switch p.Operator {

		case "<":
			yes = price < p.Threshold
		case ">":
			yes = price > p.Threshold
		case "<=":
			yes = price <= p.Threshold
		case ">=":
			yes = price >= p.Threshold
		case "==":
			yes = price == p.Threshold
		default:
			// @note No need to do anything
			// Check is going to be return negative result
		}

		return yes

	}

	var yes bool

	switch p.Field {

	case "fastest":
		yes = checkThreshold(gasPrice.Fastest)
	case "fast":
		yes = checkThreshold(gasPrice.Fast)
	case "average":
		yes = checkThreshold(gasPrice.Average)
	case "safeLow":
		yes = checkThreshold(gasPrice.SafeLow)
	default:
		// @note Not doing anything here, because result is negative

	}

	return yes

}

// PrepareNotification - Prepares notification text, that will be sent to user
// once we've decided whether this user is eligible to receive this notification
// or not
//
// @note Whether eligible or not, that's decided based upon what that user has set
// in their criteria for receiving notification
func (p *Payload) PrepareNotification(gasPrice *CurrentGasPrice) string {

	var notification string

	switch p.Field {

	case "fastest":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %f", p.Field, gasPrice.Fastest)
	case "fast":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %f", p.Field, gasPrice.Fast)
	case "average":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %f", p.Field, gasPrice.Average)
	case "safeLow":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %f", p.Field, gasPrice.SafeLow)
	default:
		// @note Not doing anything here, because result is negative

	}

	return notification

}

// Response - Subscription/ unsubscription confirmation messages
// to be received in this form
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Resources - These are the resources which are supposed to be accessed
// by multiple go routines ( can be simultaneously )
type Resources struct {
	Latest        *CurrentGasPrice
	Subscriptions map[string]*Subscriber
	Lock          *sync.RWMutex
}

// Subscriber - This is one Telegram User, who has interacted with `gasbot`
//
// We're going to be keep track of their subscription interest in this section
type Subscriber struct {
	User     *telebot.User
	Criteria *Payload
}

// CanSendNotification - Checks whether recent gas price update we received
// can that be sent to subscribed user
//
// It'll be sent, if & only if it satisfies criteria set by user
func (s *Subscriber) CanSendNotification(gasPrice *CurrentGasPrice) bool {
	return s.Criteria.SatisfiedBy(gasPrice)
}
