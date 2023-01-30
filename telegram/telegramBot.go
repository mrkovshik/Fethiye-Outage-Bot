package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	)
type botStuff struct {
	updates tgbotapi.UpdatesChannel
	bot *tgbotapi.BotAPI
	chatID  int64
}

type neigborhoodMap map [string] string

var fethiyeMap = neigborhoodMap {
	"/Menteseoglu":  "menteseoglu" , 
	"/Babatasi": "babatasi",
	"/Ciftlik": "ciftlik_fethiye",
}
var bodrumMap = neigborhoodMap{
	"/Yokusbasi":  "yokusbasi",
	"/Tepecik": "tepecik",
	"/Ciftlik": "ciftlik_bodrum",
}
 

var cityMap = map [string] neigborhoodMap {
"/Fethiye": fethiyeMap,
"/Bodrum": bodrumMap,

}

func pickNeneigborhood (updates tgbotapi.UpdatesChannel,bot *tgbotapi.BotAPI, n map [string] neigborhoodMap){
	var msg tgbotapi.MessageConfig
	msgString:="Pick a neigborhood:\n "
	for i:= range n {
		msgString+=i+"\n" 
	}	

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgString+"If you want to get back to the main menu type /back")
	bot.Send(msg)	
	for update := range updates {
		if update.Message != nil { // If we got a message
		
			if _,ok:=cityMap[update.Message.Text]; ok{

			}

			}
 

		}
}



func pickCity (updates tgbotapi.UpdatesChannel,bot *tgbotapi.BotAPI){
	var msg tgbotapi.MessageConfig
	msgString:="Pick a city:\n "
	for i:= range cityMap{
		msgString+=i+"\n" 
	}	
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgString+"If you want to get back to the main menu type /back")
	bot.Send(msg)	
	for update := range updates {
		if update.Message != nil { // If we got a message
		
			if _,ok:=cityMap[update.Message.Text]; ok{
pickNeneigborhood(updates, bot,cityMap[update.Message.Text].neigborhoods )
			}

			}
 

		}
}


func dialogMain (b botStuff){
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(b.chatID, "Hi there, "+update.Message.From.FirstName+"! Welcome to Fethiye Outage Bot! You can /check out the outage schedule in your neigborhood and even /subscribe to get notifications about closest outages")
	for update := range b.updates {
	if update.Message != nil { // If we got a message
		
		for update := range b.updates {
			if update.Message != nil { // If we got a message
}

msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hi there, "+update.Message.From.FirstName+"! Welcome to Fethiye Outage Bot! Choose your city to get information about closest outages:\n /Bodrum\n/Fethiye\n/Marmaris\n/Dalaman\n/Milas\n/Mentese")
		b.bot.Send(msg)	

			continue		
		
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Come on, "+update.Message.From.FirstName+" ")
		}
		b.bot.Send(msg)
	}
	
	}
}
	

func main() {

  // подключаемся к боту с помощью токена
api:=os.Getenv("TELEGRAM_APITOKEN")

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
			botCred:= botStuff{ updates, bot, update.Message.Chat.ID}
			go dialogMain(botCred)
}
	}
}