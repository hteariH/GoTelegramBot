package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var STICKER_SREZKI = "CAACAgIAAxkBAANNZfw-z3y76LA4KRohD4x339CmeS4AAl4tAAL8dcBIPkvpr8s75to0BA"
var STICKER_PENDOS = "CAACAgIAAxkBAANTZfxBwnzzLhYf8gXmtEgGn80uNlwAAp8sAAJCF1lLWBhOY9sCPzs0BA"
var STICKER_PENTA = "CAACAgIAAxkBAAN6ZjfrzN_Udx8D_JGFTAWnYepe66UAAuQ6AAIzJjlKreL0DOfMhNQ1BA"
var messages = []string{
	"До перемоги України у війні над Росією залишилося всього %s днів. Цей день стане новою сторінкою в історії нашої країни та усього вільного світу.",
	"Через %s днів завершиться одна з найтрагічніших сторінок в історії України, і над країною запанує мир, освячений нашою перемогою над агресором.",
	"Кожен із цих %s днів наближає нас до моменту, коли Україна остаточно звільниться від гніту та насильства з боку агресора. Ми на порозі великої перемоги!",
	"Вже зовсім скоро, через %s днів, Україна поставить крапку в цій війні, продемонструвавши всьому світу свою силу, стійкість та прагнення до свободи.",
	"Залишилося всього %s днів до того моменту, коли український народ з гордістю скаже: ми перемогли, ми вистояли, ми вільні!",
	"Кожен день боротьби наближає нас до перемоги. Всього через %s днів Україна святкуватиме довгоочікуваний мир, здобутий ціною величезних зусиль.",
	"Через %s днів Україна покаже всьому світу, що сила духу, єдність та правда здатні перемогти будь-яку агресію. Цей день вже близько!",
	"Всього %s днів відділяють нас від моменту, коли світ побачить: Україна не лише вистояла, але й здобула тріумфальну перемогу в цій жорстокій війні.",
	"Коли минуть ці %s днів, Україна остаточно утвердить свою свободу, і весь світ захоплюватиметься стійкістю нашого народу. Перемога вже на горизонті!",
	"Залишилося всього %s днів, щоб над Україною знову засяяло сонце миру та свободи, ознаменувавши перемогу над ворогом і повернення спокійного життя.",
	"Через %s днів ми зможемо сказати: наша боротьба увінчалася перемогою, і майбутнє України стало світлим і вільним.",
	"Останні %s днів війни — це час, коли український народ демонструє свою незламність, готуючись до перемоги.",
	"Ще трохи терпіння, ще %s днів — і Україна відсвяткує тріумфальне завершення війни.",
	"Ми стоїмо на порозі великого дня, який настане через %s днів і стане символом перемоги України.",
	"Через %s днів весь світ стане свідком того, як Україна здобула заслужену перемогу в боротьбі за свободу.",
	"Залишилося всього %s днів до того моменту, коли ми зможемо пишатися здобутою перемогою і мирним життям.",
	"Ще %s днів — і Україна увійде в нову еру миру та відновлення після великої перемоги над агресором.",
	"Кожен день наближає нас до перемоги. Ще %s днів, і мрія про вільну Україну стане реальністю.",
	"Україна вже на шляху до перемоги. Через %s днів ми завершимо цю війну на нашу користь.",
	"Через %s днів ми святкуватимемо перемогу, яка стане доказом сили, єдності та мужності українського народу.",
	"Ще %s днів, і Україна відновить мир та свободу, які були несправедливо порушені.",
	"Через %s днів ми разом відзначимо перемогу, яка стала результатом нашої спільної боротьби.",
	"Останні %s днів війни — це час, коли наша країна демонструє незламну віру в перемогу.",
	"Ще %s днів, і мирне небо над Україною стане символом нашої перемоги та свободи.",
	"Через %s днів український народ відзначатиме день, який стане початком нового етапу в історії країни.",
	"Ми на порозі історичного моменту. Ще %s днів — і Україна здобуде заслужену перемогу.",
	"Кожен день боротьби наближає нас до перемоги. Ще %s днів терпіння та віри.",
	"Україна переможе через %s днів, і ця перемога стане тріумфом добра над злом.",
	"Залишилося всього %s днів до моменту, коли ми зможемо зітхнути з полегшенням і сказати: ми перемогли.",
	"Через %s днів ми побачимо кінець війни і початок нової мирної епохи для України.",
	"Ще трохи, ще %s днів — і Україна переможе, відновивши свою незалежність і територіальну цілісність.",
	"Через %s днів ми зможемо з гордістю сказати: ми вистояли, ми перемогли, ми вільні.",
	"Кожен день боротьби — це ще один крок до перемоги. Залишилося всього %s днів.",
	"Останні %s днів війни — це час, коли віра в перемогу стає ще сильнішою.",
	"Ще трохи, ще %s днів — і Україна відзначатиме день перемоги над агресором.",
	"Через %s днів український прапор гордо майорітиме над звільненими територіями.",
	"Залишилося всього %s днів до того, як мирне життя повернеться в кожен український дім.",
	"Ми стоїмо на порозі великої перемоги. Ще %s днів боротьби та віри.",
	"Через %s днів український народ святкуватиме свободу, здобуту ціною величезних зусиль.",
	"Ще трохи терпіння — через %s днів наша перемога стане реальністю.",
	"Україна вже на шляху до тріумфу. Залишилося всього %s днів до перемоги.",
	"Через %s днів ми відзначимо день, який увійде в історію як день перемоги над агресором.",
	"Останні %s днів війни — це час, коли віра в мир і свободу надихає нас.",
	"Ще %s днів, і Україна стане прикладом незламності для всього світу.",
	"Через %s днів наша боротьба завершиться перемогою, яка змінить хід історії.",
	"Залишилося всього %s днів до моменту, коли мир і спокій повернуться в Україну.",
}
var f1chat int64 = -1001663174934

