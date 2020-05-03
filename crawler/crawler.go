package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
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
	GenPdf       bool
	htmlPath     string
	htmTemplate  string
	StartUrl     string
	Name         string
	OutPath      string
	Path         string
	Domain       string
	Host         string
	bodySelector string
	menuSelector string
	menuMapping  map[string]string
	Cookies      map[string]string
	uri          *url.URL
	PoolSize     int
	colly        *colly.Collector
}

//添加内容
func (c *Crawler) AddBodyListener() {
	c.colly.OnHTML(c.bodySelector, func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		imgs := e.ChildAttrs("img", "src")
		//html, err := e.DOM.Html()
		isErr(err)
		tempMap := map[string]byte{}
		//go没有set
		for _, img := range imgs {
			tempMap[img] = 1
		}
		//处理图片url
		for s := range tempMap {
			absoluteURL := e.Request.AbsoluteURL(s)
			html = strings.Replace(html, s, absoluteURL, -1)
		}

		//处理a标签url
		hrefs := e.ChildAttrs("a", "href")
		for _, href := range hrefs {
			absoluteURL := e.Request.AbsoluteURL(href)
			localUrl := c.menuMapping[absoluteURL]
			if len(localUrl) > 0 {
				html = strings.Replace(html, href, localUrl, -1)
			}
		}

		//处理a标签url
		scripts := e.ChildAttrs("script", "src")
		for _, href := range scripts {
			absoluteURL := e.Request.AbsoluteURL(href)
			html = strings.Replace(html, href, absoluteURL, -1)
		}

		//处理a标签url
		links := e.ChildAttrs("link", "href")
		for _, href := range links {
			absoluteURL := e.Request.AbsoluteURL(href)
			html = strings.Replace(html, href, absoluteURL, -1)
		}

		replace := strings.Replace(c.htmTemplate, "{content}", html, -1)
		replace = strings.Replace(replace, "{start_url}", c.StartUrl, -1)
		replace = strings.Replace(replace, "{Domain}", c.Domain, -1)
		absoluteURL := e.Request.AbsoluteURL(e.Request.URL.RequestURI())
		i := c.menuMapping[absoluteURL]
		log.Printf("url is %s,index is %s", absoluteURL, i)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", c.htmlPath, i), []byte(replace), os.ModePerm)
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
	htmlPath := outPath + "/" + name + "/html/"
	crawler := &Crawler{
		Pdfg:         pdfg,
		DelHtml:      false,
		StartUrl:     startUrl,
		Name:         name,
		OutPath:      outPath,
		htmlPath:     htmlPath,
		Path:         parse.Path,
		Domain:       parse.Scheme + "://" + parse.Hostname(),
		Host:         parse.Hostname(),
		uri:          parse,
		PoolSize:     10,
		colly:        collector,
		htmTemplate:  htmTemplate,
		bodySelector: bodySelector,
		menuSelector: menuSelector,
		GenPdf:       true,
	}
	return crawler, nil
}

//解析菜单
func (c *Crawler) ParseMenu() ([]string, error) {
	var menus []string
	c.colly.OnHTML("html", func(e *colly.HTMLElement) {
		attrs := e.ChildAttrs(c.menuSelector, "href")
		m := map[string]string{}
		for i := range attrs {
			absoluteURL := e.Request.AbsoluteURL(attrs[i])
			menus = append(menus, absoluteURL)
			m[absoluteURL] = fmt.Sprintf("%.4d.html", i)
		}
		c.menuMapping = m
	})
	_ = c.colly.Visit(c.StartUrl)
	c.colly.Wait()
	return menus, nil
}

func (c *Crawler) InitParam() {
	if c.Cookies != nil {
		cookies := make([]*http.Cookie, len(c.Cookies))
		i := 0
		for k, v := range c.Cookies {
			cookie := new(http.Cookie)
			cookie.Name = k
			cookie.Value = v
			cookie.Path = "/"
			cookie.Domain = c.Host
			cookies[i] = cookie
			i++
		}
		_ = c.colly.SetCookies("http://grace.hcoder.net", cookies)
	}
}
func (c *Crawler) Run() {
	c.InitParam()
	menus, err := c.ParseMenu()
	if isErr(err) {
		return
	}
	pool, err := ants.NewPool(c.PoolSize)
	if isErr(err) {
		return
	}
	_ = os.MkdirAll(c.htmlPath, os.ModePerm)
	defer pool.Release()
	var wg sync.WaitGroup
	log.Println(len(c.menuMapping))
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
	if c.GenPdf {
		err = c.Htm2Pdf()
		isErr(err)
	}
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
