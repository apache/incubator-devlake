package helper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/merico-dev/lake/utils"
	"github.com/panjf2000/ants/v2"
)

type WorkerScheduler struct {
	waitGroup    *sync.WaitGroup
	pool         *ants.Pool
	ticker       *time.Ticker
	workerErrors *[]error
	ctx          context.Context
}

// NewWorkerScheduler 创建一个并行执行的调度器，控制最大运行数和每秒最大运行数量
// NewWorkerScheduler Create a parallel scheduler to control the maximum number of runs and the maximum number of runs per second
// 注意: task执行是无序的
// Warning: task execution is out of order
func NewWorkerScheduler(workerNum int, maxWork int, maxWorkDuration time.Duration, ctx context.Context) (*WorkerScheduler, error) {
	var waitGroup sync.WaitGroup
	workerErrors := make([]error, 0)
	pWorkerErrors := &workerErrors
	var mux sync.Mutex
	pool, err := ants.NewPool(workerNum, ants.WithPanicHandler(func(i interface{}) {
		mux.Lock()
		defer mux.Unlock()
		workerErrors = append(*pWorkerErrors, i.(error))
		pWorkerErrors = &workerErrors
	}))
	if err != nil {
		return nil, err
	}
	var ticker *time.Ticker
	if maxWork > 0 {
		ticker = time.NewTicker(maxWorkDuration / time.Duration(maxWork))
	}
	scheduler := &WorkerScheduler{
		waitGroup:    &waitGroup,
		pool:         pool,
		ticker:       ticker,
		workerErrors: pWorkerErrors,
		ctx:          ctx,
	}
	return scheduler, nil
}

func (s *WorkerScheduler) Submit(task func() error) error {
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}
	s.waitGroup.Add(1)
	return s.pool.Submit(func() {
		defer s.waitGroup.Done()
		defer func() {
			r := recover()
			if r != nil {
				panic(fmt.Errorf("%s\n%s", r, utils.GatherCallFrames()))
			}
		}()
		select {
		case <-s.ctx.Done():
			panic(s.ctx.Err())
		default:
		}
		if s.ticker != nil {
			<-s.ticker.C
		}
		err := task()
		if err != nil {
			panic(fmt.Errorf("%w\n%s", err, utils.GatherCallFrames()))
		}
	})
}

func (s *WorkerScheduler) WaitUntilFinish() {
	s.waitGroup.Wait()
	if s.workerErrors != nil && len(*s.workerErrors) > 0 {
		panic(fmt.Errorf("%s", *s.workerErrors))
	}
}

func (s *WorkerScheduler) Release() {
	s.pool.Release()
	if s.ticker != nil {
		s.ticker.Stop()
	}
}
