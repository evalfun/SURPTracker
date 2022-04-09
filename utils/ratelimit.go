package utils

import (
	"sync"
	"time"
)

type SimpleRateLimiter struct {
	resetDuration int
	rateLimit     int
	record        sync.Map
	closed        bool
}

func NewSimpleRateLimiter(resetDuration int, rateLimit int) *SimpleRateLimiter {
	this := &SimpleRateLimiter{
		resetDuration: resetDuration,
		rateLimit:     rateLimit,
	}
	go this.resetRecordProc()
	return this
}

func (this *SimpleRateLimiter) BelowRate(id string) bool {
	_rate, ok := this.record.Load(id)
	if !ok {
		this.record.Store(id, 0)
		return true
	}
	rate := _rate.(int)
	this.record.Store(id, rate+1)
	if rate > this.rateLimit {
		return false
	}
	return true
}

func (this *SimpleRateLimiter) resetRecordProc() {
	for this.closed == false {
		if this.resetDuration <= 1 {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
			if time.Now().Unix()%int64(this.resetDuration) != 0 {
				continue
			}
		}
		this.record.Range(func(key, value interface{}) bool {
			this.record.Delete(key)
			return true
		})
	}
}
