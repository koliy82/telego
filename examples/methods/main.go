package main

import (
	"fmt"
	"os"

	"github.com/koliy82/telego"
	tu "github.com/koliy82/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Create Bot
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Call method getMe
	botUser, _ := bot.GetMe()
	fmt.Printf("Bot User: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	for update := range updates {
		if update.Message != nil {
			// Retrieve chat ID
			chatID := update.Message.Chat.ID

			// Call method sendMessage (https://core.telegram.org/bots/api#sendmessage).
			// Send a message to sender with the same text (echo bot).
			sentMessage, _ := bot.SendMessage(
				tu.Message(
					tu.ID(chatID),
					update.Message.Text,
				),
			)

			fmt.Printf("Sent Message: %v\n", sentMessage)
		}
	}
}
