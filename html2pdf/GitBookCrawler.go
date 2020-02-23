package html2pdf

import (
	"errors"
	"github.com/Caisin/golearn/crawler"
	"github.com/gocolly/colly"
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
func (c *GitBookCrawler) ParseMenu(e *colly.HTMLElement) ([]string, error) {
	var menus []string
	attrs := e.ChildAttrs("ul.summary a", "href")
	if attrs == nil {
		return nil, errors.New("menus is empty")
	}
	for i := range attrs {
		url := e.Request.AbsoluteURL(attrs[i])
		menus = append(menus, url)
	}
	return menus, nil
}
func (c *GitBookCrawler) ParseBody(e *colly.HTMLElement) (string, error) {
	text := e.ChildText("section.markdown-section")
	var err error
	if text == "" {
		err = errors.New("body is empty")
	}
	return text, err
}
