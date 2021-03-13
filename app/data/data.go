package data

import (
	"fmt"
	"log"
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
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f", p.Field, gasPrice.Fastest)
	case "fast":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f", p.Field, gasPrice.Fast)
	case "average":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f", p.Field, gasPrice.Average)
	case "safeLow":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f", p.Field, gasPrice.SafeLow)
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
	Bot           *telebot.Bot
	Latest        *CurrentGasPrice
	Subscriptions map[string]*Subscriber
	Lock          *sync.RWMutex
}

// Notify - As soon as new gas price update is received, it'll
// iteratively go through each of subscribers & check whether they're
// eligible to receive this notification or not
//
// If yes, they'll be attempted to be notified
func (r *Resources) Notify() {

	r.Lock.RLock()
	defer r.Lock.RUnlock()

	for k, v := range r.Subscriptions {

		// These are currently ongoing subscription phases
		// no need to touch them
		//
		// They're not yet confirmed, hopefully will be in sometime future
		if v.InProgress {
			continue
		}

		if !v.CanSendNotification(r.Latest) {
			continue
		}

		if err := v.SendNotification(r.Bot, r.Latest); err != nil {
			log.Printf("❌ Failed to notify @%s that gas price has reached their desired threshold : %s\n", k, err.Error())
			continue
		}

		log.Printf("🔔 Notified @%s as Gas Price reached their desired threshold\n", k)

	}

}

// InitSubscription - Subscription is a multi step operation, this is very beginning of it
// More info to be collected from user
func (r *Resources) InitSubscription(user *telebot.User, txType string) bool {

	// Valid tx types to which clients can subscribe to
	if !(txType == "fastest" || txType == "fast" || txType == "average" || txType == "safeLow") {
		return false
	}

	r.Lock.RLock()
	defer r.Lock.RUnlock()

	_, ok := r.Subscriptions[user.Username]
	if ok {
		return false
	}

	// -- Acquire lock
	r.Lock.Lock()

	r.Subscriptions[user.Username] = &Subscriber{
		InProgress: true,
		User:       user,
		Criteria: &Payload{
			Field: txType,
		},
	}

	r.Lock.Unlock()
	// -- Released lock

	return true

}

// SetSubscriptionOperator - This is second step of subscription to gas price feed
// where user selects what's relational operator to be used when checking whether some gas price
// update needs to be sent to them or not
func (r *Resources) SetSubscriptionOperator(user *telebot.User, operator string) bool {

	// Valid relational operators
	if !(operator == "<" || operator == ">" || operator == "<=" || operator == ">=" || operator == "==") {
		return false
	}

	r.Lock.RLock()
	defer r.Lock.RUnlock()

	sub, ok := r.Subscriptions[user.Username]
	if !ok {
		return false
	}

	if !(sub.InProgress && sub.Criteria.Field != "" && sub.Criteria.Operator == "") {
		return false
	}

	sub.Criteria.Operator = operator
	return true

}

// SetSubscriptionThreshold - When gas price of certain category reaches this value ( in Gwei ), user
// wants us to notify them
func (r *Resources) SetSubscriptionThreshold(user *telebot.User, threshold float64) bool {

	// Threshold i.e. gas price value against which comparison to be performed
	// needs to be >= 1.0 Gwei
	if !(threshold >= 1.0) {
		return false
	}

	r.Lock.RLock()
	defer r.Lock.RUnlock()

	sub, ok := r.Subscriptions[user.Username]
	if !ok {
		return false
	}

	if !(sub.InProgress && sub.Criteria.Field != "" && sub.Criteria.Operator != "" && sub.Criteria.Threshold == 0) {
		return false
	}

	sub.Criteria.Threshold = threshold
	return true

}

// ConfirmSubscription - This subscription is confirmed, now it can
// be attempted to be processed when latest gas price feed is received
func (r *Resources) ConfirmSubscription(user *telebot.User) bool {

	r.Lock.Lock()
	defer r.Lock.Unlock()

	sub, ok := r.Subscriptions[user.Username]
	if !ok {
		return false
	}

	sub.InProgress = false
	return true

}

// Subscriber - This is one Telegram User, who has interacted with `gasbot`
//
// We're going to be keep track of their subscription interest in this section
type Subscriber struct {
	InProgress bool
	User       *telebot.User
	Criteria   *Payload
}

// CanSendNotification - Checks whether recent gas price update we received
// can that be sent to subscribed user
//
// It'll be sent, if & only if it satisfies criteria set by user
func (s *Subscriber) CanSendNotification(gasPrice *CurrentGasPrice) bool {
	return s.Criteria.SatisfiedBy(gasPrice)
}

// SendNotification - Sends notification to user, letting them know recommended
// gas price has reached certain threshold, of their interest
func (s *Subscriber) SendNotification(handle *telebot.Bot, gasPrice *CurrentGasPrice) error {

	_, err := handle.Send(s.User, s.Criteria.PrepareNotification(gasPrice))
	return err

}
