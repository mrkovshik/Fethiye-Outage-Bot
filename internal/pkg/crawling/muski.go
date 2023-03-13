package crawling

import (
	// "database/sql"
	// "fmt"

	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/util"
	"github.com/pkg/errors"
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

func (om OutageMuski) parseTable(table *goquery.Selection) ([]outage.Outage, error) {
	var rowSlice []outage.Outage
	rows := table.Find("tr")
	for i := range rows.Nodes {
		if i <= 2 {
			continue
		}
		parsedRow := outage.Outage{}
		cells := rows.Eq(i).Find("td")
		for j := range cells.Nodes {
			switch j {
			case 2:
				parsedRow.City = cells.Eq(j).Text()
			case 3:
				parsedRow.District = cells.Eq(j).Text()
			case 4:
				parsedDur, err := strconv.ParseInt(strings.TrimSuffix(cells.Eq(j).Text(), " Saat"), 0, 64)
				if err != nil {
					err = errors.Wrap(err, "Error parsing time from muski Table")
					return []outage.Outage{}, err
				}
				parsedRow.Duration = time.Duration(parsedDur) * time.Hour
			case 5:
				parsedTime, err := time.Parse(oldTimeFormat, strings.Trim(cells.Eq(j).Text(), " "))
				if err != nil {
					err = errors.Wrap(err, "Error parsing time from muski Table")
					return []outage.Outage{}, err
				}
				parsedRow.StartDate = parsedTime.Add(-3 * time.Hour)
				parsedRow.EndDate = parsedRow.StartDate.Add(parsedRow.Duration)
				parsedRow.Alerted=false
			}
		}
		if parsedRow.EndDate.UTC().After(time.Now().UTC()) {
			parsedRow.Notes = ""
			parsedRow.Resource = om.Resource
			parsedRow.SourceURL = om.Url
			nc, err := util.Normalize(parsedRow.City)
			if err != nil {
				return []outage.Outage{}, err
			}
			parsedRow.CityNormalized = strings.Join(nc, " ")
			nd, err := util.Normalize(parsedRow.District)
			if err != nil {
				return []outage.Outage{}, err
			}
			parsedRow.DistrictNormalized = strings.Join(nd, " ")
			rowSlice = append(rowSlice, parsedRow)
		}
	}

	return rowSlice, nil
}

func (om OutageMuski) Crawl() ([]outage.Outage, error) {
	//TODO add validation here
	client := &http.Client{}
	req, err := http.NewRequest("GET", om.Url, nil)
	if err != nil {
		err = errors.Wrap(err, "Error wrapping request for Muski URL")
		return []outage.Outage{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "Error requesting Muski URL")
		return []outage.Outage{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(err, "Response from Muski is not OK")
		return []outage.Outage{}, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading Muski response")
		return []outage.Outage{}, err
	}
	table := doc.Find("table#plansiz")
	rowSlice, err := om.parseTable(table)
	if err != nil {
		return []outage.Outage{}, err
	}
	rowSlice = om.expandDistr(rowSlice)
	outages := []outage.Outage{}
	outages = append(outages, rowSlice...)
	return outages, err
}
