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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func at(hour int) time.Time {
	return time.Date(2026, 8, 9, hour, 0, 0, 0, time.UTC)
}

func TestFirstDeploymentAfter(t *testing.T) {
	// Sorted ascending, as loadSubProjectDeployments guarantees.
	deployments := []deployedAt{
		{Id: "deploy-08", FinishedDate: at(8)},
		{Id: "deploy-12", FinishedDate: at(12)},
		{Id: "deploy-18", FinishedDate: at(18)},
	}

	cases := []struct {
		name        string
		deployments []deployedAt
		mergedDate  *time.Time
		expectedId  string
	}{
		{
			name:        "picks the earliest deployment after the merge",
			deployments: deployments,
			mergedDate:  ptrTime(at(10)),
			expectedId:  "deploy-12",
		},
		{
			name:        "merge before every deployment picks the first",
			deployments: deployments,
			mergedDate:  ptrTime(at(1)),
			expectedId:  "deploy-08",
		},
		{
			name:        "merge after every deployment has none to link",
			deployments: deployments,
			mergedDate:  ptrTime(at(20)),
			expectedId:  "",
		},
		{
			name:        "a deployment finishing exactly at merge time does not count",
			deployments: deployments,
			mergedDate:  ptrTime(at(12)),
			expectedId:  "deploy-18",
		},
		{
			name:        "no deployments at all",
			deployments: nil,
			mergedDate:  ptrTime(at(10)),
			expectedId:  "",
		},
		{
			name:        "unmerged pull request",
			deployments: deployments,
			mergedDate:  nil,
			expectedId:  "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := firstDeploymentAfter(tc.deployments, tc.mergedDate)
			if tc.expectedId == "" {
				assert.Nil(t, got)
				return
			}
			assert.NotNil(t, got)
			assert.Equal(t, tc.expectedId, got.Id)
		})
	}
}

func TestComputeTimeSpan(t *testing.T) {
	start := at(10)
	end := at(12)

	t.Run("whole minutes between two times", func(t *testing.T) {
		got := computeTimeSpan(&start, &end)
		assert.NotNil(t, got)
		assert.Equal(t, int64(120), *got)
	})

	t.Run("partial minutes round up", func(t *testing.T) {
		later := start.Add(90 * time.Second)
		got := computeTimeSpan(&start, &later)
		assert.NotNil(t, got)
		assert.Equal(t, int64(2), *got)
	})

	t.Run("negative spans are discarded", func(t *testing.T) {
		assert.Nil(t, computeTimeSpan(&end, &start))
	})

	t.Run("missing endpoints yield nil", func(t *testing.T) {
		assert.Nil(t, computeTimeSpan(nil, &end))
		assert.Nil(t, computeTimeSpan(&start, nil))
		assert.Nil(t, computeTimeSpan(nil, nil))
	})
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
