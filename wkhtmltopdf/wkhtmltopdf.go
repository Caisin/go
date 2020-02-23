// Package wkhtmltopdf contains wrappers around the wkhtmltopdf commandline tool
package wkhtmltopdf

import (
	"io"
)

type page interface {
	Args() []string
	InputFile() string
	Reader() io.Reader
}

// cover page
type cover struct {
	Input string
	pageOptions
}

// table of contents
type toc struct {
	Include bool
	allTocOptions
}

type allTocOptions struct {
	pageOptions
	tocOptions
}
