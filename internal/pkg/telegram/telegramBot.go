package telegram

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
)

type tool struct {
}

func newTool() *tool {
	return &tool{}
}

func (t *tool) formatDateAndMakeLocal(tm time.Time) string {
	return tm.Add(3 * time.Hour).String()[:19]
}

func (t *tool) sanitize(s string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	s = re.ReplaceAllString(s, " ")
	return s
}

func (t *tool) escapeSimbols(s string) string {
	re := regexp.MustCompile(`[\\` + "`*_\\[\\]()#+\\-.!]")
    return re.ReplaceAllStringFunc(s, func(match string) string {
        return "\\" + match
    })
}

func BotRunner(ds *district.DistrictStore, store *postgres.OutageStore) {
	var err error
	tool := newTool()
	//mapping the functions for templates
	dialogTemplape := template.New("dialogTemplape").Funcs(template.FuncMap{
		"escape": tool.escapeSimbols,
		"format": tool.formatDateAndMakeLocal,
	})
	//parsing the template file
	t, err := dialogTemplape.ParseFiles("dialog_templates.tpl")
	if err != nil {
		log.Fatal(err)
	}
	// reading the token from envirinment and connecting
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
		var buffer bytes.Buffer
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if update.Message.Text == "/start" {
				if err := t.ExecuteTemplate(&buffer, "startMsg", nil); err != nil {
					log.Fatal(err)
				}
			} else {
				guessDistr, err := ds.GetFuzzyMatch(tool.sanitize(update.Message.Text))
				if err != nil {
					log.Fatal(err)
				}
				userOutages, err := store.GetActiveOutagesByCityDistrict(guessDistr.Name, guessDistr.City)
				if err != nil {
					log.Fatal(err)
				}
				if guessDistr.City == "no matches" {
					if err := t.ExecuteTemplate(&buffer, "badQuery", update.Message); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := t.ExecuteTemplate(&buffer, "confirmDistr", guessDistr); err != nil {
						log.Fatal(err)
					}
					if err := t.ExecuteTemplate(&buffer, "listOutages", userOutages); err != nil {
						log.Fatal(err)
					}
				}
			}
			msg.Text = buffer.String()
			msg.ParseMode = "MarkdownV2" //This parse mode enables format tags in TG
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
