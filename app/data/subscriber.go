package data

import "gopkg.in/tucnak/telebot.v2"

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

// SendNotification - Sends notification to user, letting them know recommended
// gas price has reached certain threshold, of their interest
func (s *Subscriber) SendNotification(handle *telebot.Bot, gasPrice *CurrentGasPrice) error {

	_, err := handle.Send(s.User, s.Criteria.PrepareNotification(gasPrice))
	return err

}
