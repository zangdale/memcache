package memcache

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type cache struct {
	dataMap    map[string]*list.Element
	list       *list.List
	removeRole RemoveRoleInter
	sy         sync.Mutex

	/////
	logger LogInter
	/////
	upMoveToBack  bool
	getMoveToBack bool
}

func New(role ...RemoveRoleInter) *cache {
	var removeRole RemoveRoleInter
	if len(role) == 0 {
		removeRole = MultiRole(DefaultRole)
	} else {
		removeRole = MultiRole(role...)
	}
	return &cache{
		dataMap:    make(map[string]*list.Element),
		list:       list.New(),
		removeRole: removeRole,
		sy:         sync.Mutex{},

		logger: log.Default(),
		// 更新移动到最新
		upMoveToBack: true,
		// 获取一次移动数据到最新
		getMoveToBack: false,
	}
}

// 获取一次移动数据到最新
func (c *cache) WithGetMoveToBack() *cache {
	c.getMoveToBack = true
	return c
}

// 获取一次移动数据到最新
func (c *cache) WithLog(lg LogInter) *cache {
	c.logger = lg
	return c
}

// 更新移动到最新
func (c *cache) WithUpMoveToBack() *cache {
	c.upMoveToBack = true
	return c
}

func (c *cache) Add(key string, v *Value) {
	c.sy.Lock()
	defer c.sy.Unlock()

	v.t = time.Now()
	v.key = key

	if old, ok := c.dataMap[key]; ok {
		// have
		if c.upMoveToBack {
			c.list.MoveToBack(old)
		} else {
			v.t = old.Value.(*Value).t
		}

		old.Value = v
		c.dataMap[key] = old

		c.logger.Printf("memcache add rewrite: %s", v)
		return
	}
	//no have
	e := c.list.PushBack(v)
	c.dataMap[key] = e

	c.logger.Printf("memcache add: %s", v)
	c.checkRemove()
}

func (c *cache) Get(key string) *Value {
	c.sy.Lock()
	defer c.sy.Unlock()

	var v *Value = new(Value)
	var t time.Time
	if old, ok := c.dataMap[key]; ok {
		// have
		v = old.Value.(*Value)
		t = v.t
		if c.getMoveToBack {
			v.t = time.Now()
			old.Value = v
			c.list.MoveToBack(old)
		}
	}

	c.logger.Printf("memcache get: %s", v)

	return &Value{
		key:  key,
		t:    t,
		Data: v.Data,
	}
}

func (c *cache) checkRemove() {
	d := c.list.Front()
	dv := d.Value.(*Value)

	for c.removeRole.Check(c.list.Len(), dv.t) {
		c.list.Remove(d)
		delete(c.dataMap, dv.key)
		c.logger.Printf("memcache delete: %s", dv)

		d = c.list.Front()
		dv = d.Value.(*Value)
	}
}
