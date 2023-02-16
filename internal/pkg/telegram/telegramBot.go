package telegram

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
)

type tool struct {
}

func newTool() *tool {
	return &tool{}
}

func (d *tool) formatDate(t time.Time) string {
	return t.Add(3 * time.Hour).String()[:19]
}
func sanitize (s string) string{
	myString := s
myString = strings.ReplaceAll(myString, "\\", "")
myString = strings.ReplaceAll(myString, "`", "")
myString = strings.ReplaceAll(myString, "*", "")
myString = strings.ReplaceAll(myString, "_", "")
myString = strings.ReplaceAll(myString, "[", "")
myString = strings.ReplaceAll(myString, "]", "")
myString = strings.ReplaceAll(myString, "(", "")
myString = strings.ReplaceAll(myString, ")", "")
myString = strings.ReplaceAll(myString, "#", "")
myString = strings.ReplaceAll(myString, "+", "")
myString = strings.ReplaceAll(myString, "-", "")
myString = strings.ReplaceAll(myString, ".", "")
myString = strings.ReplaceAll(myString, "!", "")
myString = strings.ReplaceAll(myString, "@", "")
myString = strings.ReplaceAll(myString, ",", "")
myString = strings.ReplaceAll(myString, "'", "")
return myString
 }

 func (t *tool)escapeSimbols (s string) string{
	myString := s
myString = strings.ReplaceAll(myString, "\\", "\\\\")
myString = strings.ReplaceAll(myString, "`", "\\`")
myString = strings.ReplaceAll(myString, "*", "\\*")
myString = strings.ReplaceAll(myString, "_", "\\_")
myString = strings.ReplaceAll(myString, "[", "\\[")
myString = strings.ReplaceAll(myString, "]", "\\]")
myString = strings.ReplaceAll(myString, "(", "\\(")
myString = strings.ReplaceAll(myString, ")", "\\)")
myString = strings.ReplaceAll(myString, "#", "\\#")
myString = strings.ReplaceAll(myString, "+", "\\+")
myString = strings.ReplaceAll(myString, "-", "\\-")
myString = strings.ReplaceAll(myString, ".", "\\.")
myString = strings.ReplaceAll(myString, "!", "\\!")
return myString
 }
const badQuery = `
I am sorry, but I can't find anythithg like _'{{.Text}}' _

Maybe we should try again?
Please print your city and neigbourhood divided by space, for example _'Fethie Taşyaka'_"
`
const listOutages = `
{{if eq (len .) 0 }}
*There is no outages planned in your neigborhood in the closest time*
{{else}}
*Here are the closest outages found for your neigborhood:*
{{range .}}
*{{.Resource}} outage* from {{escape (format (.StartDate))}} to {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}

*In the next areas and streets:*
{{escape (.Notes)}}{{end}}
{{end}}
{{end}}`
const confirmDistr = `
Did you mean _*{{.City}} {{.Name}}*_?`

func buildAnswer(d district.District, o []outage.Outage) (string, error) {
	var buffer bytes.Buffer
	var err error
	confifmDistrTemp := template.Must(template.New("confifmDistrTemp").Parse(confirmDistr))
	if err := confifmDistrTemp.Execute(&buffer, d); err != nil {
		return "Error", err
	}
	tool := newTool()
	listOutagesTemp := template.New("listOutagesTemp").Funcs(template.FuncMap{
		"escape": tool.escapeSimbols,
		"format": tool.formatDate,
	})
	listOutagesTemp, err = listOutagesTemp.Parse(listOutages)
	if err != nil {
		return "Error", err
	}
	if err := listOutagesTemp.Execute(&buffer, o); err != nil {
		return "Error", err
	}
	return buffer.String(), err
}
func BotRunner(ds *district.DistrictStore, store *postgres.OutageStore) {

	var buffer bytes.Buffer
	// подключаемся к боту с помощью токена
	api := os.Getenv("OUTAGE_TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		fmt.Println("telegram ApI error", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if update.Message.Text == "/start" {
				msg.Text = "Please print your city and neigbourhood divided by space, for example 'Fethie Taşyaka'"
			} else {

				guessDistr, err := ds.GetFuzzyMatch(sanitize(update.Message.Text))
				if err != nil {
					log.Fatal(err)
				}
				userOutages, err := store.GetActiveOutagesByCityDistrict(guessDistr.Name, guessDistr.City)
				if err != nil {
					log.Fatal(err)
				}
				if guessDistr.City == "no matches" {
					badQueryTemp := template.Must(template.New("badQueryTemp").Parse(badQuery))
					if err := badQueryTemp.Execute(&buffer, update.Message); err != nil {
						log.Fatal(err)
					}
					msg.Text = buffer.String()
				} else {
					msg.Text, err = buildAnswer(guessDistr, userOutages)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
						msg.ParseMode = "MarkdownV2" //This parse mode enables format tags in TG
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
