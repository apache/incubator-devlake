/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/apache/incubator-devlake/utils"
	"github.com/panjf2000/ants/v2"
)

type WorkerScheduler struct {
	waitGroup    *sync.WaitGroup
	pool         *ants.Pool
	subPool      *ants.Pool
	ticker       *time.Ticker
	workerErrors *[]error
	ctx          context.Context
}

// NewWorkerScheduler 创建一个并行执行的调度器，控制最大运行数和每秒最大运行数量
// NewWorkerScheduler Create a parallel scheduler to control the maximum number of runs and the maximum number of runs per second
// 注意: task执行是无序的
// Warning: task execution is out of order
func NewWorkerScheduler(workerNum int, maxWork int, maxWorkDuration time.Duration, ctx context.Context, maxRetry int) (*WorkerScheduler, error) {
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
	subPool, err := ants.NewPool(workerNum*maxRetry, ants.WithPanicHandler(func(i interface{}) {
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
		subPool:      subPool,
		ticker:       ticker,
		workerErrors: pWorkerErrors,
		ctx:          ctx,
	}
	return scheduler, nil
}

func (s *WorkerScheduler) Submit(task func() error, pool ...*ants.Pool) error {
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}
	s.waitGroup.Add(1)
	// this is expensive, enable by EnvVar
	cf := "set Environment Varaible ASYNC_CF=true to enable callframes information"
	if os.Getenv("ASYNC_CF") == "true" {
		cf = utils.GatherCallFrames()
	}
	var currentPool *ants.Pool
	if pool == nil {
		currentPool = s.pool
	} else {
		currentPool = pool[0]
	}

	return currentPool.Submit(func() {
		var err error
		defer s.waitGroup.Done()
		defer func() {
			if err == nil {
				r := recover()
				if r != nil {
					err = fmt.Errorf("%s", r)
				}
			}
			if err != nil {
				panic(fmt.Errorf("%s\n%s", err, cf))
			}
		}()
		if pool == nil && s.ticker != nil {
			for s.subPool.Running() != 0 {
				<-s.ticker.C
			}
		}
		if s.ticker != nil {
			<-s.ticker.C
		}
		select {
		case <-s.ctx.Done():
			err = s.ctx.Err()
		default:
			err = task()
		}
	})
}

func (s *WorkerScheduler) WaitUntilFinish() error {
	s.waitGroup.Wait()
	if s.workerErrors != nil && len(*s.workerErrors) > 0 {
		return fmt.Errorf("%s", *s.workerErrors)
	}
	return nil
}

func (s *WorkerScheduler) Release() {
	s.pool.Release()
	s.subPool.Release()
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

func (s *WorkerScheduler) Add(delta int) {
	s.waitGroup.Add(delta)
}

func (s *WorkerScheduler) Done() {
	s.waitGroup.Done()
}
