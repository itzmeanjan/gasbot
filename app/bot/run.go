package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/data"
	"github.com/itzmeanjan/gasbot/app/gasz"
	"gopkg.in/tucnak/telebot.v2"
)

// Run - Starts Telegram bot & keeps serving
// requests
func Run(resources *data.Resources) error {

	endpoint := &telebot.WebhookEndpoint{
		PublicURL: config.GetURL(),
	}

	webhook := &telebot.Webhook{
		Listen:   fmt.Sprintf(":%d", config.GetPort()),
		Endpoint: endpoint,
	}

	token := config.GetToken()
	if token == "" {
		return errors.New("bad Token")
	}

	settings := telebot.Settings{
		Token:  token,
		Poller: webhook,
		Client: gasz.GetHTTPClient(),
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return err
	}

	resources.Bot = bot

	bot.Handle("/start", func(m *telebot.Message) {

		log.Printf("📩 [ /start ] : From @%s\n", m.Sender.Username)

		bot.Send(m.Sender, "Ethereum Gas Price Notifier @ https://gasz.in\n\n/latest : Get latest gas price recommendation\n/subscribe : Get notified when gas price reaches threshold\n\nBuilt & maintained by Anjan Roy<anjanroy@yandex.com>\n\nFind more about me @ https://itzmeanjan.in")

	})

	bot.Handle("/latest", func(m *telebot.Message) {

		log.Printf("📩 [ /latest ] : From @%s\n", m.Sender.Username)

		// Send latest gas price feed, which was received
		// from `gasz` subscription
		bot.Send(m.Sender, resources.Latest.Sendable())

	})

	bot.Handle("/subscribe", func(m *telebot.Message) {

		log.Printf("📩 [ /subscribe ] : From @%s\n", m.Sender.Username)

		txType, operator, threshold, err := parseSubscriptionPayload(m.Payload)
		if err != nil {

			bot.Send(m.Sender, fmt.Sprintf("❗️ I received : %s", err.Error()))
			return

		}

		if err := resources.Subscribe(m.Sender, txType, operator, threshold); err != nil {

			bot.Send(m.Sender, fmt.Sprintf("❗️ I received : %s", err.Error()))
			return

		}

		bot.Send(m.Sender, "🎉  Subscription confirmed")

	})

	bot.Handle("/unsubscribe", func(m *telebot.Message) {

		log.Printf("📩 [ /unsubscribe ] : From @%s\n", m.Sender.Username)

		if err := resources.Unsubscribe(m.Sender); err != nil {

			bot.Send(m.Sender, fmt.Sprintf("❗️ I received : %s", err.Error()))
			return

		}

		bot.Send(m.Sender, "🙂  Unsubscription confirmed")

	})

	// These are only commands supported by `gasbot`
	if err := bot.SetCommands([]telebot.Command{
		{Text: "latest", Description: "Ask for latest gas price recommendation"},
		{Text: "subscribe", Description: "Get notified when gas price reaches threshold"},
		{Text: "unsubscribe", Description: "Unsubscribe from subscribed updates"},
	}); err != nil {

		return err

	}

	log.Printf("✅ Starting @%s\n", bot.Me.Username)

	// This is a blocking call
	bot.Start()

	return nil

}
