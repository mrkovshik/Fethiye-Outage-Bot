package telegram

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
)


func BotRunner (ds *district.DistrictStore, muskiStore *postgres.OutageStore) {
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
		   msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		   if update.Message.Text=="/start"{
			   msg.Text = "Please print your city and neigbourhood divided by space, for example 'Fethie Taşyaka'"
	   } else {
		   guessDistr, err:=ds.GetFuzzyMatch(update.Message.Text)
		   if err != nil {
			   fmt.Println("Fuzzy search error", err)
		   }
		   userOutages, err:=muskiStore.GetActiveOutagesByCityDistrict(guessDistr.Name, guessDistr.City)
		   if err != nil {
			   fmt.Println("Outages search error", err)
		   }
		   if guessDistr.City=="no matches" {
			
			msg.Text="I am sorry, but I can't find anythithg like '" + update.Message.Text + "' Maybe we should try again?\n\n" +"Please print your city and neigbourhood divided by space, for example 'Fethie Taşyaka'"
		} else {
		   msg.Text="Did you mean '" + guessDistr.City +" "+ guessDistr.Name + "'?\n\n"
		   if len(userOutages)==0 {
			   msg.Text+= "There is no outages planned in your neigborhood in the closest time"
		   } else {
			   msg.Text+= "Here are the closest outages found for your neigborhood:\n\n"
			   for _,i:=range userOutages{
				   msg.Text+= i.Resource +" outage from " + i.StartDate.Add(3*time.Hour).String()[:19] + " to " + i.EndDate.Add(3*time.Hour).String()[:19] + "\n\n"
			   }
		   }
		}
	   }
   
	   if _, err := bot.Send(msg); err != nil {
		   log.Panic(err)
	   }
}
   }
}