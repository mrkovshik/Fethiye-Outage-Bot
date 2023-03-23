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
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/alert"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/subscribtion"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"
	"github.com/robfig/cron"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"
	"go.uber.org/zap"
)

type stateMap map[int64]userState

type userState struct {
	previousKeyboard tgbotapi.ReplyKeyboardMarkup
	currentKeyboard  tgbotapi.ReplyKeyboardMarkup
	previousContext  string
	currentContext   string
	PickedCity       string
	PickedDistrict   string
	PickedPeriod     int
	lastActivityTime time.Time
	CurrentLanguage  LanguagePack
}

func (u stateMap) cleanUpMap(c time.Duration) {
	for id, state := range u {
		if time.Now().UTC().Sub(state.lastActivityTime) > c {
			delete(u, id)
			fmt.Printf("\nState deleted\n\n")
		}
	}
}

func (u *userState) goBack(l LanguagePack) string {
	u.currentContext = u.previousContext
	u.currentKeyboard = u.previousKeyboard
	switch {
	case u.previousContext == "Pick sub city" || u.previousContext == "Pick check city" || u.previousContext == "Subscribtion settings":
		u.previousContext = "menu"
		u.previousKeyboard = l.MenuKeyboard
	case u.previousContext == "Pick sub distr":
		u.previousContext = "Pick sub city"
		u.previousKeyboard = l.CityKeyboard
	case u.previousContext == "Pick check distr":
		u.previousContext = "Pick check city"
		u.previousKeyboard = l.CityKeyboard
	case u.previousContext == "Change location city":
		u.previousContext = "Subscribtion settings"
		u.previousKeyboard = l.SettingsKeyboard
	}

	return "go_back"
}

