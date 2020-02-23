package html2pdf

import (
	wk "github.com/Caisin/golearn/wkhtmltopdf"
	"io/ioutil"
	"log"
	"strings"
)

func Html2pdf() {
	dir := "E:/code/Python/html2pdf/out/Go语言圣经（中文版）/html/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	pdfg, err := wk.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	pdfg.MarginRight.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.PageSize.Set(wk.PageSizeA4)
	str := ""
	for i := range files {
		path := dir + files[i].Name()
		file, err := ioutil.ReadFile(path)
		if err == nil {
			str += string(file)
		}
	}
	reader := strings.NewReader(str)
	page := wk.NewPageReader(reader)
	pdfg.AddPage(page)
	_ = pdfg.Create()
	_ = pdfg.WriteFile("./Go语言圣经.pdf")

}