var freedom int64 = -1002266390232

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

	chats := make(map[int64]bool)

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	location, _ := time.LoadLocation("Europe/Kiev")
	now := time.Now()
	nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), 13, 10, 0, 0, location)
	if nextRunTime.Before(now) { // if it's after 14:00, schedule for 14:00 of next day.
		nextRunTime = nextRunTime.Add(24 * time.Hour)
	}
	delay := nextRunTime.Sub(now)
	ticker2 := time.NewTicker(delay)
	//ticker2 := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker2.C:
				futureDate := time.Date(2025, time.May, 2, 0, 0, 0, 0, time.UTC)
				daysLeft := futureDate.Sub(time.Now()).Hours() / 24
				if daysLeft > 0 {
					msg := messages[rand.Intn(len(messages))]
					msg = strings.Replace(msg, "%s", fmt.Sprintf("%.0f", daysLeft), -1)
					fmt.Println(msg)
					for chatID := range chats {
						message := tgbotapi.NewMessage(chatID, msg)
						_, err := bot.Send(message)
						if err != nil {
							log.Printf("Failed to send message to chat ID %d: %v", chatID, err)
						}
					}

					ticker2.Reset(24 * time.Hour)
				} else if daysLeft == 0 {
					for chatID := range chats {
						message := tgbotapi.NewMessage(chatID, "Вітаю, Україна Перемогла!!!")
						_, err := bot.Send(message)
						if err != nil {
							log.Printf("Failed to send message to chat ID %d: %v", chatID, err)
						}
					}

					ticker2.Stop() // Stops the ticker when the future date has reached.
				} else {
					ticker2.Stop()
				}
			}
		}
	}()
	sendNextF1Session()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		chats[chatID] = true
		log.Printf("Current unique chats: %v", chats)

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
	if update.Message.Chat.Title != "holywars" && update.Message.Chat.ID != -1002266390232 && update.Message.Chat.ID != -1001663174934 {
		sendMessageByToMe(update, text, bot)
	}
	if containsFixableEmbed(text) {
		var fixedMessage = fixEmbedText(text)
		sendMessageBy(update, fixedMessage, bot)
		deleteMessage(update, bot)
	} else if containsF1NextRequest(text) {
		sendNextF1Session()
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

func sendNextF1Session() {
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
            background-color: #282b30;
            margin: 0;
            padding: 0;
            color: white;
        }
        .container {
            max-width: 1200px;
            padding: 50px;
        }
        h1 {
            text-align: center;
            color: #white;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            font-size: 0.9em;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        table thead {
            background-color: #424549;
			color: #ffffff;
            cursor: pointer;
        }
        table th, table td {
            padding: 12px 10px;
            text-align: center;
        }
        table th.sortable:hover {
            background-color: #007965;
        }
        table tbody tr:nth-child(even) {
            background-color: #424549;
        }
        table tbody tr:hover {
            background-color: #4a4c51;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>LFM Pro Series Stats</h1>
        <table>
            <thead> 
				<REPLACE0>
            </thead>
            <tbody>
				<REPLACE1>
            </tbody>
        </table>
    </div>
    <script>
       function sortTable(columnIndex) {
    var table = document.querySelector("table tbody");
    var rows = Array.from(table.rows);
    var ascending = table.getAttribute("data-sort-order") !== "asc";

    rows.sort(function (rowA, rowB) {
        var cellA = rowA.cells.length > columnIndex? rowA.cells[columnIndex].innerText : '';
        var cellB = rowB.cells.length > columnIndex? rowB.cells[columnIndex].innerText : '';

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

	resp, err := http.Get("https://gitlab.com/pst-pepega/pst-scripts/-/raw/main/LFM_Proseries_Data/lfm_proseries_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	var html strings.Builder
	var headerHtml strings.Builder
	firstLine := true
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		headerHtml.WriteString("\t<tr>\n")
		html.WriteString("\t<tr>\n")
		i := 0
		for _, value := range record {

			if firstLine {
				headerHtml.WriteString(fmt.Sprintf("\t\t<th class=\"sortable\" onclick=\"sortTable(%d)\">%s</th>\n", i, value))
				i++
				log.Println("writing first line: " + value)
			} else {
				if i < 10 || i == 18 || i == 20 {
					html.WriteString(fmt.Sprintf("\t\t<td>%s</td>\n", value))

				} else {
					num, err := strconv.ParseFloat(value, 64)
					if err != nil {
						fmt.Println(err)
						html.WriteString(fmt.Sprintf("\t\t<td>%s</td>\n", value))
						//return
					} else {
						html.WriteString(fmt.Sprintf("\t\t<td>%.2f</td>\n", num))
					}

				}
				log.Println("writing line: " + value)
				i++
				//}
			}

		}
		if firstLine {
			firstLine = false
		}
		headerHtml.WriteString("\t</tr>\n")
		html.WriteString("\t</tr>\n")
	}
	htmlContent = strings.Replace(htmlContent, "<REPLACE0>", headerHtml.String(), -1)
	htmlContent = strings.Replace(htmlContent, "<REPLACE1>", html.String(), -1)
	log.Println("writing to file")
	err = ioutil.WriteFile("./public/lfm_proseries_data.html", []byte(htmlContent), 0644)
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
	msg := tgbotapi.NewMessage(6991628262, update.Message.Chat.Title+":"+update.Message.From.UserName+" sent: "+fixedMessage)
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
		strings.Contains(text, ".tiktok.com")
}
func fixEmbedText(text string) string {

	text = strings.ReplaceAll(text, "://twitter.com/", "://fxtwitter.com/")
	text = strings.ReplaceAll(text, "://x.com/", "://fixupx.com/")
	text = strings.ReplaceAll(text, "://www.instagram.com/", "://www.ddinstagram.com/")
	text = strings.ReplaceAll(text, ".tiktok.com/", ".tiktokez.com/")
	return text
}
