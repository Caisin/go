package classpath

import "errors"
import "strings"

type CompositeEntry []Entry

func newCompositeEntry(pathList string) CompositeEntry {
	var entries []Entry
	for _, path := range strings.Split(pathList, pathListSeparator) {
		e := newEntry(path)
		entries = append(entries, e)
	}
	return entries
}
func (composite CompositeEntry) readClass(className string) ([]byte, Entry, error) {
	for _, ce := range composite {
		data, from, err := ce.readClass(className)
		if err == nil {
			return data, from, nil
		}
	}
	return nil, nil, errors.New("class not found: " + className)
}
func (composite CompositeEntry) String() string {
	strs := make([]string, len(composite))
	for i, e := range composite {
		strs[i] = e.String()
	}
	return strings.Join(strs, pathListSeparator)
}
