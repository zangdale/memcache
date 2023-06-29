package memcache

import (
	"time"
)

type RemoveRoleInter interface {
	Check(length int, t time.Time) (remove bool)
}

var DefaultRole RemoveRoleInter = NewCountRole(10)

func NewCountRole(count int) RemoveRoleInter {
	if count < 0 {
		count = 10
	}
	return &countRole{
		count: count,
	}
}

type countRole struct {
	count int
}

func (r *countRole) Check(length int, t time.Time) (remove bool) {
	return length > r.count
}

func NewTimeOutRole(timeout time.Duration) RemoveRoleInter {
	if timeout < 0 {
		timeout = 15 * time.Minute
	}
	return &timeoutRole{
		timeout: timeout,
	}
}

type timeoutRole struct {
	timeout time.Duration
}

func (r *timeoutRole) Check(length int, t time.Time) (remove bool) {
	return time.Since(t) > r.timeout
}

func MultiRole(roles ...RemoveRoleInter) RemoveRoleInter {
	return &multiRole{
		roles: roles,
	}
}

type multiRole struct {
	roles []RemoveRoleInter
}

func (r *multiRole) Check(length int, t time.Time) (remove bool) {
	for _, v := range r.roles {
		if v.Check(length, t) {
			return true
		}
	}
	return false
}
