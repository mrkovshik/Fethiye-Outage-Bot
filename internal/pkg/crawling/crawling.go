package crawling

import "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"

type Crawler interface {
	Crawl() []outage.Outage
}

func CrawlOutages(c Crawler) []outage.Outage {
	return c.Crawl()
}
