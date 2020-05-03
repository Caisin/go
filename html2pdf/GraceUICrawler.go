package html2pdf

import (
	"github.com/Caisin/golearn/crawler"
)

type GraceUICrawler struct {
	*crawler.Crawler
}

func NewGraceUICrawler(startUrl, name, outPath, htmlTemplate, bodySelector, menuSelector string) (*GraceUICrawler, error) {
	newCrawler, err := crawler.NewCrawler(startUrl, name, outPath, htmlTemplate, bodySelector, menuSelector)
	if err != nil {
		return nil, err
	}
	graceUICrawler := GraceUICrawler{Crawler: newCrawler}
	return &graceUICrawler, nil
}
