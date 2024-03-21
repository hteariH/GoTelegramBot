package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

		processAndSendMessage(text, update, bot)
	}
}

func processAndSendMessage(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if containsFixableEmbed(text) {
		var fixedMessage = fixEmbedText(text)
		sendMessageBy(update, fixedMessage, bot)
		deleteMessage(update, bot)
	}

	if containsF1NextRequest(text) {
		url := "https://f1-live-motorsport-data.p.rapidapi.com/races/2024"
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", os.Getenv("rapidapi_key"))
		req.Header.Add("X-RapidAPI-Host", "f1-live-motorsport-data.p.rapidapi.com")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		var racecalendar RaceCalendar

		err := json.Unmarshal(data, &racecalendar)
		if err != nil {
			fmt.Println("error:", err)
		}
		nextSession, nextSessionDateString, eventName, found := getNextSession(racecalendar)

		if found {
			date := formatDate(nextSessionDateString)
			message := fmt.Sprintf("The next session is %s %s on time %s\n", eventName, nextSession, date)
			sendMessage(update, message, bot)
		} else {
			fmt.Println("There are no upcoming sessions.")
			sendMessage(update, "There are no upcoming sessions.", bot)
		}

		//fmt.Printf("%+v", sessions)
	}

}

func formatDate(dateString string) string {

	// Parsing the time from the string.
	t, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Now, use time.Format to write the time in the format you wanted.
	// Note: In layout string you need an example of input date 'Mon Jan 2 15:04:05 MST 2006' with values replaced what you want to see in your converted date.
	// For example if you want to see '15' in 'hour' place of date, replace '15' in layout with 'hour' and so on for other parts of date.
	location, _ := time.LoadLocation("Europe/Kyiv")
	t = t.In(location)
	output := t.Format("15:04 02-01-2006")
	fmt.Printf("The reformatted time is: %s\n", output)
	return output
}

func getNextSession(data RaceCalendar) (next string, nextDate string, eventName string, found bool) {
	// Get current time.
	// You can replace this with any time you're comparing to
	now := time.Now()

	for _, result := range data.Results {
		for _, session := range result.Sessions {
			timestamp, err := time.Parse(time.RFC3339, session.Date)
			if err != nil {
				fmt.Printf("An error occurred while parsing date: %s", err)
				return "", "", "", false
			}
			if timestamp.After(now) {
				return session.SessionName, session.Date, result.Name, true
			}
		}
	}

	return "", "", "", false
}

func containsF1NextRequest(text string) bool {
	return strings.HasPrefix(text, "/f1next")
}

func sendMessageBy(update tgbotapi.Update, fixedMessage string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.From.UserName+" sent: "+fixedMessage)
	bot.Send(msg)
}

func sendMessage(update tgbotapi.Update, fixedMessage string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fixedMessage)
	bot.Send(msg)
}

func deleteMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	deleteConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err := bot.DeleteMessage(deleteConfig)
	if err != nil {
		log.Println("Unable to delete message:", err)
	}
}

func containsFixableEmbed(text string) bool {
	return strings.Contains(text, "://twitter.com/") ||
		strings.Contains(text, "://x.com/") ||
		strings.Contains(text, "://www.instagram.com/p/") ||
		strings.Contains(text, "://www.instagram.com/reel")
}
func fixEmbedText(text string) string {

	text = strings.ReplaceAll(text, "://twitter.com/", "://fxtwitter.com/")
	text = strings.ReplaceAll(text, "://x.com/", "://fixupx.com/")
	text = strings.ReplaceAll(text, "://www.instagram.com/", "://www.ddinstagram.com/")
	return text
}
