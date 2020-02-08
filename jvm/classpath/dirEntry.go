package classpath

import (
	"io/ioutil"
	"path/filepath"
)

type DirEntry struct{ absDir string }

func newDirEntry(path string) *DirEntry {
	absDir, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &DirEntry{absDir}
}
func (dir *DirEntry) readClass(className string) ([]byte, Entry, error) {
	fileName := filepath.Join(dir.absDir, className)
	data, err := ioutil.ReadFile(fileName)
	return data, dir, err
}
func (dir *DirEntry) String() string {
	return dir.absDir
}
