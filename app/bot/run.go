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

	// This is step two of subscription process, where
	// user is asked to put relational operator to be used
	// when checking whether some gas price update needs to be
	// pushed to them or not
	//
	// <, >, <=, >=, == these 5 are allowed operators
	subStepTwo := func() *telebot.ReplyMarkup {

		// -- Buttons for letting user input their choice
		lesserThan := telebot.InlineButton{
			Unique: "lesserThan",
			Text:   "<",
		}

		bot.Handle(&lesserThan, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.SetSubscriptionOperator(c.Sender, "<") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please send threshold gas price in Gwei")

		})

		greaterThan := telebot.InlineButton{
			Unique: "greaterThan",
			Text:   ">",
		}

		bot.Handle(&greaterThan, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.SetSubscriptionOperator(c.Sender, ">") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please send threshold gas price in Gwei")

		})

		lesserThanOrEqualsTo := telebot.InlineButton{
			Unique: "lesserThanOrEqualsTo",
			Text:   "<=",
		}

		bot.Handle(&lesserThanOrEqualsTo, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.SetSubscriptionOperator(c.Sender, "<=") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please send threshold gas price in Gwei")

		})

		greaterThanOrEqualsTo := telebot.InlineButton{
			Unique: "greaterThanOrEqualsTo",
			Text:   ">=",
		}

		bot.Handle(&greaterThanOrEqualsTo, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.SetSubscriptionOperator(c.Sender, ">=") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please send threshold gas price in Gwei")

		})

		equalsTo := telebot.InlineButton{
			Unique: "equalsTo",
			Text:   "==",
		}

		bot.Handle(&equalsTo, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.SetSubscriptionOperator(c.Sender, "==") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please send threshold gas price in Gwei")

		})
		// -- Buttons end here, along with their respective handler
		// definitions

		return &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					lesserThan,
					greaterThan,
				},
				{
					lesserThanOrEqualsTo,
					greaterThanOrEqualsTo,
				},
				{
					equalsTo,
				},
			},
		}

	}

	// This is very first step of subscription
	//
	// Here user is asked to select what is his/ her
	// tx type category for which they would like to get notified
	subStepOne := func() *telebot.ReplyMarkup {

		// -- Buttons for letting user input their choice
		fastestTxButton := telebot.InlineButton{
			Unique: "fastest",
			Text:   "Fastest",
		}

		bot.Handle(&fastestTxButton, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.InitSubscription(c.Sender, "fastest") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please choose evaluation operator", subStepTwo())

		})

		fastTxButton := telebot.InlineButton{
			Unique: "fast",
			Text:   "Fast",
		}

		bot.Handle(&fastTxButton, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.InitSubscription(c.Sender, "fast") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please choose evaluation operator", subStepTwo())

		})

		averageTxButton := telebot.InlineButton{
			Unique: "average",
			Text:   "Average",
		}

		bot.Handle(&averageTxButton, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.InitSubscription(c.Sender, "average") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please choose evaluation operator", subStepTwo())

		})

		safeLowTxButton := telebot.InlineButton{
			Unique: "safeLow",
			Text:   "SafeLow",
		}

		bot.Handle(&safeLowTxButton, func(c *telebot.Callback) {

			bot.Respond(c, &telebot.CallbackResponse{ShowAlert: false})

			if !resources.InitSubscription(c.Sender, "safeLow") {

				bot.Edit(c.Message, "Something went wrong 😞")
				return

			}

			bot.Edit(c.Message, "Please choose evaluation operator", subStepTwo())

		})
		// -- Buttons end here, along with their respective handler
		// definitions

		return &telebot.ReplyMarkup{
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
		}

	}

	bot.Handle("/subscribe", func(m *telebot.Message) {

		log.Printf("📩 [ /subscribe ] : From @%s\n", m.Sender.Username)

		bot.Send(m.Sender, "Please choose tx category", subStepOne())

	})

	log.Printf("✅ Starting @%s\n", bot.Me.Username)

	// This is a blocking call
	bot.Start()

	return nil

}
