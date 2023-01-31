package outage

import "time"

type Outage struct {
	Resource  string
	City      string
	District  string
	StartDate time.Time
	Duration  int
	EndDate   time.Time
	SourceURL string
}