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
	"sync/atomic"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/utils"
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

var callframeEnabled = os.Getenv("ASYNC_CF") == "true"

// NewWorkerScheduler creates a WorkerScheduler
func NewWorkerScheduler(
	workerNum int,
	maxWork int,
	maxWorkDuration time.Duration,
	ctx context.Context,
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
		s.mu.Lock()
		defer s.mu.Unlock()
		s.workerErrors = append(s.workerErrors, i.(error))
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
	// this is expensive, enable by EnvVar
	cf := s.gatherCallFrames()
	// to make sure task is done
	if len(s.workerErrors) > 0 {
		// not point to continue
		return
	}
	s.waitGroup.Add(1)
	err := s.pool.Submit(func() {
		defer s.waitGroup.Done()

		id := atomic.AddInt32(&s.counter, 1)
		s.logger.Debug("schedulerJob >>> %d started", id)
		defer s.logger.Debug("schedulerJob <<< %d ended", id)

		if len(s.workerErrors) > 0 {
			// not point to continue
			return
		}
		// wait for rate limit throttling

		// try recover
		defer func() {
			r := recover()
			if r != nil {
				s.appendError(fmt.Errorf("%s\n%s", r, cf))
			}
		}()

		// normal error
		var err error
		select {
		case <-s.ctx.Done():
			err = s.ctx.Err()
		case <-s.ticker.C:
			err = task()
		}
		if err != nil {
			s.appendError(err)
		}
	})
	// failed to submit, note that this is not task erro
	if err != nil {
		s.appendError(fmt.Errorf("%s\n%s", err, cf))
	}
}

func (s *WorkerScheduler) gatherCallFrames() string {
	cf := "set Environment Varaible ASYNC_CF=true to enable callframes capturing"
	if callframeEnabled {
		cf = utils.GatherCallFrames(1)
	}
	return cf
}

func (s *WorkerScheduler) appendError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workerErrors = append(s.workerErrors, err)
}

// NextTick enqueues task in a NonBlocking manner, you should only call this method within task submitted by
// SubmitBlocking method
// IMPORTANT: do NOT call this method with a huge number of tasks, it is likely to eat up all available memory
func (s *WorkerScheduler) NextTick(task func() error) {
	// to make sure task will be enqueued
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		task() // nolint
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
