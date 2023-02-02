package crawling

import (
	// "database/sql"
	// "fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
)



type OutageMuski struct {
	Url      string
	Resource string
}



const oldTimeFormat = "02.01.2006 15:04"

func (om OutageMuski) expandDistr(s []outage.Outage) []outage.Outage {
	addition := make([]outage.Outage, 0)
	for i := range s {
		distList := strings.Split(s[i].District, ", ")
		if len(distList) > 1 {
			s[i].District = distList[0]
			for j := 1; j < len(distList); j++ {
				addition = append(addition, s[i])
				addition[j-1].District = distList[j]
			}
			s = append(s, addition...)
		}
	}
	return s
}

func (om OutageMuski) parseTable(table *goquery.Selection) []outage.Outage {
		rowSlice := make([]outage.Outage, 0)
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		if i > 2 {
			parsedRow:= outage.Outage{}
				row.Find("td").Each(func(j int, cell *goquery.Selection) {
				parsedRow.Resource = om.Resource
				parsedRow.SourceURL = om.Url
				switch {
				case j == 2:
					parsedRow.City = cell.Text()
				case j == 3:
					parsedRow.District = cell.Text()

				case j == 4:
					parsedDur, err := strconv.ParseInt(strings.Trim(cell.Text(), " Saat"), 0, 64)
					if err != nil {
						log.Fatal(err)
					}
					parsedRow.Duration = time.Duration(parsedDur) * time.Hour
				case j == 5:
					parsedTime,err := time.Parse(oldTimeFormat,strings.Trim(cell.Text(), " ") )
					parsedRow.StartDate= parsedTime.Add(-3*time.Hour)

					if err != nil {
						log.Fatal(err)
					}
					parsedRow.EndDate = parsedRow.StartDate.Add(parsedRow.Duration)
				}
				})
				if parsedRow.EndDate.UTC().After(time.Now().UTC()){
				rowSlice = append(rowSlice, parsedRow)
				}
		}
	})
	return rowSlice
}

func (osrc OutageMuski) Crawl() []outage.Outage {
	outages := make([]outage.Outage, 0)
	doc, err := goquery.NewDocument(osrc.Url)
	if err != nil {
		log.Fatal(err)
	}
	table := doc.Find("table#plansiz")
	rowSlice := osrc.parseTable(table)
	rowSlice = osrc.expandDistr(rowSlice)
	outages = append(outages, rowSlice...)
	return outages
}
