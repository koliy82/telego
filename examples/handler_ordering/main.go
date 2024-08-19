package main

import (
	"fmt"
	"os"

	"github.com/koliy82/telego"
	th "github.com/koliy82/telego/telegohandler"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	// Stop handling updates
	defer bh.Stop()

	// Stop getting updates
	defer bot.StopLongPolling()

	// Should not be here, because the order of handlers does meter.
	//
	// bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
	// 	 fmt.Println("Message:", message.Text)
	// })
	//
	// When you are defining handlers only the first matched handler will process update, that means that in this
	// specific example, matching any message will "shadow" any other handler that matches on messages.
	// For any message it will match and handling of start command will never happen.
	//
	// The general idea is that you should define the most specific handlers first (like for specific command or etc.)
	// and only then more generic handlers (like any message).

	// Will match any message with command `/start`
	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		fmt.Println("Start")
	}, th.CommandEqual("start"))

	// Will match to any message
	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		fmt.Println("Message:", message.Text)
	})

	// Start handling
	bh.Start()
}
