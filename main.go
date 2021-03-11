package main

import (
	"log"

	"github.com/itzmeanjan/gasbot/app/bot"
)

func main() {

	log.Printf("🚀 gasbot - Telegram Bot for Ethereum Gas Price Notification")

	if err := bot.Run(); err != nil {
		log.Printf("🚫 Bot stopped : %s\n", err.Error())
	}

}
