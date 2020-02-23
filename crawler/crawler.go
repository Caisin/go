package crawler

import (
	"errors"
	"fmt"
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
	Pdfg     *wk.PDFGenerator
	DelHtml  bool
	htmlPath string
	StartUrl string
	Name     string
	OutPath  string
	Path     string
	Domain   string
	uri      *url.URL
	PoolSize int
}

func NewCrawler(startUrl, name, outPath string) (*Crawler, error) {
	parse, err := url.Parse(startUrl)
	if err != nil {
		return nil, err
	}
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	return &Crawler{
		Pdfg:     pdfg,
		DelHtml:  false,
		StartUrl: startUrl,
		Name:     name,
		OutPath:  outPath,
		htmlPath: outPath + "/html/",
		Path:     parse.Path,
		Domain:   parse.Scheme + "://" + parse.Hostname(),
		uri:      parse,
		PoolSize: 20,
	}, nil
}

//解析菜单
func (c *Crawler) Request() ([]string, error) {
	return nil, errors.New("not implement")
}

//解析菜单
func (c *Crawler) ParseMenu() ([]string, error) {
	return nil, errors.New("not implement")
}

//解析菜单
func (c *Crawler) ParseBody(url string) (string, error) {
	return "", errors.New("not implement")
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
	defer pool.Release()
	var wg sync.WaitGroup
	for i := range menus {
		wg.Add(1)
		_ = pool.Submit(func() {
			c.body2File(&wg, i)
		})
	}
	wg.Wait()
	err = c.Htm2Pdf()
	isErr(err)
}

func (c *Crawler) body2File(wg *sync.WaitGroup, idx int) {
	defer wg.Done()
	body, err := c.ParseBody()
	isErr(err)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%5d.html", c.htmlPath, idx), []byte(body), 0777)
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
	c.Pdfg.MarginTop.Set(10)
	c.Pdfg.MarginBottom.Set(10)
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
