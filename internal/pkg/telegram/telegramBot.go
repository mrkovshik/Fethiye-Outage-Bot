package telegram

import (
	"bytes"
	"fmt"

	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/subscribtion"
	"github.com/robfig/cron"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"
	"go.uber.org/zap"
)
type stateMap map [int64] userState

type userState struct {
	previousKeyboard tgbotapi.ReplyKeyboardMarkup
	currentKeyboard  tgbotapi.ReplyKeyboardMarkup
	previousContext  string
	currentContext   string
	PickedCity       string
	PickedDistrict   string
	PickedPeriod     int
	lastActivityTime time.Time
}

func (u stateMap)cleanUpMap (c time.Duration){
	for id,state:= range u{
		if  time.Now().UTC().Sub(state.lastActivityTime) > c{
			delete(u,id)
			fmt.Printf("\nState deleted\n\n")
		} 
	}
}

func (u *userState) goBack() string {
	u.currentContext = u.previousContext
	u.currentKeyboard = u.previousKeyboard
	switch{
	case u.previousContext=="Pick sub city"||u.previousContext=="Pick check city"||u.previousContext=="Subscribtion settings":
		u.previousContext="menu"
		u.previousKeyboard=MenuKeyboard
		case u.previousContext=="Pick sub distr":
		u.previousContext="Pick sub city"
		u.previousKeyboard=CityKeyboard
	case u.previousContext=="Pick check distr":
		u.previousContext="Pick check city"
		u.previousKeyboard=CityKeyboard
	case u.previousContext=="Change location city":
		u.previousContext="Subscribtion settings"
		u.previousKeyboard=SettingsKeyboard
	}
		
	return "go_back"
}

func (u *userState) toMain() string {
	u.previousContext = "menu"
	u.currentContext = "menu"
	u.previousKeyboard = MenuKeyboard
	u.currentKeyboard = MenuKeyboard
	return "mainMenu_greet"
}

func formatDateAndMakeLocal(tm time.Time) string {
	return tm.Add(3 * time.Hour).String()[:19]
}

func escapeSimbols(s string) string {
	re := regexp.MustCompile(`[\\` + "`*_\\[\\]()#+\\-.!]")
	return re.ReplaceAllStringFunc(s, func(match string) string {
		return "\\" + match
	})
}

func (u *userState) processCity() string {

	switch u.PickedCity {
	case "Bodrum":
		u.currentKeyboard = BodrumKeyboard
	case "Dalaman":
		u.currentKeyboard = DalamanKeyboard
	case "Datça":
		u.currentKeyboard = DatcaKeyboard
	case "Fethiye":
		u.currentKeyboard = FethiyeKeyboard
	case "Kavaklıdere":
		u.currentKeyboard = KavaklidereKeyboard
	case "Köyceğiz":
		u.currentKeyboard = KoycegizKeyboard
	case "Marmaris":
		u.currentKeyboard = MarmarisKeyboard
	case "Menteşe":
		u.currentKeyboard = MenteseKeyboard
	case "Milas":
		u.currentKeyboard = MilasKeyboard
	case "Ortaca":
		u.currentKeyboard = OrtacaKeyboard
	case "Seydikemer":
		u.currentKeyboard = SeydikemerKeyboard
	case "Ula":
		u.currentKeyboard = UlaKeyboard
	case "Yatağan":
		u.currentKeyboard = YataganKeyboard
	case "GO BACK":
		return u.goBack()
	default:
		return "claim_buttons"
	}
	u.previousContext = u.currentContext
	u.previousKeyboard = CityKeyboard
	switch u.currentContext {
	case "Pick sub city":
		u.currentContext = "Pick sub distr"
	case "Pick check city":
		u.currentContext = "Pick check distr"
	case "Change location city":
		u.currentContext = "Change location distr"
	}
	return "pickCity_confirm"
}

