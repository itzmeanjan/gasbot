package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/itzmeanjan/gasbot/app/bot"
	"github.com/itzmeanjan/gasbot/app/config"
	"github.com/itzmeanjan/gasbot/app/data"
	"github.com/itzmeanjan/gasbot/app/gasz"
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
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT)

	comm := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		select {

		case <-interruptChan:

			cancel()
			<-time.After(time.Duration(1) * time.Second)

		case <-comm:

			// @note Need to handle it better
			// New subscriber can be spawned
			cancel()

		}

		// Stopping process
		log.Printf("\n✅ Gracefully shut down `gasbot`\n")
		os.Exit(0)

	}()

	resources := &data.Resources{
		Latest:        &data.CurrentGasPrice{},
		Subscriptions: make(map[string]*data.Subscriber),
		Lock:          &sync.RWMutex{},
	}

	// Starting gas price subscriber as different go routine
	go gasz.SubscribeToLatest(ctx, comm, resources)

	// Main go routine works as bot
	if err := bot.Run(resources); err != nil {
		log.Printf("🚫 Bot stopped : %s\n", err.Error())
	}

}
