package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/itzmeanjan/gasbot/app/bot"
	"github.com/itzmeanjan/gasbot/app/config"
)

func main() {

	log.Printf("🚀 gasbot - Telegram Bot for Ethereum Gas Price Notification")

	abs, err := filepath.Abs(".env")
	if err != nil {

		log.Printf("❗️ Failed to find absolute path of `.env` file : %s\n", err.Error())
		os.Exit(1)

	}

	if err := config.Read(abs); err != nil {

		log.Printf("❗️ Failed read `.env` file : %s\n", err.Error())
		os.Exit(1)

	}

	// Channel for catching interrupts
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go func() {

		// This is a blocking call
		//
		// To be unblocked only when interrupt is received by this process
		<-interruptChan

		// Stopping process
		log.Printf("\n✅ Gracefully shut down `gasbot`\n")
		os.Exit(0)

	}()

	if err := bot.Run(); err != nil {
		log.Printf("🚫 Bot stopped : %s\n", err.Error())
	}

}
