package layer

import (
	"context"
	"fmt"
	"net"
	"sort"
)

var m = make(map[string]CreatorFunc)

func Register(name string, ds CreatorFunc) {
	m[name] = ds
}

func Layers() []string {
	rs := make([]string, 0, len(m))
	for name := range m {
		rs = append(rs, name)
	}
	sort.Strings(rs)
	return rs
}

func MustGet(name string) CreatorFunc {
	v, ok := Get(name)
	if !ok {
		panic(fmt.Sprintf("not found layer:%s", name))
	}
	return v
}

func Get(name string) (CreatorFunc, bool) {
	v, ok := m[name]
	if !ok {
		return nil, false
	}
	return v, true
}

type CreatorFunc func(params interface{}) (ILayer, error)

type ILayer interface {
	Name() string
	MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error)
}

var Default = &defaultLayer{}

type defaultLayer struct {
}

func (d *defaultLayer) Name() string {
	return "default"
}

func (d *defaultLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	return conn, nil
}
