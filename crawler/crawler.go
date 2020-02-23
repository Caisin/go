package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	wk "github.com/Caisin/golearn/wkhtmltopdf"
	"github.com/panjf2000/ants/v2"
)

type Crawler struct {
	Pdfg         *wk.PDFGenerator
	DelHtml      bool
	htmlPath     string
	htmTemplate  string
	StartUrl     string
	Name         string
	OutPath      string
	Path         string
	Domain       string
	bodySelector string
	menuSelector string
	menuIndex    map[string]int
	uri          *url.URL
	PoolSize     int
	colly        *colly.Collector
}

//添加内容
func (c *Crawler) AddBodyListener() {
	c.colly.OnHTML(c.bodySelector, func(e *colly.HTMLElement) {
		html, err := e.DOM.Parent().Html()
		imgs := e.ChildAttrs("img", "src")
		//html, err := e.DOM.Html()
		isErr(err)
		tempMap := map[string]byte{}
		for _, img := range imgs {
			tempMap[img] = 1
		}
		for s := range tempMap {
			absoluteURL := e.Request.AbsoluteURL(s)
			html = strings.Replace(html, s, absoluteURL, -1)
		}
		replace := strings.Replace(c.htmTemplate, "{content}", html, -1)
		replace = strings.Replace(replace, "{start_url}", c.StartUrl, -1)
		absoluteURL := e.Request.AbsoluteURL(e.Request.URL.RequestURI())
		i := c.menuIndex[absoluteURL]
		log.Printf("url is %s,index is %d", absoluteURL, i)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%.4d.html", c.htmlPath, i), []byte(replace), os.ModePerm)
	})
}
func NewCrawler(startUrl, name, outPath, htmTemplate, bodySelector, menuSelector string) (*Crawler, error) {
	parse, err := url.Parse(startUrl)
	if err != nil {
		return nil, err
	}
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	collector := colly.NewCollector()
	htmlPath := outPath + "/html/"
	crawler := &Crawler{
		Pdfg:         pdfg,
		DelHtml:      false,
		StartUrl:     startUrl,
		Name:         name,
		OutPath:      outPath,
		htmlPath:     htmlPath,
		Path:         parse.Path,
		Domain:       parse.Scheme + "://" + parse.Hostname(),
		uri:          parse,
		PoolSize:     10,
		colly:        collector,
		htmTemplate:  htmTemplate,
		bodySelector: bodySelector,
		menuSelector: menuSelector,
	}
	return crawler, nil
}

//解析菜单
func (c *Crawler) ParseMenu() ([]string, error) {
	var menus []string
	c.colly.OnHTML("html", func(e *colly.HTMLElement) {
		attrs := e.ChildAttrs(c.menuSelector, "href")
		m := map[string]int{}
		for i := range attrs {
			absoluteURL := e.Request.AbsoluteURL(attrs[i])
			menus = append(menus, absoluteURL)
			m[absoluteURL] = i
		}
		c.menuIndex = m
	})
	_ = c.colly.Visit(c.StartUrl)
	c.colly.Wait()
	return menus, nil
}

//解析菜单
func (c *Crawler) ParseBody() string {
	panic("not implement")
}

func (c *Crawler) Run() {
	menus, err := c.ParseMenu()
	if isErr(err) {
		return
	}
	pool, err := ants.NewPool(c.PoolSize)
	if isErr(err) {
		return
	}
	os.MkdirAll(c.htmlPath, os.ModePerm)
	defer pool.Release()
	var wg sync.WaitGroup
	log.Println(len(c.menuIndex))
	c.AddBodyListener()
	for i := range menus {
		wg.Add(1)
		link := menus[i]
		_ = pool.Submit(func() {
			_ = c.colly.Visit(link)
			wg.Done()
		})
	}
	wg.Wait()
	err = c.Htm2Pdf()
	isErr(err)
}

func isErr(err error) bool {
	if err != nil {
		log.Fatalf("%s\n%v", err.Error(), err)
		return true
	}
	return false
}
func (c *Crawler) PdfGSetting() {
	c.Pdfg.MarginRight.Set(0)
	c.Pdfg.MarginLeft.Set(0)
	c.Pdfg.MarginTop.Set(4)
	c.Pdfg.MarginBottom.Set(4)
	c.Pdfg.PageSize.Set(wk.PageSizeA4)
}

func (c *Crawler) Htm2Pdf() error {

	files, err := ioutil.ReadDir(c.htmlPath)
	if err != nil {
		return err
	}
	c.PdfGSetting()
	str := ""
	for i := range files {
		path := c.htmlPath + files[i].Name()
		file, err := ioutil.ReadFile(path)
		if err == nil {
			str += string(file)
		}
	}
	reader := strings.NewReader(str)
	page := wk.NewPageReader(reader)
	c.Pdfg.AddPage(page)
	err = c.Pdfg.Create()
	if err != nil {
		return err
	}
	err = c.Pdfg.WriteFile(c.OutPath + "/" + c.Name + ".pdf")
	if err != nil {
		return err
	}
	log.Println("生成Pdf成功...")
	log.Printf("输出目录:%s", c.OutPath)
	if c.DelHtml {
		err = os.RemoveAll(c.htmlPath)
		if err != nil {
			log.Fatalf("删除html文件夹[%s]失败\n%v", c.htmlPath, err)
		}
	}
	return nil
}
