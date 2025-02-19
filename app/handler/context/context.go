package context

import "sync"

type Map map[string]interface{}

type Context struct {
	store Map
	lock  sync.RWMutex
}

func New() *Context {
	return &Context{}
}

func (c *Context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *Context) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}
