package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"testing"
)

func TestColly(t *testing.T) {
	c := colly.NewCollector()
	//debug
	//c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.OnHTML("html ul a[href]", func(e *colly.HTMLElement) {
		fmt.Println(e.Attr("href"))
	})
	_ = c.Visit("http://localhost/gitbook")

}

func TestIota(t *testing.T) {
	type Flags uint
	const (
		FlagUp           Flags = 1 << iota // is up
		FlagBroadcast                      //supports broadcast access capability
		FlagLoopback                       // is a loopback interface
		FlagPointToPoint                   // belongs to a point-to-point link
		FlagMulticast                      // supports multicast access capability
	)
	log.Print(FlagUp, FlagBroadcast, FlagLoopback, FlagPointToPoint, FlagMulticast)
	r := [...]int{99: -1}
	log.Print(r)
}
