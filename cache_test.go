package memcache

import (
	"testing"
)

func TestNew(t *testing.T) {
	got := New(3)
	got.Add(&Value{Key: "1", Value: 1})
	t.Logf("%v", got.Get("2"))
	got.Add(&Value{Key: "2", Value: 2})
	t.Logf("%v", got.Get("2"))
	got.Add(&Value{Key: "3", Value: 3})
	t.Logf("%v", got.Get("1"))
	got.Add(&Value{Key: "4", Value: 4})
	t.Logf("%v", got.Get("1"))
	t.Logf("%v", got.Get("4"))
	got.Add(&Value{Key: "5", Value: 5})
	t.Logf("%v", got.Get("5"))
	got.Add(&Value{Key: "5", Value: 55})
	t.Logf("%v", got.Get("5"))
	t.Logf("%v", got.Get("4"))
	t.Logf("%v", got.Get("3"))
	t.Logf("%v", got.Get("2"))
	t.Logf("%v", got.Get("1"))
}
