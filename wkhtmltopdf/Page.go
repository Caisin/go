package wkhtmltopdf

import "io"

// Page is the input struct for each page
type Page struct {
	Input string
	PageOptions
}

// InputFile returns the input string and is part of the page interface
func (p *Page) InputFile() string {
	return p.Input
}

// Args returns the argument slice and is part of the page interface
func (p *Page) Args() []string {
	return p.PageOptions.Args()
}

// Reader returns the io.Reader and is part of the page interface
func (p *Page) Reader() io.Reader {
	return nil
}

// NewPage creates a new input page from a local or web resource (filepath or URL)
func NewPage(input string) *Page {
	return &Page{
		Input:       input,
		PageOptions: NewPageOptions(),
	}
}
