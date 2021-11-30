package redis

import (
	"fmt"
	"sync"
	"time"
)

// Counter is a counter interface.
type Counter interface {
	Add(time.Time, int64)
	Reset()
	Value(time.Time) int64
}

// Group is a counter group.
type Group struct {
	mu   sync.RWMutex
	vecs map[string]Counter

	// New optionally specifies a function to generate a counter.
	// It may not be changed concurrently with calls to other functions.
	New func(string, time.Time) Counter
}

// Add add a counter by a specified key, if counter not exists then make a new one and return new value.
func (g *Group) Add(access time.Time, key string, value int64) {

	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()

	if !ok {

		k := fmt.Sprintf("rolling_%s", key)
		vec = g.New(k, access)
		g.mu.Lock()
		if g.vecs == nil {
			g.vecs = make(map[string]Counter)
		}
		if _, ok = g.vecs[key]; !ok {
			g.vecs[key] = vec
		}
		g.mu.Unlock()
	}

	vec.Add(access, value)

}

// Value get a counter value by key.
func (g *Group) Value(access time.Time, key string) int64 {
	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()
	if ok {
		return vec.Value(access)
	}
	return 0
}

// Reset reset a counter by key.
func (g *Group) Reset(key string) {
	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()
	if ok {
		vec.Reset()
	}
}
