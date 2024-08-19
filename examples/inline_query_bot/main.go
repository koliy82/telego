package main

import (
	"fmt"
	"os"

	"github.com/koliy82/telego"
	tu "github.com/koliy82/telego/telegoutil"
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
	defer bot.StopLongPolling()

	for update := range updates {
		// Receive inline query
		if update.InlineQuery != nil {
			iq := update.InlineQuery
			name := iq.From.FirstName

			// Answer inline query request
			_ = bot.AnswerInlineQuery(tu.InlineQuery(
				iq.ID,

				tu.ResultArticle(
					"hello",
					"Hello",
					tu.TextMessage( // Hello message with an inline query
						fmt.Sprintf("Hello %s\n\nYour query:\n```%#+v```", name, iq),
					).WithParseMode(telego.ModeMarkdownV2),
				).WithDescription(fmt.Sprintf("Query: %q", iq.Query)),

				tu.ResultArticle(
					"bey",
					"Bye",
					tu.TextMessage( // Bye message with an inline query
						fmt.Sprintf("Bye %s\n\nYour query:\n```%#+v```", name, iq),
					).WithParseMode(telego.ModeMarkdownV2),
				).WithDescription(fmt.Sprintf("Query: %q", iq.Query)),
			).WithIsPersonal())
		}
	}
}
