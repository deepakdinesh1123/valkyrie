package concurrency

import (
	"sync"
	"sync/atomic"
)

type SafeWaitGroup struct {
	sync.WaitGroup
	counter int64
}

func (sgw *SafeWaitGroup) Add(delta int) {
	atomic.AddInt64(&sgw.counter, int64(delta))
	sgw.WaitGroup.Add(delta)
}

func (sgw *SafeWaitGroup) Done() {
	atomic.AddInt64(&sgw.counter, -1)
	sgw.WaitGroup.Done()
}

func (sgw *SafeWaitGroup) Wait() {
	sgw.WaitGroup.Wait()
}

func (sgw *SafeWaitGroup) Count() int64 {
	return atomic.LoadInt64(&sgw.counter)
}
