package wkhtmltopdf

import "sync"

//the cached mutexed path as used by findPath()
type stringStore struct {
	val string
	sync.Mutex
}

func (ss *stringStore) Get() string {
	ss.Lock()
	defer ss.Unlock()
	return ss.val
}

func (ss *stringStore) Set(s string) {
	ss.Lock()
	ss.val = s
	ss.Unlock()
}
