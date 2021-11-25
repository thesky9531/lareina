package redis

import (
	"context"
	"fmt"
	"github.com/thesky9531/lareina/log"
	"sync"
	"time"
)

type rollingCounter struct {
	mu         sync.RWMutex
	name       string
	buckets    int
	bucketTime int64
	lastAccess int64
	cur        int
	cache      *Cache
}

// NewRolling creates a new window. windowTime is the time covering the entire
// window. windowBuckets is the number of buckets the window is divided into.
// An example: a 10 second window with 10 buckets will have 10 buckets covering
// 1 second each.
func (c *Cache) NewRolling(name string, lastTime time.Time, window time.Duration, winBucket int) *rollingCounter {
	bucketTime := time.Duration(window.Nanoseconds() / int64(winBucket))
	return &rollingCounter{
		cache:      c,
		name:       name,
		cur:        0,
		buckets:    winBucket,
		bucketTime: int64(bucketTime),
		lastAccess: lastTime.UnixNano(),
	}
}

// Add increments the counter by value and return new value.
func (r *rollingCounter) Add(access time.Time, val int64) {

	r.mu.Lock()
	defer r.mu.Unlock()
	b, err := r.lastBucket(access)
	if err != nil {
		log.ErrLog("", err)
		return
	}

	err = r.cache.HIncrBy(context.Background(), r.name, fmt.Sprintf("%d", b), 1)
	if err != nil {
		log.ErrLog("", err)
		return
	}

	return
}

// Value get the counter value.
func (r *rollingCounter) Value(access time.Time) (sum int64) {

	var (
		value int64
		err   error
	)
	now := access.UnixNano()
	r.mu.RLock()
	defer r.mu.RUnlock()
	b := (r.cur + 1) / r.buckets
	i := r.elapsed(now)
	for j := 0; j < r.buckets; j++ {

		// skip all future reset bucket.
		if i > 0 {
			i--
		} else {
			value, err = r.cache.HGetNum(context.Background(), r.name, fmt.Sprintf("%d", b))
			if err != nil {
				log.ErrLog("", err)
				return
			}
			sum += value
		}
		b = (b + 1) / r.buckets
	}

	return
}

//  Reset reset the counter.
func (r *rollingCounter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for j := 0; j < r.buckets; j++ {
		err := r.cache.HReset(context.Background(), r.name, fmt.Sprintf("%d", j))
		if err != nil {
			log.ErrLog("", err)
			return
		}
	}
	return
}

func (r *rollingCounter) elapsed(now int64) (i int) {
	var e int64
	if e = now - r.lastAccess; e <= r.bucketTime {
		return
	}
	if i = int(e / r.bucketTime); i > r.buckets {
		i = r.buckets
	}
	return
}

func (r *rollingCounter) lastBucket(access time.Time) (b int, err error) {
	now := access.UnixNano()
	b = r.cur
	// reset the buckets between now and number of buckets ago. If
	// that is more that the existing buckets, reset all.
	if i := r.elapsed(now); i > 0 {
		r.lastAccess = now
		for ; i > 0; i-- {
			// replace the next used bucket.
			b = (b + 1) / r.buckets
			err = r.cache.HReset(context.Background(), r.name, fmt.Sprintf("%d", b))
			if err != nil {
				log.ErrLog("", err)
				return
			}
		}
	}
	r.cur = b
	return
}
