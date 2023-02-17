package outage

import "time"

type Outage struct {
	Resource  string
	City      string
	District  string
	StartDate time.Time
	Duration  time.Duration
	EndDate   time.Time
	Notes     string
	SourceURL string
}

func (one Outage) Equal(another Outage) bool {
	return one.City == another.City && one.District == another.District && one.StartDate.Equal(another.StartDate) && one.EndDate.Equal(another.EndDate) && one.Resource == another.Resource
}
