package safecounter

import "sync"

type SafeCounter struct {
	uint
	sync.RWMutex
}

func (sc *SafeCounter) Increment() {
	sc.Lock()
	sc.uint++
	sc.Unlock()
}

func (sc *SafeCounter) Decrement() {
	sc.Lock()
	sc.uint--
	sc.Unlock()
}

func (sc *SafeCounter) Uint() uint {
	sc.RLock()
	u := sc.uint
	sc.RUnlock()
	return u
}

// Reset resets the counter count value to 0 and returns the value that count
// was on immediately prior to being set to 0
func (sc *SafeCounter) Reset() uint {
	sc.Lock()
	u := sc.uint
	sc.uint = 0
	sc.Unlock()
	return u
}
