package concurrency

import (
	"sync"
	"sync/atomic"
)

type SafeWaitGroup struct {
	sync.WaitGroup
	counter int32
}

func (sgw *SafeWaitGroup) Add(delta int32) {
	atomic.AddInt32(&sgw.counter, delta)
	sgw.WaitGroup.Add(int(delta))
}

func (sgw *SafeWaitGroup) Done() {
	atomic.AddInt32(&sgw.counter, -1)
	sgw.WaitGroup.Done()
}

func (sgw *SafeWaitGroup) Wait() {
	sgw.WaitGroup.Wait()
}

func (sgw *SafeWaitGroup) Count() int32 {
	return atomic.LoadInt32(&sgw.counter)
}
