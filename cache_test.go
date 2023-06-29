package memcache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	got := New(NewCountRole(3), NewTimeOutRole(time.Millisecond))
	got.Add("1", &Value{Data: 1})
	t.Logf("%v", got.Get("2"))
	got.Add("2", &Value{Data: 2})
	t.Logf("%v", got.Get("2"))
	got.Add("3", &Value{Data: 3})
	t.Logf("%v", got.Get("1"))
	got.Add("4", &Value{Data: 4})
	t.Logf("%v", got.Get("1"))
	t.Logf("%v", got.Get("4"))
	got.Add("5", &Value{Data: 5})
	t.Logf("%v", got.Get("5"))
	got.Add("5", &Value{Data: 55})
	t.Logf("%v", got.Get("5"))
	t.Logf("%v", got.Get("4"))
	t.Logf("%v", got.Get("3"))
	t.Logf("%v", got.Get("2"))
	t.Logf("%v", got.Get("1"))
}
