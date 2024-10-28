package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var STICKER_SREZKI = "CAACAgIAAxkBAANNZfw-z3y76LA4KRohD4x339CmeS4AAl4tAAL8dcBIPkvpr8s75to0BA"
var STICKER_PENDOS = "CAACAgIAAxkBAANTZfxBwnzzLhYf8gXmtEgGn80uNlwAAp8sAAJCF1lLWBhOY9sCPzs0BA"
var STICKER_PENTA = "CAACAgIAAxkBAAN6ZjfrzN_Udx8D_JGFTAWnYepe66UAAuQ6AAIzJjlKreL0DOfMhNQ1BA"

var f1chat int64 = -1001663174934

var cache *Cache

func main() {
	go func() {
		fs := http.FileServer(http.Dir("./public"))
		http.Handle("/", fs)

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := os.MkdirAll("./public", os.ModePerm)
	if err != nil {
		log.Fatalf("failed creating directory: %s", err)
	}
	cache = NewCache()

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

	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for t := range ticker.C {
			//checkF1Notification(bot)
			fmt.Println("Tick at", t)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		text := update.Message.Text

		processAndSendMessage(text, update, bot)
	}
}

//func checkF1Notification(bot *tgbotapi.BotAPI) {
//	//var f1chat int64 = -1001663174934
//	//nextSession, nextSessionDateString, eventName, found := getNextMessageWithCache()
//	if found {
//		t, _ := time.Parse(time.RFC3339, nextSessionDateString)
//		fmt.Println(nextSessionDateString)
//		//t, _ := time.Parse(time.RFC3339, "2024-03-23T12:23:00+00:00")
//		fmt.Println(t)
//		if time.Until(t) <= 5*time.Minute && time.Until(t) >= 4*time.Minute {
//			//message := tgbotapi.NewMessage(f1chat, eventName+" "+nextSession+" is about to start!")
//			//bot.Send(message)
//		} else if time.Until(t) <= 60*time.Minute && time.Until(t) >= 59*time.Minute {
//			//message := tgbotapi.NewMessage(f1chat, eventName+" "+nextSession+" starts in about one hour!")
//			//bot.Send(message)
//		}
//	}

//}

func getNextMessageWithCache() (string, string, string, bool) {
	url := "https://f1-live-motorsport-data.p.rapidapi.com/races/2024"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", os.Getenv("rapidapi_key"))
	req.Header.Add("X-RapidAPI-Host", "f1-live-motorsport-data.p.rapidapi.com")

	var data []byte
	f1data, cfound := cache.Get(url)
	if cfound {
		data = f1data
	} else {

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		data, _ = io.ReadAll(res.Body)
		cache.Set(url, data)
	}
	var racecalendar RaceCalendar

	err := json.Unmarshal(data, &racecalendar)
	if err != nil {
		fmt.Println("error:", err)
	}
	nextSession, nextSessionDateString, eventName, found := getNextSession(racecalendar)
	return nextSession, nextSessionDateString, eventName, found
}

func processAndSendMessage(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.Chat.ID != f1chat && update.Message.Chat.Title != "holywars" {
		sendMessageByToMe(update, text, bot)
	}
	if containsFixableEmbed(text) {
		var fixedMessage = fixEmbedText(text)
		sendMessageBy(update, fixedMessage, bot)
		deleteMessage(update, bot)
	} else if containsF1NextRequest(text) {
		sendNextF1Session(update, bot)
		//fmt.Printf("%+v", sessions)
	} else if update.Message.Chat.ID == f1chat {
		if textIsAboutCuts(text) {
			share := tgbotapi.NewStickerShare(update.Message.Chat.ID, STICKER_SREZKI)
			bot.Send(share)
		} else if textIsAboutPendos(text) {
			share := tgbotapi.NewStickerShare(update.Message.Chat.ID, STICKER_PENDOS)
			bot.Send(share)
		} else if textIsAboutPenta(text) {
			share := tgbotapi.NewStickerShare(update.Message.Chat.ID, STICKER_PENTA)
			bot.Send(share)
			//sendMessage(update, "ебучий дрогобыщец", bot)
		}
	}

}

func sendNextF1Session(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	htmlContent := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LFM Pro Series Stats</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            margin: 0;
            padding: 0;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 2em auto;
            padding: 1em;
        }
        h1 {
            text-align: center;
            color: #333;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            font-size: 0.9em;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        table thead {
            background-color: #009879;
            color: #ffffff;
            cursor: pointer;
        }
        table th, table td {
            padding: 12px 15px;
            text-align: center;
        }
        table th.sortable:hover {
            background-color: #007965;
        }
        table tbody tr:nth-child(even) {
            background-color: #f3f3f3;
        }
        table tbody tr:hover {
            background-color: #e9f1f7;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>LFM Pro Series Stats</h1>
        <table>
            <thead>
                <tr>
                    <th class="sortable" onclick="sortTable(0)">First Name</th>
                    <th class="sortable" onclick="sortTable(1)">Last Name</th>
                    <th class="sortable" onclick="sortTable(2)">Races</th>
                    <th class="sortable" onclick="sortTable(3)">Wins</th>
                    <th class="sortable" onclick="sortTable(4)">Podiums</th>
                    <th class="sortable" onclick="sortTable(5)">Poles</th>
                    <th class="sortable" onclick="sortTable(6)">Top 5s</th>
                    <th class="sortable" onclick="sortTable(7)">Top 10s</th>
                    <th class="sortable" onclick="sortTable(8)">Avg Finish</th>
                    <th class="sortable" onclick="sortTable(9)">Avg Qualify</th>
                    <th class="sortable" onclick="sortTable(10)">Points</th>
                    <th class="sortable" onclick="sortTable(11)">Winrate</th>
                    <th class="sortable" onclick="sortTable(12)">Podium Rate</th>
                </tr>
            </thead>
            <tbody>
               <REPLACE>
            </tbody>
        </table>
    </div>
    <script>
        function sortTable(columnIndex) {
            var table = document.querySelector("table tbody");
            var rows = Array.from(table.rows);
            var ascending = table.getAttribute("data-sort-order") !== "asc";

            rows.sort(function(rowA, rowB) {
                var cellA = rowA.cells[columnIndex].innerText;
                var cellB = rowB.cells[columnIndex].innerText;

                var numA = parseFloat(cellA) || cellA;
                var numB = parseFloat(cellB) || cellB;

                return ascending ? numA > numB ? 1 : -1 : numA < numB ? 1 : -1;
            });

            table.innerHTML = "";
            rows.forEach(row => table.appendChild(row));
            table.setAttribute("data-sort-order", ascending ? "asc" : "desc");
        }
    </script>
</body>
</html>

	`

	file, _ := os.Open("./public/lfm_proseries_data.csv")
	reader := csv.NewReader(file)
	// Create a StringBuilder to build HTML
	var html strings.Builder
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		html.WriteString("\t<tr>\n")
		for _, value := range record {
			html.WriteString(fmt.Sprintf("\t\t<td>%s</td>\n", value))
			log.Println("writing line: " + value)
		}
		html.WriteString("\t</tr>\n")
	}
	htmlContent = strings.Replace(htmlContent, "<REPLACE>", html.String(), -1)
	log.Println("writing to file")
	err := ioutil.WriteFile("./public/output.html", []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

}

func textIsAboutPenta(text string) bool {
	text = strings.ToLower(text)
	if strings.Contains(text, "пента") {
		return true
	}
	if strings.Contains(text, "пєнта") {
		return true
	}
	if strings.Contains(text, "пенти") {
		return true
	}
	if strings.Contains(text, "пєнти") {
		return true
	}
	return false
}

func textIsAboutPendos(text string) bool {
	text = strings.ToLower(text)
	if strings.Contains(text, "пендос") {
		return true
	}
	if strings.Contains(text, "пєндос") {
		return true
	}
	if strings.Contains(text, "піндос") {
		return true
	}
	if strings.Contains(text, "пиндос") {
		return true
	}
	return false
}

func textIsAboutCuts(text string) bool {
	text = strings.ToLower(text)
	if strings.Contains(text, "срезки") {
		return true
	}
	if strings.Contains(text, "зрізки") {
		return true
	}
	return false
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

func sendMessageByToMe(update tgbotapi.Update, fixedMessage string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(6991628262, update.Message.From.UserName+" sent: "+fixedMessage)
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
		strings.Contains(text, "://www.instagram.com/reel") ||
		strings.Contains(text, "://www.tiktok.com")
}
func fixEmbedText(text string) string {

	text = strings.ReplaceAll(text, "://twitter.com/", "://fxtwitter.com/")
	text = strings.ReplaceAll(text, "://x.com/", "://fixupx.com/")
	text = strings.ReplaceAll(text, "://www.instagram.com/", "://www.ddinstagram.com/")
	text = strings.ReplaceAll(text, "://www.tiktok.com/", "://www.tiktokez.com/")
	return text
}
