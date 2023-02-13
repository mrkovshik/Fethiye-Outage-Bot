package crawling

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
)

type OutageAydem struct {
	Url      string
	Resource string
}

type AydemData struct {
	Area            string `json:"Sehir"`
	City            string `json:"Ilce"`
	District        string `json:"Mahalle"`
	Street          string `json:"Sokak"`
	OutageStartDate string `json:"Planlanan_Baslangic_Zamani"`
	OutageEndDate   string `json:"Planlanan_Sona_Erme_Zamani"`
}

const aydemTimeFormat = "2006-01-02 15:04"

func (oa OutageAydem) MergeStreets(o []outage.Outage) []outage.Outage {
	dict := make(map[string]string)
	res := make([]outage.Outage, 0)
	for _, i := range o {
		key := i.City + i.District + i.StartDate.String()
		if _, ok := dict[key]; !ok {
			res = append(res, i)
		}
		if i.Notes != "" {
			dict[key] += i.Notes + "/ "
		} else {
			dict[key] += ""
		}
	}
	for i, j := range res {
		key := j.City + j.District + j.StartDate.String()
		res[i].Notes = dict[key]
	}
	return res
}

func (oa OutageAydem) ConvertToOutage(ad []AydemData) []outage.Outage {
	res := make([]outage.Outage, 0)
	for _, i := range ad {
		if i.Area == "MUĞLA" {
			parsedEndDate, err := time.Parse(aydemTimeFormat, i.OutageEndDate[:16])
			if err != nil {
				log.Fatal(err)
			}
			parsedEndDate = parsedEndDate.Add(-3 * time.Hour)
			if parsedEndDate.After(time.Now().UTC()) {
				parsedStartDate, err := time.Parse(aydemTimeFormat, i.OutageStartDate[:16])
				if err != nil {
					log.Fatal(err)
				}
				parsedStartDate = parsedStartDate.Add(-3 * time.Hour)
				o := outage.Outage{}
				o.Resource = oa.Resource
				o.City = i.City
				o.District = strings.Trim((strings.Trim(i.District, " Mh.")), " MH.")
				o.EndDate = parsedEndDate
				o.StartDate = parsedStartDate
				o.Notes = i.Street
				o.SourceURL = oa.Url
				res = append(res, o)
			}
		}
	}
	return res
}

func (oa OutageAydem) Crawl() []outage.Outage {
	response, err := http.Get(oa.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var outages []AydemData
	err = json.Unmarshal(body, &outages)
	if err != nil {
		log.Fatal(err)
	}
	res := oa.ConvertToOutage(outages)
	return oa.MergeStreets(res)
}
