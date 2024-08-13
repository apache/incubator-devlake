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

package api

import (
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCollectorStateManager(t *testing.T) {
	time0 := errors.Must1(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z"))
	time1 := errors.Must1(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))
	time2 := errors.Must1(time.Parse(time.RFC3339, "2022-01-01T00:00:00Z"))
	for _, tc := range []struct {
		name                      string
		state                     *models.CollectorLatestState
		syncPolicy                *models.SyncPolicy
		expectedIsIncremental     bool
		expectedSince             *time.Time
		expectedNewStateTimeAfter *time.Time
	}{
		{
			name:                      "syncPolicy has no timeAfter - First run",
			state:                     &models.CollectorLatestState{LatestSuccessStart: nil},
			syncPolicy:                &models.SyncPolicy{TimeAfter: nil},
			expectedIsIncremental:     false,
			expectedSince:             nil,
			expectedNewStateTimeAfter: nil,
		},
		{
			name:                      "syncPolicy has no timeAfter - Second run",
			state:                     &models.CollectorLatestState{LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TimeAfter: nil},
			expectedIsIncremental:     true,
			expectedSince:             &time1,
			expectedNewStateTimeAfter: nil,
		},
		{
			name:                      "syncPolicy has no timeAfter - Third run with timeAfter specified",
			state:                     &models.CollectorLatestState{LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time1},
			expectedIsIncremental:     true,
			expectedSince:             &time1,
			expectedNewStateTimeAfter: nil,
		},
		{
			name:                      "syncPolicy has timeAfter - First run",
			state:                     &models.CollectorLatestState{LatestSuccessStart: nil},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time1},
			expectedIsIncremental:     false,
			expectedSince:             &time1,
			expectedNewStateTimeAfter: &time1,
		},
		{
			name:                      "syncPolicy has timeAfter - Second run with a later timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time2},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time2},
			expectedIsIncremental:     true,
			expectedSince:             &time2,
			expectedNewStateTimeAfter: &time1,
		},
		{
			name:                      "syncPolicy has timeAfter - Third run with a earlier timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time0},
			expectedIsIncremental:     false,
			expectedSince:             &time0,
			expectedNewStateTimeAfter: &time0,
		},
		{
			name:                      "syncPolicy has timeAfter - Fourth run with a same timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time2},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time1},
			expectedIsIncremental:     true,
			expectedSince:             &time2,
			expectedNewStateTimeAfter: &time1,
		},
		{
			name:                      "Full sync - with timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TriggerSyncPolicy: models.TriggerSyncPolicy{FullSync: true}},
			expectedIsIncremental:     false,
			expectedSince:             &time1,
			expectedNewStateTimeAfter: &time1,
		},
		{
			name:                      "Full sync - with newer timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time2, TriggerSyncPolicy: models.TriggerSyncPolicy{FullSync: true}},
			expectedIsIncremental:     false,
			expectedSince:             &time2,
			expectedNewStateTimeAfter: &time2,
		},
		{
			name:                      "Full sync - with older timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: &time1, LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TimeAfter: &time0, TriggerSyncPolicy: models.TriggerSyncPolicy{FullSync: true}},
			expectedIsIncremental:     false,
			expectedSince:             &time0,
			expectedNewStateTimeAfter: &time0,
		},
		{
			name:                      "Full sync - without timeAfter",
			state:                     &models.CollectorLatestState{TimeAfter: nil, LatestSuccessStart: &time1},
			syncPolicy:                &models.SyncPolicy{TriggerSyncPolicy: models.TriggerSyncPolicy{FullSync: true}},
			expectedIsIncremental:     false,
			expectedSince:             nil,
			expectedNewStateTimeAfter: nil,
		},
	} {
		started := time.Now()
		t.Run(tc.name, func(t *testing.T) {
			mockBasicRes := newMockBasicRes(tc.state)
			stateManager, err := NewCollectorStateManager(mockBasicRes, tc.syncPolicy, "table", "params")
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedSince, stateManager.since)
			assert.Equal(t, tc.expectedIsIncremental, stateManager.isIncremental)
			assert.Nil(t, stateManager.Close())
			assert.Equal(t, tc.expectedNewStateTimeAfter, stateManager.state.TimeAfter)
			// LatestSuccessStart should be updated
			assert.GreaterOrEqual(t, stateManager.state.LatestSuccessStart.Unix(), started.Unix())
			// First and update should both be called once
			mockBasicRes.AssertExpectations(t)
		})
	}
}

func newMockBasicRes(state *models.CollectorLatestState) *mockcontext.BasicRes {
	// Refresh Global Variables and set the sql mock
	return unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.CollectorLatestState)
			*dst = *state
		}).Return(nil).Once()
		mockDal.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	})
}
