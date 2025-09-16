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

package tasks

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/stretchr/testify/assert"
)

type MockLogger struct {
	debugCalls []string
}

func (m *MockLogger) Debug(args ...interface{}) {
	if len(args) >= 1 {
		if str, ok := args[0].(string); ok {
			m.debugCalls = append(m.debugCalls, str)
		}
	}
}

func (m *MockLogger) Info(args ...interface{})       {}
func (m *MockLogger) Warn(args ...interface{})       {}
func (m *MockLogger) Error(args ...interface{})      {}
func (m *MockLogger) Fatal(args ...interface{})      {}
func (m *MockLogger) Panic(args ...interface{})      {}
func (m *MockLogger) Nested(name string) interface{} { return m }

// MockTaskContext implements a simple task context for testing
type MockTaskContext struct {
	logger *MockLogger
}

func (m *MockTaskContext) GetLogger() interface{} {
	return m.logger
}

func TestJobCollector_ResponseParser_SkipsJobsWithNilStartedAt(t *testing.T) {
	// Create a mock task context
	mockLogger := &MockLogger{}
	mockTaskCtx := &MockTaskContext{logger: mockLogger}

	now := time.Now()

	// Check run with zero StartedAt should be skipped
	checkRunSkipped := GraphqlQueryCheckRun{
		Id:          "skipped-job-1",
		Name:        "Skipped Job",
		DatabaseId:  456,
		Status:      "completed",
		Conclusion:  "skipped",
		StartedAt:   nil,
		CompletedAt: &now,
	}

	checkRunNormal := GraphqlQueryCheckRun{
		Id:          "normal-job-1",
		Name:        "Normal Job",
		DatabaseId:  789,
		Status:      "completed",
		Conclusion:  "success",
		StartedAt:   &now,
		CompletedAt: &now,
	}

	checkRuns := []GraphqlQueryCheckRun{checkRunSkipped, checkRunNormal}

	responseParser := func(checkRuns []GraphqlQueryCheckRun, runId int) (messages []json.RawMessage, err errors.Error) {
		for _, checkRun := range checkRuns {
			dbCheckRun := &DbCheckRun{
				RunId:                runId,
				GraphqlQueryCheckRun: &checkRun,
			}
			if dbCheckRun.StartedAt == nil || dbCheckRun.StartedAt.IsZero() {
				mockTaskCtx.GetLogger().(*MockLogger).Debug("collector: checkRun.StartedAt is nil or zero: " + dbCheckRun.Id)
				continue
			}
			messages = append(messages, errors.Must1(json.Marshal(dbCheckRun)))
		}
		return
	}

	// Execute the response parser
	messages, err := responseParser(checkRuns, 123)

	// Verify results
	assert.Nil(t, err)
	assert.Len(t, messages, 1, "Should only process jobs with valid StartedAt")

	// Verify the processed message is the correct job
	var processedJob DbCheckRun
	unmarshalErr := json.Unmarshal(messages[0], &processedJob)
	assert.Nil(t, unmarshalErr)
	assert.Equal(t, "normal-job-1", processedJob.Id)
	assert.Equal(t, "Normal Job", processedJob.Name)
	assert.NotNil(t, processedJob.StartedAt)

	assert.Len(t, mockLogger.debugCalls, 1)
	assert.Contains(t, mockLogger.debugCalls[0], "skipped-job-1")
}

func TestJobCollector_ResponseParser_SkipsJobsWithZeroStartedAt(t *testing.T) {
	mockLogger := &MockLogger{}
	mockTaskCtx := &MockTaskContext{logger: mockLogger}

	// Create test data with a job that has zero StartedAt
	now := time.Now()
	zeroTime := time.Time{} // Zero time

	// Check run with zero StartedAt should be skipped
	checkRunZero := GraphqlQueryCheckRun{
		Id:          "zero-time-job-1",
		Name:        "Zero Time Job",
		DatabaseId:  456,
		Status:      "completed",
		Conclusion:  "skipped",
		StartedAt:   &zeroTime,
		CompletedAt: &now,
	}

	// Check run with valid StartedAt should be processed
	checkRunNormal := GraphqlQueryCheckRun{
		Id:          "normal-job-1",
		Name:        "Normal Job",
		DatabaseId:  789,
		Status:      "completed",
		Conclusion:  "success",
		StartedAt:   &now,
		CompletedAt: &now,
	}

	checkRuns := []GraphqlQueryCheckRun{checkRunZero, checkRunNormal}

	responseParser := func(checkRuns []GraphqlQueryCheckRun, runId int) (messages []json.RawMessage, err errors.Error) {
		for _, checkRun := range checkRuns {
			dbCheckRun := &DbCheckRun{
				RunId:                runId,
				GraphqlQueryCheckRun: &checkRun,
			}
			if dbCheckRun.StartedAt == nil || dbCheckRun.StartedAt.IsZero() {
				mockTaskCtx.GetLogger().(*MockLogger).Debug("collector: checkRun.StartedAt is nil or zero: " + dbCheckRun.Id)
				continue
			}
			messages = append(messages, errors.Must1(json.Marshal(dbCheckRun)))
		}
		return
	}

	// Execute the response parser
	messages, err := responseParser(checkRuns, 123)

	// Verify results
	assert.Nil(t, err)
	assert.Len(t, messages, 1, "Should only process jobs with valid StartedAt")

	// Verify the processed message is the correct job
	var processedJob DbCheckRun
	unmarshalErr := json.Unmarshal(messages[0], &processedJob)
	assert.Nil(t, unmarshalErr)
	assert.Equal(t, "normal-job-1", processedJob.Id)
	assert.Equal(t, "Normal Job", processedJob.Name)
	assert.NotNil(t, processedJob.StartedAt)
	assert.False(t, processedJob.StartedAt.IsZero())

	assert.Len(t, mockLogger.debugCalls, 1)
	assert.Contains(t, mockLogger.debugCalls[0], "zero-time-job-1")
}

