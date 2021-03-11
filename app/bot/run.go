package bot

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/gasz"
	"gopkg.in/tucnak/telebot.v2"
)

// GetHTTPClient - HTTP client to be used for talking to
// Telegram servers
//
// When new message will be received from telegram over
// Webhook, it'll be responded back using this custom HTTP client
func GetHTTPClient() *http.Client {

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

	return client

}

// Run - Starts Telegram bot
func Run() error {

	endpoint := &telebot.WebhookEndpoint{
		PublicURL: config.GetURL(),
	}

	webhook := &telebot.Webhook{
		Listen:   fmt.Sprintf(":%d", config.GetPort()),
		Endpoint: endpoint,
	}

	token := config.GetToken()
	if token == "" {
		return errors.New("Bad Token")
	}

	settings := telebot.Settings{
		Token:  token,
		Poller: webhook,
		Client: GetHTTPClient(),
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return err
	}

	bot.Handle("/start", func(m *telebot.Message) {

		log.Printf("📩 [ /start ] : From %s\n", m.Sender.Username)

		bot.Send(m.Sender, "Ethereum Gas Price Notifier @ https://gasz.in\n\nBuilt & maintained by Anjan Roy<anjanroy@yandex.com>\n\nFind more about me @ https://itzmeanjan.in")

	})

	bot.Handle("/latest", func(m *telebot.Message) {

		log.Printf("📩 [ /latest ] : From %s\n", m.Sender.Username)

		gasPrice, err := gasz.CurrentGasPrice()
		if err != nil {

			log.Printf("❗️ Failed to get latest gas price : %s\n", err.Error())
			return

		}

		bot.Send(m.Sender, gasPrice.Sendable())

	})

	log.Printf("✅ Starting %s\n", bot.Me.Username)

	// This is a blocking call
	bot.Start()

	return nil

}
