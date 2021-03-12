package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/gasz"
	"gopkg.in/tucnak/telebot.v2"
)

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

	bot.Handle("/start", func(m *telebot.Message) {

		log.Printf("📩 [ /start ] : From %s\n", m.Sender.Username)

		bot.Send(m.Sender, "Ethereum Gas Price Notifier @ https://gasz.in\n\n/latest : Latest Ethereum Gas Price recommendation\n\nBuilt & maintained by Anjan Roy<anjanroy@yandex.com>\n\nFind more about me @ https://itzmeanjan.in")

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

	log.Printf("✅ Starting `%s`\n", bot.Me.Username)

	// This is a blocking call
	bot.Start()

	return nil

}
