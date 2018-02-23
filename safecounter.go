package main

import "sync"

type safeCounter struct {
	uint
	sync.RWMutex
}

func (sc *safeCounter) Increment() {
	sc.Lock()
	sc.uint++
	sc.Unlock()
}

func (sc *safeCounter) Decrement() {
	sc.Lock()
	sc.uint--
	sc.Unlock()
}

func (sc *safeCounter) Uint() uint {
	sc.RLock()
	u := sc.uint
	sc.RUnlock()
	return u
}

// Reset resets the counter count value to 0 and returns the value that count
// was on immediately prior to being set to 0
func (sc *safeCounter) Reset() uint {
	sc.Lock()
	u := sc.uint
	sc.uint = 0
	sc.Unlock()
	return u
}
