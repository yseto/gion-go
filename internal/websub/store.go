package websub

import (
	"slices"
	"sync"
)

type mutex struct {
	urls []string
	sync.RWMutex
}

var topicURLs mutex

func init() {
	topicURLs = mutex{}
}

func (m *mutex) insert(url string) bool {
	m.Lock()
	defer m.Unlock()
	if slices.Contains(m.urls, url) {
		return false
	}
	m.urls = append(m.urls, url)
	return true
}

func (m *mutex) contains(url string) bool {
	m.RLock()
	defer m.RUnlock()
	return slices.Contains(m.urls, url)
}

func (m *mutex) delete(url string) bool {
	m.Lock()
	defer m.Unlock()

	if index := slices.Index(m.urls, url); index >= 0 {
		m.urls = slices.Delete(m.urls, index, index+1)
		return true
	}
	return false
}
