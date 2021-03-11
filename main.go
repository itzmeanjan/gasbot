package main

import (
	"log"

	"github.com/itzmeanjan/gasbot/app/gasz"
)

func main() {

	log.Printf("🚀 gasbot - Telegram Bot for Ethereum Gas Price Notification")

	gasPrice, err := gasz.CurrentGasPrice()
	if err != nil {
		log.Printf("🚫 Failed to get latest gas price : %s\n", err.Error())
		return
	}

	log.Printf("%v\n", gasPrice)

}
