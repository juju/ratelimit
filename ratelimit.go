// The ratelimit package provides an efficient token bucket implementation.
// See http://en.wikipedia.org/wiki/Token_bucket.
package ratelimit

import (
	"sync"
	"time"
)

// TODO what about aborting requests?

// TokenBucket represents a token bucket
// that fills at a predetermined rate.
// Methods on TokenBucket may be called
// concurrently.
type TokenBucket struct {
	mu           sync.Mutex
	startTime    time.Time
	capacity     int64
	fillInterval time.Duration
	availTick    int64
	avail        int64
}

// New returns a new token bucket that fills at the
// rate of one token every fillInterval, up to the given
// maximum capacity. Both arguments must be
// positive.
func New(fillInterval time.Duration, capacity int64) *TokenBucket {
	if fillInterval <= 0 {
		panic("token bucket fill interval is not > 0")
	}
	if capacity <= 0 {
		panic("token bucket capacity is not > 0")
	}
	return &TokenBucket{
		startTime:    time.Now(),
		capacity:     capacity,
		fillInterval: fillInterval,
	}
}

// Get gets count tokens from the bucket, waiting
// until the tokens are available.
func (tb *TokenBucket) Get(count int64) {
	if d := tb.GetNB(count); d > 0 {
		time.Sleep(d)
	}
}

// GetNB gets count tokens from the bucket without
// blocking. It returns the time to wait until the
// tokens are actually available.
//
// Note that if the request is irrevocable - there
// is no way to return tokens to the bucket once
// this method commits us to taking them.
func (tb *TokenBucket) GetNB(count int64) time.Duration {
	return tb.getNB(time.Now(), count)
}

// getNB is the internal version of GetNB - it takes
// the current time as an argument to enable easy testing.
func (tb *TokenBucket) getNB(now time.Time, count int64) time.Duration {
	if count <= 0 {
		return 0
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	currentTick := int64(now.Sub(tb.startTime) / tb.fillInterval)
	tb.adjust(currentTick)

	tb.avail -= count
	if tb.avail >= 0 {
		return 0
	}
	endTick := currentTick-tb.avail
	endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
	return endTime.Sub(now)
}

// adjust adjusts the current bucket capacity based
// on the current tick.
func (tb *TokenBucket) adjust(currentTick int64) {
	if tb.avail >= tb.capacity {
		return
	}
	tb.avail += currentTick - tb.availTick
	if tb.avail > tb.capacity {
		tb.avail = tb.capacity
	}
	tb.availTick = currentTick
}
