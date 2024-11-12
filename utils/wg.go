package utils

import (
	"sync"
	"sync/atomic"
)

type CustomWaitGroup struct {
	wg    sync.WaitGroup
	count int32 // Counter to track the number of active tasks
}

func (cwg *CustomWaitGroup) Add(delta int) {
	atomic.AddInt32(&cwg.count, int32(delta))
	cwg.wg.Add(delta)
}

func (cwg *CustomWaitGroup) Done() {
	atomic.AddInt32(&cwg.count, -1)
	cwg.wg.Done()
}

func (cwg *CustomWaitGroup) Wait() {
	cwg.wg.Wait()
}

func (cwg *CustomWaitGroup) GetCount() int32 {
	return atomic.LoadInt32(&cwg.count)
}
