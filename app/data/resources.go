package data

import (
	"errors"
	"log"
	"sync"

	"gopkg.in/tucnak/telebot.v2"
)

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

// Subscribe - Subscribes to gas price of certain tx category, with condition to be
// evaluated on their sake when new gas price is seen, to decide whether we need to
// notify subscriber or not
func (r *Resources) Subscribe(user *telebot.User, txType string, operator string, threshold float64) error {

	// Valid tx types to which clients can subscribe to
	if !(txType == "fastest" || txType == "fast" || txType == "average" || txType == "safeLow") {
		return errors.New("txType ∈ {fastest, fast, average, safeLow}")
	}

	// Valid relational operators
	if !(operator == "<" || operator == ">" || operator == "<=" || operator == ">=" || operator == "==") {
		return errors.New("operator ∈ {<, >, <=, >=, ==}")
	}

	// Threshold i.e. gas price value against which comparison to be performed
	// needs to be >= 1.0 Gwei
	if !(threshold >= 1.0) {
		return errors.New("threshold >= 1.0 Gwei")
	}

	r.Lock.Lock()
	defer r.Lock.Unlock()

	if sub, ok := r.Subscriptions[user.Username]; ok {

		// Because user has already subscribed to
		// some topic, it'll simply update previous choice
		// with latest one

		sub.Criteria.Field = txType
		sub.Criteria.Operator = operator
		sub.Criteria.Threshold = threshold

		return nil

	}

	// User is subscribing for first time, putting
	// their entry in in-memory directory
	r.Subscriptions[user.Username] = &Subscriber{
		User: user,
		Criteria: &Payload{
			Field:     txType,
			Operator:  operator,
			Threshold: threshold,
		},
	}

	return nil

}

// Unsubscribe - Unsubscribes this user from his/ her subscription, if any exists
func (r *Resources) Unsubscribe(user *telebot.User) error {

	r.Lock.Lock()
	defer r.Lock.Unlock()

	if _, ok := r.Subscriptions[user.Username]; !ok {

		return errors.New("not subscribed yet")

	}

	// Removing entry
	delete(r.Subscriptions, user.Username)

	return nil

}
