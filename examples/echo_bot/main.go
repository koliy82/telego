package main

import (
	"fmt"
	"os"

	"github.com/koliy82/telego"
	tu "github.com/koliy82/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// Stop reviving updates from update channel
	defer bot.StopLongPolling()

	// Loop through all updates when they came
	for update := range updates {
		// Check if update contains a message
		if update.Message != nil {
			// Get chat ID from the message
			chatID := tu.ID(update.Message.Chat.ID)

			// Copy sent messages back to the user
			_, _ = bot.CopyMessage(
				tu.CopyMessage(
					chatID,
					chatID,
					update.Message.MessageID,
				),
			)
		}
	}
}
