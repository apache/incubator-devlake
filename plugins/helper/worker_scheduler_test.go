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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerScheduler(t *testing.T) {
	testChannel := make(chan int, 100)
	ctx, cancel := context.WithCancel(context.Background())
	s, _ := NewWorkerScheduler(5, 2, 1*time.Second, ctx, 0)
	defer s.Release()
	for i := 1; i <= 5; i++ {
		t := i
		_ = s.Submit(func() error {
			testChannel <- t
			return nil
		})
	}
	time.Sleep(1200 * time.Millisecond)
	if len(testChannel) < 2 {
		t.Fatal(`worker not start`)
	}
	if len(testChannel) > 2 {
		t.Fatal(`worker run too fast`)
	}
	time.Sleep(time.Second)
	if len(testChannel) < 4 {
		t.Fatal(`worker not run after a second`)
	}
	if len(testChannel) > 4 {
		t.Fatal(`worker run too fast after a second`)
	}
	assert.Nil(t, s.WaitUntilFinish())
	if len(*s.workerErrors) != 0 {
		t.Fatal(`worker got panic`)
	}
	if len(testChannel) != 5 {
		t.Fatal(`worker not wait until finish`)
	}
	cancel()
}

func TestNewWorkerSchedulerWithoutSecond(t *testing.T) {
	testChannel := make(chan int, 100)
	ctx, cancel := context.WithCancel(context.Background())
	s, _ := NewWorkerScheduler(5, 0, 1*time.Second, ctx, 0)
	defer s.Release()
	for i := 1; i <= 5; i++ {
		t := i
		_ = s.Submit(func() error {
			testChannel <- t
			return nil
		})
	}
	time.Sleep(5 * time.Millisecond)
	if len(testChannel) != 5 {
		t.Fatal(`worker not finish`)
	}
	assert.Nil(t, s.WaitUntilFinish())
	if len(testChannel) != 5 {
		t.Fatal(`worker not finish`)
	}
	cancel()
}

/*
func TestNewWorkerSchedulerWithPanic(t *testing.T) {
	testChannel := make(chan int, 100)
	ctx, cancel := context.WithCancel(context.Background())
	s, _ := NewWorkerScheduler(1, 1, ctx)
	defer s.Release()
	_ = s.Submit(func() error {
		testChannel <- 1
		return errors.New(`error message`)
	})
	s.WaitUntilFinish()
	if len(*s.workerErrors) != 1 {
		t.Fatal(`worker not got panic`)
	}
	cancel()
}
*/
