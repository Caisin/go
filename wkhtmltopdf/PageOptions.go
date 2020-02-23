package wkhtmltopdf

// PageOptions are options for each input page
type PageOptions struct {
	pageOptions
	headerAndFooterOptions
}

// Args returns the argument slice
func (po *PageOptions) Args() []string {
	return append(append([]string{}, po.pageOptions.Args()...), po.headerAndFooterOptions.Args()...)
}

// NewPageOptions returns a new PageOptions struct with all options
func NewPageOptions() PageOptions {
	return PageOptions{
		pageOptions:            newPageOptions(),
		headerAndFooterOptions: newHeaderAndFooterOptions(),
	}
}
