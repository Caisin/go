package html2pdf

import (
	"github.com/Caisin/golearn/crawler"
)

type GitBookCrawler struct {
	*crawler.Crawler
}

func NewGitBookCrawler(startUrl, name, outPath, htmlTemplate, bodySelector, menuSelector string) (*GitBookCrawler, error) {
	newCrawler, err := crawler.NewCrawler(startUrl, name, outPath, htmlTemplate, bodySelector, menuSelector)
	if err != nil {
		return nil, err
	}
	gitBookCrawler := GitBookCrawler{Crawler: newCrawler}
	return &gitBookCrawler, nil
}
