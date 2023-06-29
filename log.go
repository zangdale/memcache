package memcache

type LogInter interface {
	Printf(format string, v ...any)
	Println(v ...any)
}
