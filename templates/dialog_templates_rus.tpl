{{define "listOutages"}}
{{if eq (len .) 0 }}
    *В вашем районе нет запланированных отключений на ближайшие даты*
{{else}}
*В вашем районе запланированы следующие отключения:*
{{range .}}
{{if eq .Resource "power" }}
*Отключение электроэнергии* {{else}} *Отключение водоснабжения* {{end}} с {{escape (format (.StartDate))}} по {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}
*На следующих улицах и локациях:*
{{escape (.Notes)}}{{end}}
{{end}}
{{end}}
{{end}}

{{define "alert"}}
Информируем Вас о том, что в районе *{{.District}} {{.City}}* запланировано {{if eq .Resource "power" }}*Отключение электроэнергии*{{else}}*Отключение водоснабжения* {{end}} сроком с {{escape (format (.StartDate))}} по {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}
*На следующих улицах и локациях:*
{{escape (.Notes)}}{{end}}
{{end}}


{{define "name_greet"}}
Привет {{escape (.UserName)}}{{escape ("!")}}
{{end}}

{{define "mainMenu_greet"}}
Что бы вы хотели сделать?
{{end}}

{{define "pickCity_greet"}}
Пожалуйста, выберите город из списка ниже:
{{end}}

{{define "settings_greet"}}
Здесь вы можете изменять параметры своей подписки{{escape (".")}} Что бы вы хотели именить?
{{end}}

{{define "claim_buttons"}}
Пожалуйста, нажмите на кнопку из меню снизу, *не печатайте свой ответ*:
{{end}}

{{define "pickCity_confirm"}}
Вы выбрали город *{{.PickedCity}}*{{escape ("!")}} Теперь выберите ваш район из списка ниже:
{{end}}

{{define "pickPeriod_greet"}}
Теперь выберите, за сколько часов вас нужно предупредить об отключениях:
{{end}}

{{define "change_period_greet"}}
Пожалуйста выберите, за сколько часов вас нужно предупредить об отключениях:
{{end}}

{{define "pickDistr_confirm"}}
Вы выбрали район *{{.PickedDistrict}}* в городе *{{.PickedCity}}*{{escape ("!")}}
{{end}}

{{define "change_location_confirm"}}
Ваша подписка была успешно обновлена{{escape ("!")}}
{{end}}

{{define "change_period_confirm"}}
Ваша подписка была успешно обновлена{{escape ("!")}} С этого момента вы будете получать уведомления об отключениях в вашем районе за *{{.PickedPeriod}}* {{if eq .PickedPeriod 2 }} часа {{else}}часов{{end}} до их начала{{escape (".")}}
{{end}}

{{define "set_period_confirm"}}
Ваша подписка успешно оформлена{{escape ("!")}} С этого момента вы будете получать уведомления об отключениях в *{{.PickedDistrict}}* *{{.PickedCity}}* за *{{.PickedPeriod}}* {{if eq .PickedPeriod 2 }} часа {{else}}часов{{end}} до их начала{{escape (".")}}
{{end}}

{{define "show_sub"}}
Вы подписаны на получение уведомлений в районе *{{.District}}* города *{{.City}}* за *{{.Period}}* {{if or (eq .Period 2) (eq .Period 24)}} часа {{else}}часов{{end}} до их начала{{escape (".")}}
{{end}}

{{define "no_subs"}}
Похоже{{escape (",")}} что у вас еще нет подписки на оповещения{{escape (".")}} Выберите пункт *Подписаться на оповещения*{{escape (",")}} чтобы оформить подписку{{escape (".")}}
{{end}}

{{define "have_sub"}}
Похоже{{escape (",")}} что у вас уже есть подписка на оповещения{{escape (".")}} Выберите пункт *Настройки оповещений*{{escape (",")}} если хотите отменить или изменить параметры подписки{{escape (".")}}
{{end}}

{{define "cancel_confirm"}}
Ваша подписка была успешно отменена{{escape ("!")}}
{{end}}

{{define "cancel_you_sure"}}
Вы уверены{{escape (",")}} что хотите отменить подписку{{escape ("?")}}
{{end}}

{{define "go_back"}}
Возвращаемся в предыдущее меню
{{end}}

{{define "press_start"}}
{{escape ("Прошу прощения, Вы давненько не заходили, и я забыл, на чем мы остановились в прошлый раз(")}} 
{{escape ("Давайте попробуем начать с самого начала?")}}
{{end}}

{{define "change_language"}}
{{escape ("Please, pick your language")}}

{{escape ("Пожалуйста, выберите свой язык")}}

{{escape ("Lütfen dilinizi seçiniz")}}
{{end}}

{{define "change_language_confirm"}}
{{escape ("Язык интерфейса был изменен на русский")}}
{{end}}

{{define "error"}}
{{escape ("ОЙ! Похоже, что что-то пошло не так. Попробуйте повторить свой запрос попозже")}}
{{end}}