package crawler

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

// Crawler is for crawler message news
type Crawler struct {
	l *zap.SugaredLogger
}

// NewCrawler return new crawler instance
func NewCrawler(l *zap.SugaredLogger) *Crawler {
	return &Crawler{
		l: l,
	}
}

func (cr *Crawler) crawlerNews() {
	fmt.Println("start crawling")
	// Instantiate default collector
	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML(".timeline-sec > ul:first-child > li", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText(".timeline-head"))
		fmt.Println(e.ChildText("p"))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://ncov.moh.gov.vn/")
}

// Start the crawl
func (cr *Crawler) Start() {
	cr.crawlerNews()
}