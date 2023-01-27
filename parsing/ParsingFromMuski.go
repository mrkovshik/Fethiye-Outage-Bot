package parsing

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strconv"
	"strings"
	"time"
)

func ParceFromMuski() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://www.muski.gov.tr"},
		ParseFunc: quotesParse,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}
type Outage struct {
	ID int
	City string
	District string
	StartDate time.Time
	Duration time.Duration
}

const oldTimeFormat = "02.01.2006 15:04"

func quotesParse(g *geziyor.Geziyor, r *client.Response) {

	r.HTMLDoc.Find("table#plansiz tbody tr").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td").Map(func(_ int, s *goquery.Selection) string { return strings.TrimSpace(s.Text()) })
		parsedTime := cols[5]
		newtime, _ := time.Parse(oldTimeFormat, parsedTime)
		convertedDur, _ := strconv.ParseInt(strings.Trim(cols[4], " Saat"), 0, 64)
		newDuration := time.Duration(convertedDur)
		g.Exports <- map[string]interface{}{
			"City":              cols[2],
			"District":          cols[3],
			"Outage duration":   newDuration,
			"Outage start date": newtime,
		}

	})
	if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
		g.Get(r.JoinURL(href), quotesParse)
	}
}
