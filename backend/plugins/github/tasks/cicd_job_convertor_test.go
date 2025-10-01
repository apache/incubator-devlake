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
	"fmt"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/stretchr/testify/assert"
)

func GenJobIDWithReflect(jobIdGen *didgen.DomainIdGenerator) {
	connectionId := uint64(1)
	runId := 1
	lineId := 1
	jobIdGen.Generate(connectionId, runId, lineId)
}

func GenJobID() {
	connectionId := uint64(1)
	runId := 1
	lineId := 1
	fmt.Sprintf("GithubJob:%d:%d:%d", connectionId, runId, lineId)
}

func BenchmarkGenJobIDWithReflect(b *testing.B) {
	mockMeta := mockplugin.NewPluginMeta(b)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")
	mockMeta.On("Name").Return("github").Maybe()
	err := plugin.RegisterPlugin("github", mockMeta)
	assert.NoError(b, err)

	jobIdGen := didgen.NewDomainIdGenerator(&models.GithubJob{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GenJobIDWithReflect(jobIdGen)
	}
	b.StopTimer()
	//BenchmarkGenJobIDWithReflect-8   	 5611773	       208.9 ns/op
}

func BenchmarkGenJobID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenJobID()
	}
	//BenchmarkGenJobID-8   	11078593	        99.43 ns/op
}

