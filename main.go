package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strconv"
)

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://www.muski.gov.tr"},
		ParseFunc: quotesParse,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("body").Each(func(i int, s *goquery.Selection) {

		sel := (r.HTMLDoc.Find("#ContentPlaceHolder1_Repeater3_tr1_" + strconv.Itoa(i) + " > td:nth-child(3)")).Text()
		fmt.Println(sel + "gh")
		for i := 0; i < 10; i++ {

			g.Exports <- map[string]interface{}{
				"Округ":        s.Find("#ContentPlaceHolder1_Repeater3_tr1_" + strconv.Itoa(i) + " > td:nth-child(3)").Text(),
				"Район":        s.Find("#ContentPlaceHolder1_Repeater3_tr1_" + strconv.Itoa(i) + " > td:nth-child(4)").Text(),
				"Длительность": s.Find("#ContentPlaceHolder1_Repeater3_tr1_" + strconv.Itoa(i) + " > td:nth-child(5)").Text(),
				"Дата и время отключения": s.Find("#ContentPlaceHolder1_Repeater3_tr1_" + strconv.Itoa(i) + " > td:nth-child(6)").Text(),
			}
		}
	})
	if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
		g.Get(r.JoinURL(href), quotesParse)
	}
}
