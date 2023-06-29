package memcache

import (
	"fmt"
	"time"
)

type Value struct {
	key  string
	t    time.Time
	Data any
}

func (v *Value) String() string {
	return fmt.Sprintf("%s %s %t", v.key, v.t.Format(time.DateTime), v.IsNil())
}

func (v *Value) Key() string {
	return v.key
}

func (v *Value) Time() time.Time {
	return v.t
}

func (v *Value) IsNil() bool {
	return v.Data == nil
}

func (v *Value) IsZero() bool {
	return v.t.IsZero()
}