func TestConvertJobs_SkipNoStartedAt(t *testing.T) {
	job := &models.GithubJob{
		ID:        123,
		RunID:     456,
		Name:      "test-job",
		StartedAt: nil,
	}

	convert := func(inputRow interface{}) ([]interface{}, error) {
		line := inputRow.(*models.GithubJob)
		if line.StartedAt == nil {
			return nil, nil
		}
		createdAt := *line.StartedAt
		domainJob := &devops.CICDTask{
			Name: line.Name,
			TaskDatesInfo: devops.TaskDatesInfo{
				CreatedDate:  createdAt,
				StartedDate:  line.StartedAt,
				FinishedDate: line.CompletedAt,
			},
		}
		return []interface{}{domainJob}, nil
	}

	result, err := convert(job)
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestConvertJobs_WithStartedAt(t *testing.T) {
	now := time.Now()
	job := &models.GithubJob{
		ID:        123,
		RunID:     456,
		Name:      "test-job",
		StartedAt: &now,
	}

	convert := func(inputRow interface{}) ([]interface{}, error) {
		line := inputRow.(*models.GithubJob)
		if line.StartedAt == nil {
			return nil, nil
		}
		createdAt := *line.StartedAt
		domainJob := &devops.CICDTask{
			Name: line.Name,
			TaskDatesInfo: devops.TaskDatesInfo{
				CreatedDate:  createdAt,
				StartedDate:  line.StartedAt,
				FinishedDate: line.CompletedAt,
			},
		}
		return []interface{}{domainJob}, nil
	}

	result, err := convert(job)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-job", result[0].(*devops.CICDTask).Name)
}

// Test fixtures and helpers
type jobFixture struct {
	name             string
	conclusion       string
	expectedResult   string
	expectedOriginal string
}

func createTestJob(conclusion string) *models.GithubJob {
	now := time.Now()
	return &models.GithubJob{
		ID:         123,
		RunID:      456,
		RepoId:     789,
		Name:       "test-job",
		StartedAt:  &now,
		Conclusion: conclusion,
		Status:     StatusCompleted,
	}
}

func createTestJobWithCompletedAt(conclusion string) *models.GithubJob {
	now := time.Now()
	completed := now.Add(5 * time.Minute)
	return &models.GithubJob{
		ID:          123,
		RunID:       456,
		RepoId:      789,
		Name:        "test-job",
		StartedAt:   &now,
		CompletedAt: &completed,
		Conclusion:  conclusion,
		Status:      StatusCompleted,
	}
}

func getResultRule() *devops.ResultRule {
	return &devops.ResultRule{
		Success:  []string{StatusSuccess},
		Failure:  []string{StatusFailure, StatusTimedOut, StatusStartUpFailure},
		Canceled: []string{StatusCancelled},
		Skipped:  []string{StatusSkipped},
		Default:  devops.RESULT_DEFAULT,
	}
}

// convertJobForTest simulates the actual conversion logic from ConvertJobs
func convertJobForTest(inputRow interface{}) ([]interface{}, error) {
	line := inputRow.(*models.GithubJob)
	if line.StartedAt == nil {
		return nil, nil
	}
	createdAt := *line.StartedAt
	domainJob := &devops.CICDTask{
		Name: line.Name,
		TaskDatesInfo: devops.TaskDatesInfo{
			CreatedDate:  createdAt,
			StartedDate:  line.StartedAt,
			FinishedDate: line.CompletedAt,
		},
		Result:         devops.GetResult(getResultRule(), line.Conclusion),
		OriginalResult: line.Conclusion,
	}

	// Calculate duration if both dates are available
	if line.CompletedAt != nil && line.StartedAt != nil {
		domainJob.DurationSec = float64(line.CompletedAt.Sub(*line.StartedAt).Milliseconds() / 1e3)
	}

	return []interface{}{domainJob}, nil
}

func TestConvertJobs_ResultMapping(t *testing.T) {
	testCases := []jobFixture{
		{
			name:             "Success mapping",
			conclusion:       StatusSuccess,
			expectedResult:   devops.RESULT_SUCCESS,
			expectedOriginal: StatusSuccess,
		},
		{
			name:             "Failure mapping",
			conclusion:       StatusFailure,
			expectedResult:   devops.RESULT_FAILURE,
			expectedOriginal: StatusFailure,
		},
		{
			name:             "TimedOut mapping to Failure",
			conclusion:       StatusTimedOut,
			expectedResult:   devops.RESULT_FAILURE,
			expectedOriginal: StatusTimedOut,
		},
		{
			name:             "StartupFailure mapping to Failure",
			conclusion:       StatusStartUpFailure,
			expectedResult:   devops.RESULT_FAILURE,
			expectedOriginal: StatusStartUpFailure,
		},
		{
			name:             "Canceled mapping",
			conclusion:       StatusCancelled,
			expectedResult:   devops.RESULT_CANCELED,
			expectedOriginal: StatusCancelled,
		},
		{
			name:             "Skipped mapping",
			conclusion:       StatusSkipped,
			expectedResult:   devops.RESULT_SKIPPED,
			expectedOriginal: StatusSkipped,
		},
		{
			name:             "Unknown status defaults",
			conclusion:       "UNKNOWN_STATUS",
			expectedResult:   devops.RESULT_DEFAULT,
			expectedOriginal: "UNKNOWN_STATUS",
		},
		{
			name:             "Empty conclusion defaults",
			conclusion:       "",
			expectedResult:   devops.RESULT_DEFAULT,
			expectedOriginal: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := createTestJob(tc.conclusion)

			result, err := convertJobForTest(job)

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, 1)

			domainJob := result[0].(*devops.CICDTask)
			assert.Equal(t, tc.expectedResult, domainJob.Result)
			assert.Equal(t, tc.expectedOriginal, domainJob.OriginalResult)
			assert.Equal(t, "test-job", domainJob.Name)
		})
	}
}

func TestConvertJobs_DurationCalculation(t *testing.T) {
	job := createTestJobWithCompletedAt(StatusSuccess)

	result, err := convertJobForTest(job)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	domainJob := result[0].(*devops.CICDTask)
	assert.Equal(t, float64(300), domainJob.DurationSec) // 5 minutes = 300 seconds
	assert.NotNil(t, domainJob.StartedDate)
	assert.NotNil(t, domainJob.FinishedDate)
}

func TestConvertJobs_NoStartedAtFixture(t *testing.T) {
	job := &models.GithubJob{
		ID:         123,
		RunID:      456,
		Name:       "test-job",
		StartedAt:  nil, // This should cause the job to be skipped
		Conclusion: StatusSuccess,
	}

	result, err := convertJobForTest(job)

	assert.Nil(t, err)
	assert.Nil(t, result, "Jobs with no StartedAt should be skipped")
}
