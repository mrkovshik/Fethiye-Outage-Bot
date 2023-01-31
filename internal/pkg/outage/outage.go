package outage

import "time"

type Outage struct {
	Resource  string
	City      string
	District  string
	StartDate time.Time
	Duration  time.Duration
	EndDate   time.Time
	SourceURL string
}