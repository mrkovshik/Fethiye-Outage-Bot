{{define "badQuery"}}
I am sorry, but I can't find anythithg like _*'{{.Text}}'*_

Maybe we should try again?
Please print your city and neigbourhood divided by space, for example _*'Fethie Taşyaka'*_"
{{end}}

{{define "listOutages"}}
{{if eq (len .) 0 }}
*There is no outages planned in your neigborhood in the closest time*
{{else}}
*Here are the closest outages found for your neigborhood:*
{{range .}}
*{{.Resource}} outage* from {{escape (format (.StartDate))}} to {{escape (format (.EndDate))}}{{if gt (len .Notes ) 3 }}
*In the next areas and streets:*
{{escape (.Notes)}}{{end}}
{{end}}
{{end}}
{{end}}


{{define "confirmDistr"}}
Did you mean _*{{.City}} {{.Name}}*_?{{end}}

{{define "startMsg"}}
Please print your city and neigbourhood divided by space, for example _*'Fethiye Taşyaka'*_
{{end}}