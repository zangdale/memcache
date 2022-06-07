package memcache

import (
	"container/list"
	"sync"
)

type cache struct {
	dataMap map[string]*list.Element
	list    *list.List
	count   int
	sy      sync.Mutex

	UpMoveToBack  bool
	GetMoveToBack bool
}

type Value struct {
	Key   string
	Value interface{}
}

func New(count int) *cache {
	if count < 0 {
		count = 10
	}
	return &cache{
		dataMap: make(map[string]*list.Element, count),
		list:    list.New(),
		count:   count,
		sy:      sync.Mutex{},

		UpMoveToBack:  true,
		GetMoveToBack: false,
	}
}

func (c *cache) Add(v *Value) {
	c.sy.Lock()
	defer c.sy.Unlock()

	key := v.Key
	if old, ok := c.dataMap[key]; ok {
		// have
		if c.UpMoveToBack {
			c.list.MoveToBack(old)
		}

		old.Value = v
		c.dataMap[key] = old
		return
	}
	//no have
	e := c.list.PushBack(v)
	c.dataMap[key] = e

	c.checkLen()
}

func (c *cache) Get(key string) *Value {
	c.sy.Lock()
	defer c.sy.Unlock()

	if old, ok := c.dataMap[key]; ok {
		// have
		if c.GetMoveToBack {
			c.list.MoveToBack(old)
		}
		return old.Value.(*Value)
	}
	return &Value{
		Key:   key,
		Value: nil,
	}
}

func (c *cache) checkLen() {
	if c.list.Len() > c.count {
		d := c.list.Front()
		dv := d.Value.(*Value)
		c.list.Remove(d)
		delete(c.dataMap, dv.Key)
	}
}