func TestJobCollector_ResponseParser_ProcessesValidJobs(t *testing.T) {
	// Create a mock task context
	mockLogger := &MockLogger{}
	mockTaskCtx := &MockTaskContext{logger: mockLogger}

	// Create test data with valid jobs
	now := time.Now()
	earlier := now.Add(-time.Hour)

	// Check run with zero StartedAt should be skipped
	checkRun1 := GraphqlQueryCheckRun{
		Id:          "job-1",
		Name:        "Job 1",
		DatabaseId:  456,
		Status:      "completed",
		Conclusion:  "success",
		StartedAt:   &earlier,
		CompletedAt: &now,
	}

	// Check run with valid StartedAt should be processed
	checkRun2 := GraphqlQueryCheckRun{
		Id:          "job-2",
		Name:        "Job 2",
		DatabaseId:  789,
		Status:      "in_progress",
		Conclusion:  "",
		StartedAt:   &now,
		CompletedAt: nil,
	}

	checkRuns := []GraphqlQueryCheckRun{checkRun1, checkRun2}

	responseParser := func(checkRuns []GraphqlQueryCheckRun, runId int) (messages []json.RawMessage, err errors.Error) {
		for _, checkRun := range checkRuns {
			dbCheckRun := &DbCheckRun{
				RunId:                runId,
				GraphqlQueryCheckRun: &checkRun,
			}
			if dbCheckRun.StartedAt == nil || dbCheckRun.StartedAt.IsZero() {
				mockTaskCtx.GetLogger().(*MockLogger).Debug("collector: checkRun.StartedAt is nil or zero: " + dbCheckRun.Id)
				continue
			}
			messages = append(messages, errors.Must1(json.Marshal(dbCheckRun)))
		}
		return
	}

	// Execute the response parser
	messages, err := responseParser(checkRuns, 123)

	assert.Nil(t, err)
	assert.Len(t, messages, 2, "Should process both valid jobs")

	// Verify both jobs were processed
	var job1, job2 DbCheckRun
	unmarshalErr1 := json.Unmarshal(messages[0], &job1)
	unmarshalErr2 := json.Unmarshal(messages[1], &job2)
	assert.Nil(t, unmarshalErr1)
	assert.Nil(t, unmarshalErr2)

	assert.Equal(t, "job-1", job1.Id)
	assert.Equal(t, "Job 1", job1.Name)
	assert.Equal(t, 123, job1.RunId)
	assert.NotNil(t, job1.StartedAt)

	assert.Equal(t, "job-2", job2.Id)
	assert.Equal(t, "Job 2", job2.Name)
	assert.Equal(t, 123, job2.RunId)
	assert.NotNil(t, job2.StartedAt)

	assert.Len(t, mockLogger.debugCalls, 0)
}

func TestDbCheckRun_StartedAtValidation(t *testing.T) {
	now := time.Now()
	zeroTime := time.Time{}

	testCases := []struct {
		name       string
		checkRun   *DbCheckRun
		expectNil  bool
		expectZero bool
	}{
		{
			name: "nil StartedAt",
			checkRun: &DbCheckRun{
				RunId: 123,
				GraphqlQueryCheckRun: &GraphqlQueryCheckRun{
					Id:        "test-1",
					StartedAt: nil,
				},
			},
			expectNil:  true,
			expectZero: false,
		},
		{
			name: "zero StartedAt",
			checkRun: &DbCheckRun{
				RunId: 123,
				GraphqlQueryCheckRun: &GraphqlQueryCheckRun{
					Id:        "test-2",
					StartedAt: &zeroTime,
				},
			},
			expectNil:  false,
			expectZero: true,
		},
		{
			name: "valid StartedAt",
			checkRun: &DbCheckRun{
				RunId: 123,
				GraphqlQueryCheckRun: &GraphqlQueryCheckRun{
					Id:        "test-3",
					StartedAt: &now,
				},
			},
			expectNil:  false,
			expectZero: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectNil {
				assert.True(t, tc.checkRun.StartedAt == nil, "StartedAt should be nil")
			} else {
				assert.False(t, tc.checkRun.StartedAt == nil, "StartedAt should not be nil")
			}

			if tc.expectZero && tc.checkRun.StartedAt != nil {
				assert.True(t, tc.checkRun.StartedAt.IsZero(), "StartedAt should be zero")
			} else if !tc.expectNil {
				assert.False(t, tc.checkRun.StartedAt.IsZero(), "StartedAt should not be zero")
			}
		})
	}
}
