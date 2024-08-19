package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/koliy82/telego"
	th "github.com/koliy82/telego/telegohandler"
	tu "github.com/koliy82/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Create bot: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Note: Creating a secret token like this is not secure,
	// but at least better than having a plain bot token as is in requests
	secretBytes := sha256.Sum256([]byte(botToken))
	secret := hex.EncodeToString(secretBytes[:])

	srv, url := Webhook(ctx, bot, secret)

	updates, err := bot.UpdatesViaWebhook(
		"/bot",
		telego.WithWebhookServer(srv),
		telego.WithWebhookSet(tu.Webhook(url+"/bot").WithSecretToken(secret)),
	)
	if err != nil {
		log.Fatalf("Updates via webhoo: %s", err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatalf("Bot handler: %s", err)
	}

	RegisterHandlers(bh)

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		log.Println("Stopping...")

		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer stopCancel()

		err = bot.StopWebhookWithContext(stopCtx)
		if err != nil {
			log.Println("Failed to stop webhook properly:", err)
		}

		bh.StopWithContext(stopCtx)

		done <- struct{}{}
	}()

	go bh.Start()
	log.Println("Handling updates...")

	go func() {
		err = bot.StartWebhook(":443")
		if err != nil {
			log.Fatalf("Failed to start webhook: %s", err)
		}
	}()

	<-done
	log.Println("Done")
}
