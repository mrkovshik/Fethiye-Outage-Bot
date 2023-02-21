package crawling

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/pkg/errors"
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

func (oa OutageAydem) ConvertToOutage(ad []AydemData) ([]outage.Outage, error) {
	var err error
	res := make([]outage.Outage, 0)
	for _, i := range ad {
		if i.Area == "MUĞLA" {
			parsedEndDate, err := time.Parse(aydemTimeFormat, i.OutageEndDate[:16])
			if err != nil {
				err=errors.Wrap(err, "OutageEndDate parsing error")
				return []outage.Outage{}, err
			}
			parsedEndDate = parsedEndDate.Add(-3 * time.Hour)
			if parsedEndDate.After(time.Now().UTC()) {
				parsedStartDate, err := time.Parse(aydemTimeFormat, i.OutageStartDate[:16])
				if err != nil {
					err=errors.Wrap(err, "OutageStartDate parsing error")
					return []outage.Outage{}, err
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
	return res, err
}

func (oa OutageAydem) Crawl() ([]outage.Outage, error) {
	response, err := http.Get(oa.Url)
	if err != nil {
		err=errors.Wrap(err, "Error querying Aydem URL")
		return []outage.Outage{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		err=errors.Wrap(err, "Error reading response from Aydem")
		return []outage.Outage{}, err
	}
	var outages []AydemData
	err = json.Unmarshal(body, &outages)
	if err != nil {
		err=errors.Wrap(err, "Error Unmarshalling json from Aydem")
		return []outage.Outage{}, err
	}
	res, err := oa.ConvertToOutage(outages)
	if err != nil {
		return []outage.Outage{}, err
	}
	return oa.MergeStreets(res), err
}
