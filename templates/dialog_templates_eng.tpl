{{define "listOutages"}}
{{if eq (len .) 0 }}
    *There are no outages planned in your neighborhood in the nearest future*
{{else}}
*Here are the planned outages found for your neighborhood:*
{{range .}}
*{{.Resource}} outage* from {{escape (format (.StartDate))}} to {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}
*In the following areas and streets:*
{{escape (.Notes)}}{{end}}
{{end}}
{{end}}
{{end}}

{{define "alert"}}
Please be aware that there is a *{{.Resource}} outage* planned in *{{escape (.District)}} {{.City}}* from {{escape (format (.StartDate))}} to {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}
*In the following areas and streets:*
{{escape (.Notes)}}{{end}}
{{end}}


{{define "name_greet"}}
Hi {{escape (.UserName)}}{{escape ("!")}}
{{end}}

{{define "mainMenu_greet"}}
What would you like to do?
{{end}}

{{define "pickCity_greet"}}
Please, pick your city from the list below:
{{end}}

{{define "settings_greet"}}
You can modify your subscriptions here{{escape (".")}} What would you like to change?
{{end}}

{{define "claim_buttons"}}
Please pick a button from the menu below, *do not type your answer*:
{{end}}

{{define "pickCity_confirm"}}
You have chosen *{{.PickedCity}}* city{{escape ("!")}} Now pick your neighborhood from the list below:
{{end}}

{{define "pickPeriod_greet"}}
Now pick the alert period:
{{end}}

{{define "change_period_greet"}}
 All right, pick a new alert period:
{{end}}

{{define "pickDistr_confirm"}}
You have chosen *{{escape (.PickedDistrict)}}* neighborhood in *{{.PickedCity}}* city{{escape ("!")}}
{{end}}

{{define "change_location_confirm"}}
Your subscription has been sucsessfully updated{{escape ("!")}}
{{end}}

{{define "change_period_confirm"}}
Your subscription has been sucsessfully updated{{escape ("!")}} From now on you will get notifications about outages in your neighborhood in *{{.PickedPeriod}}* hours before it starts{{escape (".")}}
{{end}}

{{define "set_period_confirm"}}
Your subscription has been sucsessfully set{{escape ("!")}} From now on you will get notifications about outages in *{{escape (.PickedDistrict)}}* *{{.PickedCity}}* in *{{.PickedPeriod}}* hours before it starts{{escape (".")}}
{{end}}

{{define "show_sub"}}
You are subscribed to get notifications about outages in *{{escape (.District)}}* neighborhood in *{{.City}}* city in *{{.Period}}* hours before it starts{{escape (".")}}
{{end}}

{{define "no_subs"}}
It seems you do not have a subscription yet{{escape (".")}} Pick *Subscribe for alerts* button to get one{{escape (".")}}
{{end}}

{{define "have_sub"}}
It seems you already have a subscription{{escape (".")}} Pick *Subscription settings* button to modify or cancel your subscription{{escape (".")}}
{{end}}

{{define "cancel_confirm"}}
Your subscription has been sucsessfully cancelled{{escape ("!")}}
{{end}}

{{define "cancel_you_sure"}}
Are you sure you want to cancel your subscription{{escape ("?")}}
{{end}}

{{define "go_back"}}
All right, let{{escape ("'")}}s go back to the previous step
{{end}}

{{define "press_start"}}
I am sorry, I am afraid I forgot what we were talking about{{escape ("(")}}
Let{{escape ("'")}}s start from the very beginning{{escape ("!")}}
{{end}}

{{define "error"}}
OOOPS{{escape ("!")}} We are very sorry, but it looks like something went wrong{{escape (".")}} Please try again later{{escape (".")}}
{{end}}