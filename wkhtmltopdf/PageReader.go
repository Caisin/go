package wkhtmltopdf

import "io"

// PageReader is one input page (a HTML document) that is read from an io.Reader
// You can add only one Page from a reader
type PageReader struct {
	Input io.Reader
	PageOptions
}

// InputFile returns the input string and is part of the page interface
func (pr *PageReader) InputFile() string {
	return "-"
}

// Args returns the argument slice and is part of the page interface
func (pr *PageReader) Args() []string {
	return pr.PageOptions.Args()
}

//Reader returns the io.Reader and is part of the page interface
func (pr *PageReader) Reader() io.Reader {
	return pr.Input
}

// NewPageReader creates a new PageReader from an io.Reader
func NewPageReader(input io.Reader) *PageReader {
	return &PageReader{
		Input:       input,
		PageOptions: NewPageOptions(),
	}
}
