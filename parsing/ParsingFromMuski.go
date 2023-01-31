package parsing

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/mrkovshik/Fethiye-Outage-Bot/postgresdb"
	"github.com/PuerkitoBio/goquery"
)

type Outage struct {
	Resource  string      
	City      string      
	District  string      
	StartDate time.Time    
	Duration  int		   
	EndDate   time.Time   
	SourceURL	  string	
}

func (o Outage) unmarshal() postgresdb.outageRow {
	return postgresdb.outageRow{
		Resource : o.Resource,
		City   :  o.City,
		District : o.District,      
		StartDate: o.StartDate,  
		Duration: o.Duration,		   
		EndDate: o.EndDate,
		SourceURL: o.SourceURL,
	}
}
type WaterOutage struct {
   Outages	[] Outage
}
 
type PowerOutage struct {
	Outages	[] Outage
}

type Crawler interface {
	Crawl()
}

const oldTimeFormat = "02.01.2006 15:04"




func expandDistr (s [] Outage) [] Outage {
	addition:=make([]Outage,0)
	for i:= range s {
		distList:=strings.Split(s[i].District,", ")
if len(distList)>1{
	s[i].District=distList[0]
	for j:=1;j<len(distList);j++ {
		addition=append(addition, s[i])
		addition[j-1].District=distList[j]
	}
	s = append(s, addition...)
	}
	}
	return s
}

func parseTable (table *goquery.Selection) [] Outage{
	var err error
	rowSlice := make([]Outage, 0)
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		if i > 2 {
		rowSlice = append(rowSlice, Outage{})
		k := i - 3
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			rowSlice[k].Resource="water"
			rowSlice[k].SourceURL="https://www.muski.gov.tr"
			switch {
			case j == 2:
				rowSlice[k].City = cell.Text()
			case j == 3:
								rowSlice[k].District = cell.Text()
				
			case j == 4:
				parsedDur, err := strconv.ParseInt(strings.Trim(cell.Text(), " Saat"), 0, 64)
				if err != nil {
					log.Fatal(err)
				}
				rowSlice[k].Duration = int(parsedDur)
			case j == 5:
				parsedTime := strings.Trim(cell.Text(), " ")
				rowSlice[k].StartDate, err = time.Parse(oldTimeFormat, parsedTime)
				if err != nil {
					log.Fatal(err)
				}
				rowSlice[k].EndDate=rowSlice[k].StartDate.Add(time.Duration(rowSlice[k].Duration)*time.Hour)
			}

		})
	}
})
rowSlice=expandDistr(rowSlice)
return  rowSlice
}

func (wo WaterOutage) Crawl()  {
	
	doc, err := goquery.NewDocument("https://www.muski.gov.tr")
	if err != nil {
		log.Fatal(err)
	}
	table := doc.Find("table#plansiz")
	rowSlice:=parseTable(table)

	for _,i:=range rowSlice {
			fmt.Printf("%+v\n", i)
		}
		wo.Outages=append(wo.Outages, rowSlice...)
}
