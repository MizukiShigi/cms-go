package context

import (
	"context"
	"sync"
)

func WithValue(parent context.Context, key string, val any) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if v, ok := parent.Value(Logging).(*sync.Map); ok {
		mapCopy := copySyncMap(v)
		mapCopy.Store(key, val)
		return context.WithValue(parent, Logging, mapCopy)
	}
	v := &sync.Map{}
	v.Store(key, val)
	return context.WithValue(parent, Logging, v)
}

func copySyncMap(m *sync.Map) *sync.Map {
	var cp sync.Map
	m.Range(func(k, v any) bool {
		cp.Store(k, v)
		return true
	})
	return &cp
}
