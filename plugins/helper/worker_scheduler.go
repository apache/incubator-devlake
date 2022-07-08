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
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	ants "github.com/panjf2000/ants/v2"
)

// WorkerScheduler runs asynchronous tasks in parallel with throttling support
type WorkerScheduler struct {
	waitGroup    sync.WaitGroup
	pool         *ants.Pool
	ticker       *time.Ticker
	workerErrors []error
	ctx          context.Context
	mu           sync.Mutex
	counter      int32
	logger       core.Logger
}

//var callframeEnabled = os.Getenv("ASYNC_CF") == "true"

// NewWorkerScheduler creates a WorkerScheduler
func NewWorkerScheduler(
	ctx context.Context,
	workerNum int,
	maxWork int,
	maxWorkDuration time.Duration,
	maxRetry int,
	logger core.Logger,
) (*WorkerScheduler, error) {
	if maxWork <= 0 {
		return nil, fmt.Errorf("maxWork less than 1")
	}
	if maxWorkDuration <= 0 {
		return nil, fmt.Errorf("maxWorkDuration less than 1")
	}
	s := &WorkerScheduler{
		ctx:    ctx,
		ticker: time.NewTicker(maxWorkDuration / time.Duration(maxWork)),
		logger: logger,
	}
	pool, err := ants.NewPool(workerNum, ants.WithPanicHandler(func(i interface{}) {
		s.checkError(i)
	}))
	if err != nil {
		return nil, err
	}
	s.pool = pool
	return s, nil
}

// SubmitBlocking enqueues a async task to ants, the task will be executed in future when timing is right.
// It doesn't return error because it wouldn't be any when with a Blocking semantic, returned error does nothing but
// causing confusion, more often, people thought it is returned by the task.
// Since it is async task, the callframes would not be available for production mode, you can export Environment
// Varaible ASYNC_CF=true to enable callframes capturing when debugging.
// IMPORTANT: do NOT call SubmitBlocking inside the async task, it is likely to cause a deadlock, call
// SubmitNonBlocking instead when number of tasks is relatively small.
func (s *WorkerScheduler) SubmitBlocking(task func() error) {
	if s.HasError() {
		return
	}
	s.waitGroup.Add(1)
	s.checkError(s.pool.Submit(func() {
		defer s.waitGroup.Done()

		id := atomic.AddInt32(&s.counter, 1)
		s.logger.Debug("schedulerJob >>> %d started", id)
		defer s.logger.Debug("schedulerJob <<< %d ended", id)

		if s.HasError() {
			return
		}

		// normal error
		select {
		case <-s.ctx.Done():
			panic(s.ctx.Err())
		case <-s.ticker.C:
			err := task()
			if err != nil {
				panic(err)
			}
		}
	}))
}

/*
func (s *WorkerScheduler) gatherCallFrames() string {
	cf := "set Environment Varaible ASYNC_CF=true to enable callframes capturing"
	if callframeEnabled {
		cf = utils.GatherCallFrames(1)
	}
	return cf
}
*/

func (s *WorkerScheduler) appendError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workerErrors = append(s.workerErrors, err)
}

func (s *WorkerScheduler) checkError(err interface{}) {
	if err == nil {
		return
	}
	s.appendError(err.(error))
}

// HasError return if any error occurred
func (s *WorkerScheduler) HasError() bool {
	return len(s.workerErrors) > 0
}

// NextTick enqueues task in a NonBlocking manner, you should only call this method within task submitted by
// SubmitBlocking method
// IMPORTANT: do NOT call this method with a huge number of tasks, it is likely to eat up all available memory
func (s *WorkerScheduler) NextTick(task func() error) {
	// to make sure task will be enqueued
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.checkError(task())
	}()
}

// Wait blocks current go-routine until all workers returned
func (s *WorkerScheduler) Wait() error {
	s.waitGroup.Wait()
	if len(s.workerErrors) > 0 {
		return fmt.Errorf("%s", s.workerErrors)
	}
	return nil
}

// Release resources
func (s *WorkerScheduler) Release() {
	s.waitGroup.Wait()
	s.pool.Release()
	if s.ticker != nil {
		s.ticker.Stop()
	}
}
