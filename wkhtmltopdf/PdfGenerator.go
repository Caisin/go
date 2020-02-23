package wkhtmltopdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var binPath stringStore

// SetPath sets the path to wkhtmltopdf
func SetPath(path string) {
	binPath.Set(path)
}

// GetPath gets the path to wkhtmltopdf
func GetPath() string {
	return binPath.Get()
}

// PDFGenerator is the main wkhtmltopdf struct, always use NewPDFGenerator to obtain a new PDFGenerator struct
type PDFGenerator struct {
	globalOptions
	outlineOptions

	Cover      cover
	TOC        toc
	OutputFile string //filename to write to, default empty (writes to internal buffer)

	binPath   string
	outbuf    bytes.Buffer
	outWriter io.Writer
	stdErr    io.Writer
	pages     []page
}

//Args returns the commandline arguments as a string slice
func (pdfg *PDFGenerator) Args() []string {
	args := append([]string{}, pdfg.globalOptions.Args()...)
	args = append(args, pdfg.outlineOptions.Args()...)
	if pdfg.Cover.Input != "" {
		args = append(args, "cover")
		args = append(args, pdfg.Cover.Input)
		args = append(args, pdfg.Cover.pageOptions.Args()...)
	}
	if pdfg.TOC.Include {
		args = append(args, "toc")
		args = append(args, pdfg.TOC.pageOptions.Args()...)
		args = append(args, pdfg.TOC.tocOptions.Args()...)
	}
	for _, page := range pdfg.pages {
		args = append(args, "page")
		args = append(args, page.InputFile())
		args = append(args, page.Args()...)
	}
	if pdfg.OutputFile != "" {
		args = append(args, pdfg.OutputFile)
	} else {
		args = append(args, "-")
	}
	return args
}

// ArgString returns Args as a single string
func (pdfg *PDFGenerator) ArgString() string {
	return strings.Join(pdfg.Args(), " ")
}

// AddPage adds a new input page to the document.
// A page is an input HTML page, it can span multiple pages in the output document.
// It is a Page when read from file or URL or a PageReader when read from memory.
func (pdfg *PDFGenerator) AddPage(p page) {
	pdfg.pages = append(pdfg.pages, p)
}

// SetPages resets all pages
func (pdfg *PDFGenerator) SetPages(p []page) {
	pdfg.pages = p
}

// ResetPages drops all pages previously added by AddPage or SetPages.
// This allows reuse of current instance of PDFGenerator with all of it's configuration preserved.
func (pdfg *PDFGenerator) ResetPages() {
	pdfg.pages = []page{}
}

// Buffer returns the embedded output buffer used if OutputFile is empty
func (pdfg *PDFGenerator) Buffer() *bytes.Buffer {
	return &pdfg.outbuf
}

// Bytes returns the output byte slice from the output buffer used if OutputFile is empty
func (pdfg *PDFGenerator) Bytes() []byte {
	return pdfg.outbuf.Bytes()
}

// SetOutput sets the output to write the PDF to, when this method is called, the internal buffer will not be used,
// so the Bytes(), Buffer() and WriteFile() methods will not work.
func (pdfg *PDFGenerator) SetOutput(w io.Writer) {
	pdfg.outWriter = w
}

// SetStderr sets the output writer for Stderr when running the wkhtmltopdf command. You only need to call this when you
// want to print the output of wkhtmltopdf (like the progress messages in verbose mode). If not called, or if w is nil, the
// output of Stderr is kept in an internal buffer and returned as error message if there was an error when calling wkhtmltopdf.
func (pdfg *PDFGenerator) SetStderr(w io.Writer) {
	pdfg.stdErr = w
}

// WriteFile writes the contents of the output buffer to a file
func (pdfg *PDFGenerator) WriteFile(filename string) error {
	return ioutil.WriteFile(filename, pdfg.Bytes(), 0666)
}

//findPath finds the path to wkhtmltopdf by
//- first looking in the current dir
//- looking in the PATH and PATHEXT environment dirs
//- using the WKHTMLTOPDF_PATH environment dir
//The path is cached, meaning you can not change the location of wkhtmltopdf in
//a running program once it has been found
func (pdfg *PDFGenerator) findPath() error {
	const exe = "wkhtmltopdf"
	pdfg.binPath = GetPath()
	if pdfg.binPath != "" {
		// wkhtmltopdf has already already found, return
		return nil
	}
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	path, err := exec.LookPath(filepath.Join(exeDir, exe))
	if err == nil && path != "" {
		binPath.Set(path)
		pdfg.binPath = path
		return nil
	}
	path, err = exec.LookPath(exe)
	if err == nil && path != "" {
		binPath.Set(path)
		pdfg.binPath = path
		return nil
	}
	dir := os.Getenv("WKHTMLTOPDF_PATH")
	if dir == "" {
		return fmt.Errorf("%s not found", exe)
	}
	path, err = exec.LookPath(filepath.Join(dir, exe))
	if err == nil && path != "" {
		binPath.Set(path)
		pdfg.binPath = path
		return nil
	}
	return fmt.Errorf("%s not found", exe)
}

// Create creates the PDF document and stores it in the internal buffer if no error is returned
func (pdfg *PDFGenerator) Create() error {
	return pdfg.run()
}

func (pdfg *PDFGenerator) run() error {
	// create command
	cmd := exec.Command(pdfg.binPath, pdfg.Args()...)

	// set stderr to the provided writer, or create a new buffer
	var errBuf *bytes.Buffer
	cmd.Stderr = pdfg.stdErr
	if cmd.Stderr == nil {
		errBuf = new(bytes.Buffer)
		cmd.Stderr = errBuf
	}

	// set output to the desired writer or the internal buffer
	if pdfg.outWriter != nil {
		cmd.Stdout = pdfg.outWriter
	} else {
		cmd.Stdout = &pdfg.outbuf
	}

	// if there is a pageReader page (from Stdin) we set Stdin to that reader
	for _, page := range pdfg.pages {
		if page.Reader() != nil {
			cmd.Stdin = page.Reader()
			break
		}
	}

	// run cmd to create the PDF
	err := cmd.Run()
	if err != nil {
		// on an error, return the contents of Stderr if it was our own buffer
		// if Stderr was set to a custom writer, just return err
		if errBuf != nil {
			if errStr := errBuf.String(); strings.TrimSpace(errStr) != "" {
				return errors.New(errStr)
			}
		}
		return err
	}
	return nil
}

// NewPDFGenerator returns a new PDFGenerator struct with all options created and
// checks if wkhtmltopdf can be found on the system
func NewPDFGenerator() (*PDFGenerator, error) {
	pdfg := NewPDFPreparer()
	return pdfg, pdfg.findPath()
}

// NewPDFPreparer returns a PDFGenerator object without looking for the wkhtmltopdf executable file.
// This is useful to prepare a PDF file that is generated elsewhere and you just want to save the config as JSON.
// Note that Create() can not be called on this object unless you call SetPath yourself.
func NewPDFPreparer() *PDFGenerator {
	return &PDFGenerator{
		globalOptions:  newGlobalOptions(),
		outlineOptions: newOutlineOptions(),
		Cover: cover{
			pageOptions: newPageOptions(),
		},
		TOC: toc{
			allTocOptions: allTocOptions{
				tocOptions:  newTocOptions(),
				pageOptions: newPageOptions(),
			},
		},
	}
}
