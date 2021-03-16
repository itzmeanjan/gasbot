package data

// Response - Subscription/ unsubscription confirmation messages
// to be received in this form
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
