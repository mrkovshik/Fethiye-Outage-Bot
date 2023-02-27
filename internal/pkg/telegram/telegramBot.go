package telegram

import (
	"bytes"
	"os"
	"regexp"
	"text/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"
	"go.uber.org/zap"
)

type tool struct {
}

func newTool() *tool {
	return &tool{}
}

func (t *tool) formatDateAndMakeLocal(tm time.Time) string {
	return tm.Add(3 * time.Hour).String()[:19]
}

func (t *tool) escapeSimbols(s string) string {
	re := regexp.MustCompile(`[\\` + "`*_\\[\\]()#+\\-.!]")
	return re.ReplaceAllStringFunc(s, func(match string) string {
		return "\\" + match
	})
}

func BotRunner(ds *district.DistrictStore, store *postgres.OutageStore, logger *zap.Logger) {
	var err error
	tool := newTool()
	//mapping the functions for templates
	dialogTemplate := template.New("dialogTemplate").Funcs(template.FuncMap{
		"escape": tool.escapeSimbols,
		"format": tool.formatDateAndMakeLocal,
	})
	//parsing the template file
	t, err := dialogTemplate.ParseFiles("./templates/dialog_templates_eng.tpl")
	if err != nil {
		logger.Fatal("Parsing templates error",
			zap.Error(err),
		)
	}
	// reading the token from envirinment and connecting
	api := os.Getenv("OUTAGE_TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		logger.Fatal("Error connecting to telegram API",
			zap.Error(err),
		)
	}
	bot.Debug = true
	logger.Info("Authorized ",
		zap.Any(
			"account:", bot.Self.UserName),
	)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		var buffer bytes.Buffer
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if update.Message.Text == "/start" {
				if err := t.ExecuteTemplate(&buffer, "startMsg", nil); err != nil {
					logger.Fatal("Executing startMsg template error",
						zap.Error(err),
					)
				}
			} else {
				words, err := util.Normalize(update.Message.Text)
				if err != nil {
					logger.Fatal("",
						zap.Error(err),
					)
				}
				if len(words) == 0 {
					if err := t.ExecuteTemplate(&buffer, "badQuery", update.Message); err != nil {
						logger.Fatal("Error executing badQuery template",
							zap.Error(err),
						)
					}
				} else {
					guessDistr, err := ds.GetFuzzyMatch(words)
					if err != nil {
						logger.Fatal("",
							zap.Error(err),
						)
					}
					userOutages, err := store.GetActiveOutagesByCityDistrict(guessDistr.NameNormalized, guessDistr.CityNormalized)
					if err != nil {
						logger.Fatal("",
							zap.Error(err),
						)
					}
					if guessDistr.City == "no matches" {
						if err := t.ExecuteTemplate(&buffer, "badQuery", update.Message); err != nil {
							logger.Fatal("Error executing badQuery template",
								zap.Error(err),
							)
						}
					} else {
						if err := t.ExecuteTemplate(&buffer, "confirmDistr", guessDistr); err != nil {
							logger.Fatal("Error executing confirmDistr template",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, "listOutages", userOutages); err != nil {
							logger.Fatal("Error executing listOutages template",
								zap.Error(err),
							)
						}
					}
				}
			}
			msg.Text = buffer.String()
			msg.ParseMode = "MarkdownV2" //This parse mode enables format tags in TG
			if _, err := bot.Send(msg); err != nil {
				logger.Fatal("Error sending message",
					zap.Error(err),
				)
			}
		}
	}
}
