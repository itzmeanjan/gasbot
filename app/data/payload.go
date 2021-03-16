package data

import (
	"fmt"
	"strings"
)

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
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f Gwei", strings.ToTitle(p.Field), gasPrice.Fastest)
	case "fast":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f Gwei", strings.ToTitle(p.Field), gasPrice.Fast)
	case "average":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f Gwei", strings.ToTitle(p.Field), gasPrice.Average)
	case "safeLow":
		notification = fmt.Sprintf("Hey 👋, Gas Price for `%s` tx has reached : %.2f Gwei", strings.ToTitle(p.Field), gasPrice.SafeLow)
	default:
		// @note Not doing anything here, because result is negative

	}

	return notification

}
