package wkhtmltopdf

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func ExampleNewPDFGenerator() {

	// Create new PDF generator
	pdfg, err := NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(OrientationLandscape)
	pdfg.Grayscale.Set(true)

	// Create a new input page from an URL
	page := NewPage("https://godoc.org/github.com/SebastiaanKlippert/go-wkhtmltopdf")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done

}

func ExampleNewPDFGeneratorFromJSON() {

	const html = `<!doctype html><html><head><title>WKHTMLTOPDF TEST</title></head><body>HELLO PDF</body></html>`

	// Client code
	pdfg := NewPDFPreparer()
	pdfg.AddPage(NewPageReader(strings.NewReader(html)))
	pdfg.Dpi.Set(600)

	// The html string is also saved as base64 string in the JSON file
	jsonBytes, err := pdfg.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	// The JSON can be saved, uploaded, etc.

	// Server code, create a new PDF generator from JSON, also looks for the wkhtmltopdf executable
	pdfgFromJSON, err := NewPDFGeneratorFromJSON(bytes.NewReader(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}

	// Create the PDF
	err = pdfgFromJSON.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Use the PDF
	fmt.Printf("PDF size %d bytes", pdfgFromJSON.Buffer().Len())
}

func TestPageReader_InputFile(t *testing.T) {
	dir := "E:/code/Python/html2pdf/out/Go语言圣经（中文版）/html/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	pdfg, err := NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.globalOptions.MarginRight.Set(0)
	pdfg.globalOptions.MarginLeft.Set(0)
	pdfg.globalOptions.MarginTop.Set(0)
	pdfg.globalOptions.MarginBottom.Set(0)
	pdfg.globalOptions.PageSize.Set(PageSizeLetter)
	str := ""
	for i := range files {
		path := dir + files[i].Name()
		file, err := ioutil.ReadFile(path)
		if err == nil {
			str += string(file)
		}
	}
	reader := strings.NewReader(str)
	page := NewPageReader(reader)
	pdfg.AddPage(page)
	_ = pdfg.Create()
	_ = pdfg.WriteFile("./Go语言圣经.pdf")
}

func TestBoolOption_Parse(t *testing.T) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	_ = pdf.OutputFileAndClose("hello.pdf")
}
