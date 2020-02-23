package html2pdf

import "github.com/Caisin/golearn/crawler"

type GitBookCrawler struct {
	crawler.Crawler
}

func (c *GitBookCrawler) ParseMenu() ([]string, error) {
	return nil, nil
}
func (c *GitBookCrawler) ParseBody(url string) (string, error) {
	return "", nil
}
