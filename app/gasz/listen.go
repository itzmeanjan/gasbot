package gasz

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/data"
)

// SubscribeToLatest - Subscribe to latest gas price feed of `gasz`
func SubscribeToLatest(ctx context.Context, comm chan<- struct{}, resources *data.Resources) {

	netDialer := &net.Dialer{
		Timeout: time.Duration(5) * time.Second,
	}

	dialer := &websocket.Dialer{
		NetDialContext:   netDialer.DialContext,
		HandshakeTimeout: time.Duration(5) * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}

	conn, _, err := dialer.DialContext(ctx, config.GetGaszSubscribeURL(), nil)
	if errors.Is(err, websocket.ErrBadHandshake) {

		log.Printf("❗️ Websocket handshake failed : %s\n", err.Error())
		close(comm)
		return

	}

	if conn == nil {

		log.Printf("❗️ Bad Websocket connection\n")
		close(comm)
		return

	}

	// Scheduling clean up from now
	defer func() {
		close(comm)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))

	// Handle ping messages received from server &
	// increment read deadline
	conn.SetPingHandler(func(appData string) error {

		conn.SetReadDeadline(time.Now().Add(time.Duration(15) * time.Second))
		return conn.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(time.Duration(1)*time.Second))

	})

	subPayload := &data.Payload{
		Type:      "subscription",
		Field:     "*",
		Threshold: 1,
		Operator:  "*",
	}

	// Send subscription request, for receiving
	// latest gas price feed
	if err := conn.WriteJSON(subPayload); err != nil {

		log.Printf("❗️ Failed to send subscription request to Gas Price feed : %s\n", err.Error())
		return

	}

	var confirmation data.Response

	// Waiting for subscription confirmation message from `gasz`
	if err := conn.ReadJSON(&confirmation); err != nil {

		log.Printf("❗️ Failed to receive subscription confirmation from Gas Price feed : %s\n", err.Error())
		return

	}

	// Subscription attempt failed
	if confirmation.Code != 1 {

		log.Printf("❗️ Gas Price feed subscription denied\n")
		return

	}

	log.Printf("✅ Subscribed to latest Gas Price feed\n")

	// Iteratively keep reading latest gas price updates
	// & put fresh data in shared resource
	for {

		var gasPrice data.CurrentGasPrice

		if err := conn.ReadJSON(&gasPrice); err != nil {

			log.Printf("❗️ Failed to receive Gas Price subscription : %s\n", err.Error())
			break

		}

		// Updating shared resource with latest gas price
		resources.Latest = &gasPrice

		log.Printf("%s\n", resources.Latest.String())

	}

}
