package config

import (
	"sync"
)

type Params struct {
	params map[string]interface{}
	sync.Mutex
}

func (c *Params) GetParam(paramName string) interface{} {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.params[paramName]
}

func (c *Params) SetParam(paramName string, value interface{}) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.params[paramName] = value
}
