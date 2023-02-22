package crawling

import (
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
)

type Crawler interface {
	Crawl() ([]outage.Outage, error)
}

func CrawlOutages(c Crawler) ([]outage.Outage, error) {
	return c.Crawl()
}
