package utils

import (
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

// Pool accepts the tasks from client, it limits the total of goroutines to a given number by recycling goroutines.
type WorkerScheduler struct {
	waitGroup    *sync.WaitGroup
	pool         *ants.Pool
	ticker       *time.Ticker
	workerErrors *[]error
}

// 创建一个并行执行的调度器，控制最大运行数和每秒最大运行数量
// Create a parallel scheduler to control the maximum number of runs and the maximum number of runs per second
// 注意: task执行是无序的
// Warning: task execution is out of order
func NewWorkerScheduler(workerNum int, maxWorkEverySeconds int) (*WorkerScheduler, error) {
	var waitGroup sync.WaitGroup
	workerErrors := make([]error, 0)
	pWorkerErrors := &workerErrors
	pool, err := ants.NewPool(workerNum, ants.WithPanicHandler(func(i interface{}) {
		workerErrors = append(*pWorkerErrors, i.(error))
		pWorkerErrors = &workerErrors
		waitGroup.Done()
	}))
	if err != nil {
		return nil, err
	}
	scheduler := &WorkerScheduler{
		waitGroup:    &waitGroup,
		pool:         pool,
		ticker:       time.NewTicker(time.Second / time.Duration(maxWorkEverySeconds)),
		workerErrors: pWorkerErrors,
	}
	return scheduler, nil
}

func (s WorkerScheduler) Submit(task func() error) error {
	s.waitGroup.Add(1)
	return s.pool.Submit(func() {
		<-s.ticker.C
		err := task()
		if err != nil {
			panic(err)
		}
		s.waitGroup.Done()
	})
}

func (s WorkerScheduler) WaitUntilFinish() {
	s.waitGroup.Wait()
}

func (s WorkerScheduler) Release() {
	s.pool.Release()
	s.ticker.Stop()
}
