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

package services

import (
	"testing"
	"time"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models"
	"github.com/apache/devlake/core/plugin"
	contextimpl "github.com/apache/devlake/impls/context"
	"github.com/apache/devlake/impls/logruslog"
	mockdal "github.com/apache/devlake/mocks/core/dal"
	mockplugin "github.com/apache/devlake/mocks/core/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type blueprintFailureNotification func(PipelineNotificationParam) errors.Error

func (f blueprintFailureNotification) PipelineStatusChanged(p PipelineNotificationParam) errors.Error {
	return f(p)
}

// Service globals follow the existing service tests; these tests must not run
// in parallel, and normal-pipeline notification goroutines are joined below.
func blueprintFailureTestDB(t *testing.T) *mockdal.Dal {
	t.Helper()
	testDB := mockdal.NewDal(t)
	withProjectTestDatabase(t, testDB)
	oldRes, oldLogger := basicRes, logger
	oldDefault, oldCustom := defaultNotificationService, customPipelineNotificationService
	basicRes = contextimpl.NewDefaultBasicRes(nil, logruslog.Global, testDB)
	logger = logruslog.Global
	defaultNotificationService, customPipelineNotificationService = nil, nil
	t.Cleanup(func() {
		basicRes, logger = oldRes, oldLogger
		defaultNotificationService, customPipelineNotificationService = oldDefault, oldCustom
	})
	return testDB
}

func failingPlanBlueprint(t *testing.T, planErr errors.Error) *models.Blueprint {
	t.Helper()
	name := "test-blueprint-plan-failure"
	source := new(mockplugin.CompositeDataSourcePluginBlueprintV200)
	source.On("MakeDataSourcePipelinePlanV200", uint64(7), mock.Anything).
		Return(models.PipelinePlan(nil), []plugin.Scope(nil), planErr).Once()
	require.NoError(t, plugin.RegisterPlugin(name, source))
	t.Cleanup(func() {
		delete(plugin.AllPlugins(), name)
		source.AssertExpectations(t)
	})
	bp := &models.Blueprint{
		Name: "scheduled-source", Mode: models.BLUEPRINT_MODE_NORMAL,
		Priority: 3, Labels: []string{"engineering"},
		Connections: []*models.BlueprintConnection{{
			PluginName: name, ConnectionId: 7,
			Scopes: []*models.BlueprintScope{{ScopeId: "repo"}},
		}},
	}
	bp.ID = 12
	bp.SkipOnFail = true
	return bp
}

func expectPipelineNotificationRead(t *testing.T, testDB *mockdal.Dal, pipeline **models.Pipeline) {
	t.Helper()
	testDB.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		switch target := args.Get(0).(type) {
		case *models.Pipeline:
			*target = **pipeline
		case *models.Blueprint:
			target.ProjectName = "engineering"
		default:
			t.Errorf("unexpected notification lookup: %T", target)
		}
	}).Return(nil).Twice()
	testDB.On("Pluck", "name", mock.Anything, mock.Anything).Return(nil).Twice()
}

