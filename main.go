package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func crawler(wg *sync.WaitGroup, urls chan string, doneUrls chan string) {
	defer wg.Done()

	log.Println("Worker started")
	c := colly.NewCollector()

	c.OnHTML("article", func(e *colly.HTMLElement) {
		metaTags := e.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, s *goquery.Selection) {

		})
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		err := c.Visit(e.Request.AbsoluteURL(link))
		if err != nil {
			fmt.Println("skipped.")
			count++
		}
	})

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		log.Println(err)
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawled:", r.URL.String())
		count++
		fmt.Println("counter:", count)
	})

	err = c.Visit(<-urls)
	if err != nil {
		log.Fatal(err)
	}
}

var count int64

func main() {
	initial := []string{"https://golang.org", "http://cnn.com", "https://yahoo.com"}
	init := "https://stackoverflow.com"
	urls := make(chan string, 100)
	doneUrls := make(chan string)
	wg := &sync.WaitGroup{}

	urls <- init

	for i, item := range initial {
		_ = i
		urls <- item
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		log.Println("Starting worker")
		go crawler(wg, urls, doneUrls)
	}
	wg.Wait()
}
