package main

import (
	"golearn/html2pdf"
	"os"
	"testing"
)

func TestHtml2pd(t *testing.T) {
	os.Setenv("WKHTMLTOPDF_PATH", "D:/work/software/wkhtmltopdf/bin")
	html2pdf.Html2pdf()
}

func TestHtml2pd2(t *testing.T) {
	html2pdf.GitBookHtml2pdf()
}
