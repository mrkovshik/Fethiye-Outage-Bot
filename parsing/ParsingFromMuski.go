package parsing

import (
	"fmt"
	"log"
	"strconv"

	"strings"

	"github.com/PuerkitoBio/goquery"
	// "github.com/geziyor/geziyor"
	// "github.com/geziyor/geziyor/client"
	// "github.com/geziyor/geziyor/export"
	// "strconv"
	// "strings"
	"time"
)

type Outage struct {
	ID int
	City string
	District string
	StartDate time.Time
	Duration time.Duration
}


const oldTimeFormat = "02.01.2006 15:04"
func ParceFromMuski() {
rowSlice:=make([]Outage,0)
	doc, err := goquery.NewDocument("https://www.muski.gov.tr")
    if err != nil {
        log.Fatal(err)
    }
    table := doc.Find("table#plansiz")
    table.Find("tr").Each(func(i int, row *goquery.Selection) {
		fmt.Println("row ",  i)
		if i>2{
rowSlice=append(rowSlice, Outage{})
k:=i-3
	 row.Find("td").Each(func(j int, cell *goquery.Selection) {

		fmt.Println("cell ",  j, cell.Text())
		switch {
		case j==2:
			rowSlice[k].City=cell.Text()
		case j==3:
			rowSlice[k].District=cell.Text()
		case j==4:
			convertedDur, err := strconv.ParseInt(strings.Trim(cell.Text(), " Saat"), 0, 64)
			if err != nil {
				log.Fatal(err)
			}
			rowSlice[k].Duration=time.Duration(convertedDur)
		case j==5:
			parsedTime:=strings.Trim(cell.Text(), " ")
			rowSlice[k].StartDate,err=time.Parse(oldTimeFormat, parsedTime)
			if err != nil {
				log.Fatal(err)
			}

		}
		

    })
}
        fmt.Println(rowSlice)
    })

}