func (u *userState) toMain(l LanguagePack) string {
	u.previousContext = "menu"
	u.currentContext = "menu"
	u.previousKeyboard = l.MenuKeyboard
	u.currentKeyboard = l.MenuKeyboard
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

func sendAlert(a alert.AlertStore, o postgres.OutageStore, bot *tgbotapi.BotAPI, t *template.Template, logger *zap.Logger, l LanguagePack) {
	var buffer bytes.Buffer
	alerts, err := a.GetActiveAlerts()
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	for _, al := range alerts {
		TheOutage, err := o.GetOutageByID(al.OutageID)
		if err != nil {
			logger.Fatal("", zap.Error(err))
		}
		if err := l.Template.ExecuteTemplate(&buffer, "alert", TheOutage); err != nil {
			if err != nil {
				logger.Fatal("Executing message template error", zap.Error(err))
			}
		}
		msg := tgbotapi.NewMessage(al.ChatID, "")
		msg.Text = buffer.String()
		msg.ParseMode = "MarkdownV2" //This parse mode enables format tags in TG
		if _, err := bot.Send(msg); err != nil {
			if err != nil {
				logger.Fatal("Sending message error", zap.Error(err))
			}
		}
		err = a.SetIsSent(al.ID)
		if err != nil {
			logger.Fatal("", zap.Error(err))
		}
	}
}

func (u *userState) processCity(l LanguagePack) string {

	switch u.PickedCity {
	case "Bodrum":
		u.currentKeyboard = l.BodrumKeyboard
	case "Dalaman":
		u.currentKeyboard = l.DalamanKeyboard
	case "Datça":
		u.currentKeyboard = l.DatcaKeyboard
	case "Fethiye":
		u.currentKeyboard = l.FethiyeKeyboard
	case "Kavaklıdere":
		u.currentKeyboard = l.KavaklidereKeyboard
	case "Köyceğiz":
		u.currentKeyboard = l.KoycegizKeyboard
	case "Marmaris":
		u.currentKeyboard = l.MarmarisKeyboard
	case "Menteşe":
		u.currentKeyboard = l.MenteseKeyboard
	case "Milas":
		u.currentKeyboard = l.MilasKeyboard
	case "Ortaca":
		u.currentKeyboard = l.OrtacaKeyboard
	case "Seydikemer":
		u.currentKeyboard = l.SeydikemerKeyboard
	case "Ula":
		u.currentKeyboard = l.UlaKeyboard
	case "Yatağan":
		u.currentKeyboard = l.YataganKeyboard
	case "GO BACK":
		return u.goBack(l)
	case "НАЗАД":
		return u.goBack(l)
	case "GERI":
		return u.goBack(l)
	default:
		return "claim_buttons"
	}
	u.previousContext = u.currentContext
	u.previousKeyboard = l.CityKeyboard
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

func BotRunner(ds *district.DistrictStore, store *postgres.OutageStore, subStore *subscribtion.SubscribtionStore, alertStore *alert.AlertStore, logger *zap.Logger, cfg config.Config) {
	var err error
	userMap := stateMap{}
	Langs, err := NewLanguages()
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	// CurrentLanguage := Langs.Eng
	// reading the token from envirinment and connecting
	api := os.Getenv("OUTAGE_TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		logger.Fatal("Error connecting to telegram API", zap.Error(err))
	}
	c := cron.New()
	err = c.AddFunc(cfg.SchedulerConfig.StateCleanUpPeriod, func() { userMap.cleanUpMap(time.Duration(cfg.BotConfig.UserStateLifeTime) * time.Hour) })
	if err != nil {
		logger.Fatal("Sceduler error", zap.Error(err))
	}
	err = c.AddFunc(cfg.SchedulerConfig.AlertSendPeriod, func() { sendAlert(*alertStore, *store, bot, Langs.Eng.Template, logger, Langs.Eng) })
	if err != nil {
		logger.Fatal("Sceduler error", zap.Error(err))
	}
	go c.Start()
	sendAlert(*alertStore, *store, bot, Langs.Eng.Template, logger, Langs.Eng)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	bot.Debug = true
	logger.Info("Authorized ", zap.Any("account:", bot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		var buffer bytes.Buffer
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			//Checking out if user has a state in stateMap
			if _, ok := userMap[update.Message.Chat.ID]; !ok {
				userMap[update.Message.Chat.ID] = userState{
					currentContext:   "start",
					previousContext:  "start",
					lastActivityTime: time.Now().UTC(),
					currentKeyboard:  Langs.Eng.MenuKeyboard,
					previousKeyboard: Langs.Eng.MenuKeyboard,
					CurrentLanguage:  Langs.Eng,
				}
			}
			currentUserState := userMap[update.Message.Chat.ID]
			currentUserState.lastActivityTime = time.Now().UTC()
			if update.Message.Text == "/start" || update.Message.Text == "/main_menu" {
				if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "name_greet", update.Message.From); err != nil {
					logger.Fatal("Executing startMsg template error", zap.Error(err))
				}
				if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
					logger.Fatal("Executing startMsg template error", zap.Error(err))
				}
			} else {
				switch {
				case currentUserState.currentContext == "menu":
					switch {
					case update.Message.Text == "Subscribe for alerts" || update.Message.Text == "Подписаться на оповещения" || update.Message.Text == "Uyarılara abone ol":
						isExist, err := subStore.SubExists(update.Message.Chat.ID)
						if err != nil {
							logger.Fatal("", zap.Error(err))
						}
						if isExist {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "have_sub", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						} else {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							currentUserState.currentContext = "Pick sub city"
							currentUserState.currentKeyboard = currentUserState.CurrentLanguage.CityKeyboard
						}
					case update.Message.Text == "Check out for outages" || update.Message.Text == "Проверить отключения" || update.Message.Text == "Kapatmaları kontrol et":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						currentUserState.currentContext = "Pick check city"
						currentUserState.currentKeyboard = currentUserState.CurrentLanguage.CityKeyboard
					case update.Message.Text == "Subscription settings" || update.Message.Text == "Настройки оповещений" || update.Message.Text == "Uyarı ayarları":
						isExist, err := subStore.SubExists(update.Message.Chat.ID)
						if err != nil {
							logger.Fatal("", zap.Error(err))
						}
						if !isExist {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "no_subs", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						} else {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "settings_greet", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							currentUserState.currentContext = "Subscribtion settings"
							currentUserState.currentKeyboard = currentUserState.CurrentLanguage.SettingsKeyboard
						}
					case update.Message.Text == "Рус/Eng/Tür":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_language", nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						currentUserState.currentContext = "Language change"
						currentUserState.currentKeyboard = currentUserState.CurrentLanguage.LanguageKeyboard

					default:
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					}
				case currentUserState.currentContext == "Language change":
					if update.Message.Text == "НАЗАД" || update.Message.Text == "GO BACK" || update.Message.Text == "GERI" {
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					} else {
						switch update.Message.Text {
						case "English":
							currentUserState.CurrentLanguage = Langs.Eng
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_language_confirm", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						case "Русский":
							currentUserState.CurrentLanguage = Langs.Rus
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_language_confirm", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						case "Türkçe":
							currentUserState.CurrentLanguage = Langs.Tur
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_language_confirm", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						default:
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						}
					}

				case currentUserState.currentContext == "Pick sub city" || currentUserState.currentContext == "Pick check city" || currentUserState.currentContext == "Change location city":
					currentUserState.PickedCity = update.Message.Text
					tmp := currentUserState.processCity(currentUserState.CurrentLanguage)
					if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, tmp, currentUserState); err != nil {
						logger.Fatal("Executing message template error",
							zap.Error(err),
						)
					}

				case currentUserState.currentContext == "Pick sub distr" || currentUserState.currentContext == "Pick check distr" || currentUserState.currentContext == "Change location distr":
					currentUserState.PickedDistrict = update.Message.Text
					//Checking out if input is valid by searching a match from DB
					match, err := ds.GetNormFromDB(currentUserState.PickedCity, currentUserState.PickedDistrict)
					if err != nil {
						logger.Fatal("", zap.Error(err))
					}
					if update.Message.Text == "НАЗАД" || update.Message.Text == "GO BACK" || update.Message.Text == "GERI" {
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.goBack(currentUserState.CurrentLanguage), currentUserState); err != nil {
							logger.Fatal("Executing message template error",
								zap.Error(err),
							)
						}
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					} else {
						//if input is valid
						if match.Name != "" {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickDistr_confirm", currentUserState); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							switch currentUserState.currentContext {
							case "Pick sub distr":
								currentUserState.previousContext = "Pick sub distr"
								currentUserState.currentContext = "Pick sub alert period"
								currentUserState.previousKeyboard = currentUserState.currentKeyboard
								currentUserState.currentKeyboard = currentUserState.CurrentLanguage.HoursKeyboard
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickPeriod_greet", currentUserState); err != nil {
									logger.Fatal("Executing message template error", zap.Error(err))
								}
							case "Pick check distr":
								userOutages, err := store.GetActiveOutagesByCityDistrict(match.NameNormalized, match.CityNormalized)
								if err != nil {
									logger.Fatal("", zap.Error(err))
								}
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "listOutages", userOutages); err != nil {
									logger.Fatal("Error executing listOutages template", zap.Error(err))
								}
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), update.Message.From); err != nil {
									logger.Fatal("Executing startMsg template error", zap.Error(err))
								}
							case "Change location distr":
								c, err := util.Normalize(currentUserState.PickedCity)
								if err != nil {
									logger.Fatal("", zap.Error(err))
								}
								d, err := util.Normalize(currentUserState.PickedDistrict)
								if err != nil {
									logger.Fatal("", zap.Error(err))
								}
								s := subscribtion.Subscribtion{
									City:               currentUserState.PickedCity,
									District:           currentUserState.PickedDistrict,
									Period:             currentUserState.PickedPeriod,
									ChatID:             update.Message.Chat.ID,
									CityNormalized:     strings.Join(c, " "),
									DistrictNormalized: strings.Join(d, " "),
								}
								err = subStore.ModifyLocation(s)
								if err != nil {
									if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
										logger.Fatal("Executing message template error", zap.Error(err))
									}
									logger.Warn("", zap.Error(err))
								} else {
									if err := alertStore.CancelByChatID(update.Message.Chat.ID); err != nil {
										logger.Fatal("", zap.Error(err))
									}
									newSub, err := subStore.GetSubsByChatID(update.Message.Chat.ID)
									if err != nil {
										logger.Fatal("", zap.Error(err))
									}
									if err := alertStore.GenerateAlertsForNewSub(*store, newSub[0]); err != nil {
										logger.Fatal("", zap.Error(err))
									}
									if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_location_confirm", nil); err != nil {
										logger.Fatal("Executing message template error", zap.Error(err))
									}
									if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), update.Message.From); err != nil {
										logger.Fatal("Executing startMsg template error", zap.Error(err))
									}
								}
							}
						} else {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						}
					}

				case currentUserState.currentContext == "Pick sub alert period":
					switch {
					case update.Message.Text == "НАЗАД" || update.Message.Text == "GO BACK" || update.Message.Text == "GERI":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.goBack(currentUserState.CurrentLanguage), currentUserState); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickCity_confirm", currentUserState); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}

					case update.Message.Text == "2 hours" || update.Message.Text == "6 hours" || update.Message.Text == "12 hours" || update.Message.Text == "24 hours" || update.Message.Text == "2 часа" || update.Message.Text == "6 часов" || update.Message.Text == "12 часов" || update.Message.Text == "24 часа" || update.Message.Text == "2 saat" || update.Message.Text == "6 saat" || update.Message.Text == "12 saat" || update.Message.Text == "24 saat":
						re := regexp.MustCompile("[0-9]+")
						h := re.FindString(update.Message.Text)
						currentUserState.PickedPeriod, err = strconv.Atoi(h)
						if err != nil {
							logger.Fatal("Converting period error", zap.Error(err))
						}
						c, err := util.Normalize(currentUserState.PickedCity)
						if err != nil {
							logger.Fatal("", zap.Error(err))
						}
						d, err := util.Normalize(currentUserState.PickedDistrict)
						if err != nil {
							logger.Fatal("", zap.Error(err))
						}
						s := subscribtion.Subscribtion{
							City:               currentUserState.PickedCity,
							District:           currentUserState.PickedDistrict,
							Period:             currentUserState.PickedPeriod,
							ChatID:             update.Message.Chat.ID,
							CityNormalized:     strings.Join(c, " "),
							DistrictNormalized: strings.Join(d, " "),
						}

						if err := subStore.Save(s); err != nil {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							logger.Warn("", zap.Error(err))
						}
						newSub, err := subStore.GetSubsByChatID(update.Message.Chat.ID)
						if err != nil {
							logger.Fatal("", zap.Error(err))
						}
						if err := alertStore.GenerateAlertsForNewSub(*store, newSub[0]); err != nil {
							logger.Fatal("", zap.Error(err))
						}

						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "set_period_confirm", currentUserState); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					default:
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}

					}
				case currentUserState.currentContext == "Change alert period":
					switch {
					case update.Message.Text == "НАЗАД" || update.Message.Text == "GO BACK" || update.Message.Text == "GERI":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.goBack(currentUserState.CurrentLanguage), currentUserState); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "settings_greet", nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}

					case update.Message.Text == "2 hours" || update.Message.Text == "6 hours" || update.Message.Text == "12 hours" || update.Message.Text == "24 hours" || update.Message.Text == "2 часа" || update.Message.Text == "6 часов" || update.Message.Text == "12 часов" || update.Message.Text == "24 часа" || update.Message.Text == "2 saat" || update.Message.Text == "6 saat" || update.Message.Text == "12 saat" || update.Message.Text == "24 saat":
						re := regexp.MustCompile("[0-9]+")
						h := re.FindString(update.Message.Text)
						currentUserState.PickedPeriod, err = strconv.Atoi(h)
						if err != nil {
							logger.Fatal("Converting period error", zap.Error(err))
						}
						s := subscribtion.Subscribtion{
							Period: currentUserState.PickedPeriod,
							ChatID: update.Message.Chat.ID,
						}

						err := subStore.ModifyPeriod(s)
						if err != nil {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							logger.Warn("", zap.Error(err))
						} else {
							if err := alertStore.CancelByChatID(update.Message.Chat.ID); err != nil {
								logger.Fatal("", zap.Error(err))
							}
							newSub, err := subStore.GetSubsByChatID(update.Message.Chat.ID)
							if err != nil {
								logger.Fatal("", zap.Error(err))
							}
							if err := alertStore.GenerateAlertsForNewSub(*store, newSub[0]); err != nil {
								logger.Fatal("", zap.Error(err))
							}
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_period_confirm", currentUserState); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						}
					}

				case currentUserState.currentContext == "Subscribtion settings":
					switch {
					case update.Message.Text == "Change location"||update.Message.Text == "Изменить локацию"||update.Message.Text == "Yer değiştir":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "pickCity_greet", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						currentUserState.currentContext = "Change location city"
						currentUserState.currentKeyboard = currentUserState.CurrentLanguage.CityKeyboard
						currentUserState.previousContext = "Subscribtion settings"
						currentUserState.previousKeyboard = currentUserState.CurrentLanguage.SettingsKeyboard
					case update.Message.Text == "Cancel subscription"||update.Message.Text == "Отменить подписку"||update.Message.Text == "Aboneliği iptal et":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "cancel_you_sure", nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))

						}
						currentUserState.currentContext = "Cancel subscribtion confirmed"
						currentUserState.currentKeyboard = currentUserState.CurrentLanguage.ConfirmKeyboard
						currentUserState.previousContext = "Subscribtion settings"
						currentUserState.previousKeyboard = currentUserState.CurrentLanguage.SettingsKeyboard

					case update.Message.Text == "Change alert period"||update.Message.Text == "Изменить время оповещения"||update.Message.Text == "Uyarı saatini değiştir":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "change_period_greet", nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
						currentUserState.currentContext = "Change alert period"
						currentUserState.currentKeyboard = currentUserState.CurrentLanguage.HoursKeyboard
						currentUserState.previousContext = "Subscribtion settings"
						currentUserState.previousKeyboard = currentUserState.CurrentLanguage.SettingsKeyboard

					case update.Message.Text == "View current subscription"||update.Message.Text == "Посмотреть текущую подписку"||update.Message.Text == "Geçerli aboneliği görüntüle":
						subs, err := subStore.GetSubsByChatID(update.Message.Chat.ID)
						if err != nil {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							logger.Warn("", zap.Error(err))
						} else {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "show_sub", subs[0]); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
						}
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					case update.Message.Text == "НАЗАД" || update.Message.Text == "GO BACK" || update.Message.Text == "GERI":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					default:
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					}
				case currentUserState.currentContext == "Cancel subscribtion confirmed":
					switch  {
					case update.Message.Text=="Yes, cancel it"||update.Message.Text=="Да, отменить"||update.Message.Text=="Evet İptal":
						s := subscribtion.Subscribtion{
							ChatID: update.Message.Chat.ID,
						}
						if err := alertStore.CancelByChatID(update.Message.Chat.ID); err != nil {
							if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
								logger.Fatal("Executing message template error", zap.Error(err))
							}
							logger.Warn("", zap.Error(err))
						} else {
							if err := subStore.CancelSubscribtion(s); err != nil {
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "error", nil); err != nil {
									logger.Fatal("Executing message template error", zap.Error(err))
								}
								logger.Warn("", zap.Error(err))
							} else {
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "cancel_confirm", nil); err != nil {
									logger.Fatal("Executing message template error", zap.Error(err))
								}
								if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
									logger.Fatal("Executing message template error", zap.Error(err))
								}
							}
						}
					case update.Message.Text=="No, let's go back"||update.Message.Text=="Нет, вернуться"||update.Message.Text=="Hayır, geri gitmek":
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.goBack(currentUserState.CurrentLanguage), nil); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					default:
						if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "claim_buttons", update.Message.From); err != nil {
							logger.Fatal("Executing message template error", zap.Error(err))
						}
					}

				default:
					if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, "press_start", nil); err != nil {
						logger.Fatal("Executing message template error", zap.Error(err))
					}
					if err := currentUserState.CurrentLanguage.Template.ExecuteTemplate(&buffer, currentUserState.toMain(currentUserState.CurrentLanguage), nil); err != nil {
						logger.Fatal("Executing message template error", zap.Error(err))
					}
				}
			}
			userMap[update.Message.Chat.ID] = currentUserState
			msg.ReplyMarkup = currentUserState.currentKeyboard
			msg.Text = buffer.String()
			msg.ParseMode = "MarkdownV2" //This parse mode enables format tags in TG
			if _, err := bot.Send(msg); err != nil {
				logger.Fatal("Error sending message", zap.Error(err))
			}
		}
	}
}
