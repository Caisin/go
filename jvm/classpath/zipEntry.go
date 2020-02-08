package classpath

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"path/filepath"
)

type ZipEntry struct {
	absPath string
	zipRC   *zip.ReadCloser
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return &ZipEntry{absPath, nil}
}

func (e *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	if e.zipRC == nil {
		err := e.openJar()
		if err != nil {
			return nil, nil, err
		}
	}

	classFile := e.findClass(className)
	if classFile == nil {
		return nil, nil, errors.New("class not found: " + className)
	}

	data, err := readClass(classFile)
	return data, e, err
}

// todo: close zip
func (e *ZipEntry) openJar() error {
	r, err := zip.OpenReader(e.absPath)
	if err == nil {
		e.zipRC = r
	}
	return err
}

func (e *ZipEntry) findClass(className string) *zip.File {
	for _, f := range e.zipRC.File {
		if f.Name == className {
			return f
		}
	}
	return nil
}

func readClass(classFile *zip.File) ([]byte, error) {
	rc, err := classFile.Open()
	if err != nil {
		return nil, err
	}
	// read class data
	data, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *ZipEntry) String() string {
	return e.absPath
}
