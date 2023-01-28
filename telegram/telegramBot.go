package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TelegaBot() {
	fmt.Println("asdf",os.Getenv("TELEGRAM_APITOKEN"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))

	if err != nil {
		fmt.Println("botapi connection error")
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
	if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот в разработке")
				msg.ReplyToMessageID = update.Message.MessageID

	
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
