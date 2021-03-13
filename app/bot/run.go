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

		bot.Send(m.Sender, "Ethereum Gas Price Notifier @ https://gasz.in\n\n/latest : Latest Ethereum Gas Price recommendation\n\nBuilt & maintained by Anjan Roy<anjanroy@yandex.com>\n\nFind more about me @ https://itzmeanjan.in")

	})

	bot.Handle("/latest", func(m *telebot.Message) {

		log.Printf("📩 [ /latest ] : From @%s\n", m.Sender.Username)

		// Send latest gas price feed, which was received
		// from `gasz` subscription
		bot.Send(m.Sender, resources.Latest.Sendable())

	})

	fastestTxButton := telebot.InlineButton{
		Unique: "fastest",
		Text:   "Fastest",
	}

	bot.Handle(fastestTxButton, func(c *telebot.Callback) {

		bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})
		bot.Send(c.Sender, "Fastest")

	})

	// -- Buttons for letting user input
	fastTxButton := telebot.InlineButton{
		Unique: "fast",
		Text:   "Fast",
	}

	bot.Handle(fastTxButton, func(c *telebot.Callback) {

		bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})
		bot.Send(c.Sender, "Fast")

	})

	averageTxButton := telebot.InlineButton{
		Unique: "average",
		Text:   "Average",
	}

	bot.Handle(averageTxButton, func(c *telebot.Callback) {

		bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})
		bot.Send(c.Sender, "Average")

	})

	safeLowTxButton := telebot.InlineButton{
		Unique: "safeLow",
		Text:   "SafeLow",
	}

	bot.Handle(safeLowTxButton, func(c *telebot.Callback) {

		bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})
		bot.Send(c.Sender, "SafeLow")

	})

	bot.Handle("/subscribe", func(m *telebot.Message) {

		log.Printf("📩 [ /subscribe ] : From @%s\n", m.Sender.Username)

		bot.Send(m.Sender, "Please choose tx category", &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					fastestTxButton,
					fastTxButton,
				},
				{
					averageTxButton,
					safeLowTxButton,
				},
			},
		})

	})

	log.Printf("✅ Starting @%s\n", bot.Me.Username)

	// This is a blocking call
	bot.Start()

	return nil

}
