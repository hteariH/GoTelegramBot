package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
)

func main() {
	botToken := os.Getenv("LinkFixer_Bot_token")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		text := update.Message.Text
		if isFixeable(text) {
			var fixedMessage = fixmsg(text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.From.UserName+" sent: "+fixedMessage)
			bot.Send(msg)
			deleteConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			_, err := bot.DeleteMessage(deleteConfig)
			if err != nil {
				log.Println("Unable to delete message:", err)
			}
		}
	}
}

func isFixeable(text string) bool {
	return strings.Contains(text, "://twitter.com/") ||
		strings.Contains(text, "://x.com/") ||
		strings.Contains(text, "://www.instagram.com/p/") ||
		strings.Contains(text, "://www.instagram.com/reel")
}
func fixmsg(text string) string {

	text = strings.ReplaceAll(text, "://twitter.com/", "://fxtwitter.com/")
	text = strings.ReplaceAll(text, "://x.com/", "://fixupx.com/")
	text = strings.ReplaceAll(text, "://www.instagram.com/", "://www.ddinstagram.com/")
	return text
}