func TestBlueprintJobRecordsPlanFailure(t *testing.T) {
	for _, notificationFails := range []bool{false, true} {
		name := "notifies after commit"
		if notificationFails {
			name = "notification failure keeps history"
		}
		t.Run(name, func(t *testing.T) {
			testDB := blueprintFailureTestDB(t)
			planErr := errors.Unauthorized.New("installation credentials rejected")
			bp := failingPlanBlueprint(t, planErr)
			started := time.Now()
			tx := mockdal.NewTransaction(t)
			testDB.On("Begin").Return(tx).Once()
			var stored *models.Pipeline
			tx.On("Create", mock.Anything).Run(func(args mock.Arguments) {
				switch row := args.Get(0).(type) {
				case *models.Pipeline:
					row.ID = 123
					stored = row
					assert.Equal(t, models.TASK_FAILED, row.Status)
					assert.Equal(t, bp.ID, row.BlueprintId)
					assert.Equal(t, bp.Name, row.Name)
					assert.Equal(t, bp.Priority, row.Priority)
					assert.Equal(t, bp.SyncPolicy, row.SyncPolicy)
					assert.Equal(t, models.PipelinePlan{}, row.Plan)
					assert.Zero(t, row.TotalTasks)
					assert.Zero(t, row.FinishedTasks)
					require.NotNil(t, row.BeganAt)
					require.NotNil(t, row.FinishedAt)
					assert.False(t, row.BeganAt.Before(started))
					assert.False(t, row.FinishedAt.Before(*row.BeganAt))
					assert.Equal(t, int(row.FinishedAt.Sub(*row.BeganAt).Seconds()), row.SpentSeconds)
					assert.Contains(t, row.Message, planErr.Error())
					assert.Contains(t, row.ErrorName, "Failed to generate pipeline plan")
				case *[]models.DbPipelineLabel:
					require.NotNil(t, stored)
					assert.Equal(t, []models.DbPipelineLabel{{PipelineId: 123, Name: "engineering"}}, *row)
				default:
					t.Errorf("planning failure must not create tasks: %T", row)
				}
			}).Return(nil).Twice()
			committed := false
			tx.On("Commit").Run(func(mock.Arguments) { committed = true }).Return(nil).Once()
			expectPipelineNotificationRead(t, testDB, &stored)
			notifications := 0
			customPipelineNotificationService = blueprintFailureNotification(func(p PipelineNotificationParam) errors.Error {
				assert.True(t, committed)
				assert.Equal(t, uint64(123), p.PipelineID)
				assert.Equal(t, "engineering", p.ProjectName)
				assert.Equal(t, models.TASK_FAILED, p.Status)
				assert.Equal(t, stored.FinishedAt, p.FinishedAt)
				notifications++
				if notificationFails {
					return errors.Default.New("notification unavailable")
				}
				return nil
			})
			BlueprintJob{Blueprint: bp}.Run()
			assert.Equal(t, 1, notifications)
			tx.AssertNotCalled(t, "Rollback")
		})
	}
}

func TestBlueprintPlanFailureRollsBackAndPreservesCause(t *testing.T) {
	for _, stage := range []string{"pipeline", "labels", "commit", "rollback"} {
		t.Run(stage, func(t *testing.T) {
			testDB := blueprintFailureTestDB(t)
			planErr := errors.Default.New("source planning failed")
			bp := failingPlanBlueprint(t, planErr)
			tx := mockdal.NewTransaction(t)
			testDB.On("Begin").Return(tx).Once()
			writeErr := errors.Default.New("database write failed")
			if stage == "pipeline" || stage == "rollback" {
				tx.On("Create", mock.Anything).Return(writeErr).Once()
			} else {
				tx.On("Create", mock.Anything).Return(nil).Once()
				if stage == "labels" {
					tx.On("Create", mock.Anything).Return(writeErr).Once()
				} else {
					tx.On("Create", mock.Anything).Return(nil).Once()
					tx.On("Commit").Return(writeErr).Once()
				}
			}
			var rollbackErr errors.Error
			if stage == "rollback" {
				rollbackErr = errors.Default.New("rollback failed")
			}
			tx.On("Rollback").Return(rollbackErr).Once()
			customPipelineNotificationService = blueprintFailureNotification(func(PipelineNotificationParam) errors.Error {
				t.Error("uncommitted failure must not be notified")
				return nil
			})
			pipeline, err := createPipelineByBlueprint(bp, &bp.SyncPolicy, true)
			assert.Nil(t, pipeline)
			assert.Same(t, planErr, err)
		})
	}
}

func TestBlueprintManualPlanningFailureDoesNotCreateHistory(t *testing.T) {
	testDB := blueprintFailureTestDB(t)
	planErr := errors.Default.New("source planning failed")
	bp := failingPlanBlueprint(t, planErr)
	bp.Enable = true
	testDB.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(0).(*models.Blueprint) = *bp
	}).Return(nil).Once()
	testDB.On("Pluck", "name", mock.Anything, mock.Anything).Return(nil).Once()
	testDB.On("All", mock.Anything, mock.Anything).Return(nil).Twice()
	pipeline, err := TriggerBlueprint(bp.ID, &models.TriggerSyncPolicy{}, false)
	assert.Nil(t, pipeline)
	assert.Same(t, planErr, err)
	testDB.AssertNotCalled(t, "Begin")
}