func BotRunner(ds *district.DistrictStore, store *postgres.OutageStore, subStore *subscribtion.SubscribtionStore, logger *zap.Logger, cfg config.Config) {
	var err error

	userMap := stateMap{}
	c := cron.New()
	err = c.AddFunc(cfg.SchedulerConfig.StateCleanUpPeriod, func() { userMap.cleanUpMap(time.Duration(cfg.BotConfig.UserStateLifeTime)*time.Minute) })
	if err != nil {
		logger.Fatal("Sceduler error",
			zap.Error(err),
		)
	}
	go c.Start()
	//mapping the functions for templates
	dialogTemplate := template.New("dialogTemplate").Funcs(template.FuncMap{
		"escape": escapeSimbols,
		"format": formatDateAndMakeLocal,
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
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, ok := userMap[update.Message.Chat.ID]; !ok {
				userMap[update.Message.Chat.ID] = userState{
					currentContext:   "start",
					previousContext: "start",
					lastActivityTime: time.Now().UTC(),
					currentKeyboard: MenuKeyboard,
					previousKeyboard: MenuKeyboard,
				}
			}			
			currentUserState := userMap[update.Message.Chat.ID]
			currentUserState.lastActivityTime=time.Now().UTC()
			if update.Message.Text == "/start" || update.Message.Text == "/main_menu" {
				if err := t.ExecuteTemplate(&buffer, "name_greet", update.Message.From); err != nil {
					logger.Fatal("Executing startMsg template error",
						zap.Error(err),
					)
				}
				if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
					logger.Fatal("Executing startMsg template error",
						zap.Error(err),
					)
				}

			} else {
				switch {
				case currentUserState.currentContext == "menu":
					switch update.Message.Text {
					case "Subscribe for alerts":
						isExist, err := subStore.SubExists(update.Message.Chat.ID)
						if err != nil {
							logger.Fatal("",
								zap.Error(err),
							)
						}
						if isExist {
							if err := t.ExecuteTemplate(&buffer, "have_sub", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						} else {
							if err := t.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							currentUserState.currentContext = "Pick sub city"
							currentUserState.currentKeyboard = CityKeyboard
						}
					case "Check out for outages":
						if err := t.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						currentUserState.currentContext = "Pick check city"
						currentUserState.currentKeyboard = CityKeyboard
					case "Subscribtion settings":
						isExist, err := subStore.SubExists(update.Message.Chat.ID)
						if err != nil {
							logger.Fatal("",
								zap.Error(err),
							)
						}
						if !isExist {
							if err := t.ExecuteTemplate(&buffer, "no_subs", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						} else {
							if err := t.ExecuteTemplate(&buffer, "settings_greet", update.Message.From); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							currentUserState.currentContext = "Subscribtion settings"
							currentUserState.currentKeyboard = SettingsKeyboard
						}
					default:
						if err := t.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
					}

				case currentUserState.currentContext == "Pick sub city" || currentUserState.currentContext == "Pick check city" || currentUserState.currentContext == "Change location city":
					currentUserState.PickedCity = update.Message.Text
					tmp := currentUserState.processCity()
					if err := t.ExecuteTemplate(&buffer, tmp, currentUserState); err != nil {
						logger.Fatal("Executing message template error",
							zap.Error(err),
						)
					}

				case currentUserState.currentContext == "Pick sub distr" || currentUserState.currentContext == "Pick check distr" || currentUserState.currentContext == "Change location distr":
					currentUserState.PickedDistrict = update.Message.Text
					//Checking out if input is valid by searching a match from DB
					match, err := ds.GetNormFromDB(currentUserState.PickedCity, currentUserState.PickedDistrict)
					if err != nil {
						logger.Fatal("",
							zap.Error(err),
						)
					}
					if update.Message.Text == "GO BACK" {
						if err := t.ExecuteTemplate(&buffer, currentUserState.goBack(), currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
					} else {
						//if input is valid
						if match.Name != "" {
							if err := t.ExecuteTemplate(&buffer, "pickDistr_confirm", currentUserState); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							switch currentUserState.currentContext {
							case "Pick sub distr":
								currentUserState.previousContext = "Pick sub distr"
								currentUserState.currentContext = "Pick sub alert period"
								currentUserState.previousKeyboard = currentUserState.currentKeyboard
								currentUserState.currentKeyboard = HoursKeyboard
								if err := t.ExecuteTemplate(&buffer, "pickPeriod_greet", currentUserState); err != nil {
									logger.Fatal("Executing message template error",
										zap.Error(err),
									)
								}
							case "Pick check distr":
								userOutages, err := store.GetActiveOutagesByCityDistrict(match.NameNormalized, match.CityNormalized)
								if err != nil {
									logger.Fatal("",
										zap.Error(err),
									)
								}
								if err := t.ExecuteTemplate(&buffer, "listOutages", userOutages); err != nil {
									logger.Fatal("Error executing listOutages template",
										zap.Error(err),
									)
								}
								if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), update.Message.From); err != nil {
									logger.Fatal("Executing startMsg template error",
										zap.Error(err),
									)
								}
							case "Change location distr":
								s := subscribtion.Subscribtion{
									City:     currentUserState.PickedCity,
									District: currentUserState.PickedDistrict,
									ChatID:   update.Message.Chat.ID,
								}
								err := subStore.ModifyLocation(s)
								if err != nil {
									if err := t.ExecuteTemplate(&buffer, "error", nil); err != nil {
										logger.Fatal("Executing message template error",
											zap.Error(err),
										)
									}
									logger.Warn("",
										zap.Error(err),
									)
								} else {
									if err := t.ExecuteTemplate(&buffer, "change_location_confirm", nil); err != nil {
										logger.Fatal("Executing message template error",
											zap.Error(err),
										)
									}
									if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), update.Message.From); err != nil {
										logger.Fatal("Executing startMsg template error",
											zap.Error(err),
										)
									}
								}
							}
						} else {
							if err := t.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						}
					}

				case currentUserState.currentContext == "Pick sub alert period":
					switch {
					case update.Message.Text == "GO BACK":
						if err := t.ExecuteTemplate(&buffer, currentUserState.goBack(), currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, "pickCity_confirm", currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						

					case update.Message.Text == "2 hours" || update.Message.Text == "6 hours" || update.Message.Text == "12 hours" || update.Message.Text == "24 hours":
						h := strings.TrimSuffix(update.Message.Text, " hours")
						currentUserState.PickedPeriod, err = strconv.Atoi(h)
						if err != nil {
							logger.Fatal("Converting period error",
								zap.Error(err),
							)
						}
						s := subscribtion.Subscribtion{
							City:     currentUserState.PickedCity,
							District: currentUserState.PickedDistrict,
							Period:   currentUserState.PickedPeriod,
							ChatID:   update.Message.Chat.ID,
						}
						if err := subStore.Save(s); err != nil {
							if err := t.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							logger.Warn("",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, "set_period_confirm", currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
					default:
						if err := t.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}

					}
				case currentUserState.currentContext == "Change alert period":
					switch {
					case update.Message.Text == "GO BACK":
						if err := t.ExecuteTemplate(&buffer, currentUserState.goBack(), currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						if err := t.ExecuteTemplate(&buffer, "settings_greet",nil ); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}

					case update.Message.Text == "2 hours" || update.Message.Text == "6 hours" || update.Message.Text == "12 hours" || update.Message.Text == "24 hours":
						h := strings.TrimSuffix(update.Message.Text, " hours")
						currentUserState.PickedPeriod, err = strconv.Atoi(h)
						if err != nil {
							logger.Fatal("Converting period error",
								zap.Error(err),
							)
						}
						s := subscribtion.Subscribtion{
							Period: currentUserState.PickedPeriod,
							ChatID: update.Message.Chat.ID,
						}

						err := subStore.ModifyPeriod(s)
						if err != nil {
							if err := t.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							logger.Warn("",
								zap.Error(err),
							)
						} else {
							if err := t.ExecuteTemplate(&buffer, "change_period_confirm", currentUserState); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						}
					}

				case currentUserState.currentContext == "Subscribtion settings":
					switch update.Message.Text {
					case "Change location":
						if err := t.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						currentUserState.currentContext = "Change location city"
						currentUserState.currentKeyboard = CityKeyboard
						currentUserState.previousContext = "Subscribtion settings"
						currentUserState.previousKeyboard = SettingsKeyboard
					case "Cancel subscribtion":
						s := subscribtion.Subscribtion{
							ChatID: update.Message.Chat.ID,
						}
						if err := subStore.RemoveSubscribtion(s); err != nil {
							if err := t.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							logger.Warn("",
								zap.Error(err),
							)
						} else {
							if err := t.ExecuteTemplate(&buffer, "cancel_confirm", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						}
					case "Change alert period":
						if err := t.ExecuteTemplate(&buffer, "change_period_greet", nil); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						currentUserState.currentContext = "Change alert period"
						currentUserState.currentKeyboard = HoursKeyboard
						currentUserState.previousContext = "Subscribtion settings"
						currentUserState.previousKeyboard = SettingsKeyboard

					case "View current subscribtion":
						subs, err := subStore.GetSubsByChatID(update.Message.Chat.ID)
						if err != nil {
							if err := t.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
							logger.Warn("",
								zap.Error(err),
							)
						} else {

							if err := t.ExecuteTemplate(&buffer, "show_sub", subs[0]); err != nil {
								logger.Fatal("Executing message template error",
									zap.Error(err),
								)
							}
						}
						if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}

					case "GO BACK":
						if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
					default:
						if err := t.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
					}
				default:
					if err := t.ExecuteTemplate(&buffer, "press_start", nil); err != nil {
						logger.Fatal("Executing message template error",
							zap.Error(err),
						)
					}
					if err := t.ExecuteTemplate(&buffer, currentUserState.toMain(), nil); err != nil {
						logger.Fatal("Executing message template error",
							zap.Error(err),
						)
					}
				}
			}
			userMap[update.Message.Chat.ID] = currentUserState
			msg.ReplyMarkup = currentUserState.currentKeyboard

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
