package GxMisc

import (
	"sync"
)

type Counter struct {
	C     uint32
	Mutex *sync.Mutex
}

func NewCounter() *Counter {
	counter := new(Counter)
	counter.C = 0
	counter.Mutex = new(sync.Mutex)
	return counter
}

func (counter *Counter) Genarate() uint32 {
	counter.Mutex.Lock()
	defer counter.Mutex.Unlock()

	counter.C += 1
	return counter.C
}