func TestBlueprintPlanFailureWithoutLabels(t *testing.T) {
	testDB := blueprintFailureTestDB(t)
	planErr := errors.Default.New("source planning failed")
	bp := failingPlanBlueprint(t, planErr)
	bp.Labels = nil
	tx := mockdal.NewTransaction(t)
	testDB.On("Begin").Return(tx).Once()
	tx.On("Create", mock.MatchedBy(func(row *models.Pipeline) bool {
		return row.Status == models.TASK_FAILED && len(row.Labels) == 0
	})).Return(nil).Once()
	tx.On("Commit").Return(nil).Once()
	pipeline, err := createPipelineByBlueprint(bp, &bp.SyncPolicy, true)
	assert.Nil(t, pipeline)
	assert.Same(t, planErr, err)
	tx.AssertNotCalled(t, "Rollback")
}

func TestBlueprintJobEmptyPlanSentinelDoesNotCreateHistory(t *testing.T) {
	testDB := blueprintFailureTestDB(t)
	bp := failingPlanBlueprint(t, ErrEmptyPlan)
	BlueprintJob{Blueprint: bp}.Run()
	testDB.AssertNotCalled(t, "Begin")
}

func TestBlueprintJobExistingCreationPaths(t *testing.T) {
	for _, scenario := range []string{"success", "empty", "pending", "insert failure"} {
		t.Run(scenario, func(t *testing.T) {
			testDB := blueprintFailureTestDB(t)
			bp := &models.Blueprint{Name: "advanced", Mode: models.BLUEPRINT_MODE_ADVANCED}
			bp.ID = 20
			if scenario != "empty" {
				bp.Plan = models.PipelinePlan{{{Plugin: "test"}}}
			}
			tx := mockdal.NewTransaction(t)
			testDB.On("Begin").Return(tx).Once()
			tx.On("LockTables", mock.Anything).Return(nil).Once()
			tx.On("UnlockTables").Return(nil).Once()
			var pending int64
			if scenario == "pending" {
				pending = 1
			}
			tx.On("Count", mock.Anything).Return(pending, nil).Once()
			var stored *models.Pipeline
			if scenario != "pending" {
				var insertErr errors.Error
				if scenario == "insert failure" {
					insertErr = errors.Default.New("insert failed")
				}
				tx.On("Create", mock.Anything).Run(func(args mock.Arguments) {
					stored = args.Get(0).(*models.Pipeline)
					stored.ID = 456
					assert.Equal(t, models.TASK_CREATED, stored.Status)
				}).Return(insertErr).Once()
			}
			notified := make(chan struct{})
			if scenario == "success" {
				tx.On("Create", mock.Anything).Run(func(args mock.Arguments) {
					task := args.Get(0).(*models.Task)
					assert.Equal(t, "test", task.Plugin)
					assert.Equal(t, uint64(456), task.PipelineId)
				}).Return(nil).Once()
				tx.On("Update", mock.Anything).Return(nil).Once()
				tx.On("Commit").Return(nil).Once()
				expectPipelineNotificationRead(t, testDB, &stored)
				customPipelineNotificationService = blueprintFailureNotification(func(p PipelineNotificationParam) errors.Error {
					assert.Equal(t, models.TASK_CREATED, p.Status)
					close(notified)
					return nil
				})
			} else {
				tx.On("Rollback").Return(nil).Once()
			}
			BlueprintJob{Blueprint: bp}.Run()
			if scenario == "success" {
				select {
				case <-notified:
				case <-time.After(5 * time.Second):
					t.Fatal("created notification did not complete")
				}
				assert.Equal(t, 1, stored.TotalTasks)
			}
		})
	}
}
