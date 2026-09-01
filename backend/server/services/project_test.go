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
	"fmt"
	"testing"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	blueprintservices "github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	dalmocks "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testOrdinaryPlugin struct {
	name string
}

func (p *testOrdinaryPlugin) Description() string { return "test ordinary plugin" }
func (p *testOrdinaryPlugin) RootPkgPath() string { return "plugins/test_ordinary" }
func (p *testOrdinaryPlugin) Name() string        { return p.name }

type testHookPlugin struct {
	name        string
	enabled     bool
	deleteCalls []struct {
		tx          dal.Transaction
		projectName string
	}
	deleteErr errors.Error
}

func (p *testHookPlugin) Description() string { return "test hook plugin" }
func (p *testHookPlugin) RootPkgPath() string { return "plugins/test_hook" }
func (p *testHookPlugin) Name() string        { return p.name }
func (p *testHookPlugin) BeforeDeleteProject(tx dal.Transaction, projectName string) errors.Error {
	if !p.enabled {
		return nil
	}
	p.deleteCalls = append(p.deleteCalls, struct {
		tx          dal.Transaction
		projectName string
	}{tx: tx, projectName: projectName})
	return p.deleteErr
}

func TestRunProjectDeleteHooks(t *testing.T) {
	t.Run("skips ordinary plugins without ProjectDeleteHook", func(t *testing.T) {
		ordinary := &testOrdinaryPlugin{name: "test-ordinary-skip"}
		assert.NoError(t, plugin.RegisterPlugin(ordinary.Name(), ordinary))

		tx := dalmocks.NewTransaction(t)
		assert.NoError(t, runProjectDeleteHooks(tx, "test-project"))
	})

	t.Run("invokes implementing plugins with exact transaction and project name", func(t *testing.T) {
		hook := &testHookPlugin{name: "test-hook-invoke", enabled: true}
		t.Cleanup(func() { hook.enabled = false })
		assert.NoError(t, plugin.RegisterPlugin(hook.Name(), hook))

		tx := dalmocks.NewTransaction(t)
		assert.NoError(t, runProjectDeleteHooks(tx, "test-project"))

		assert.Equal(t, 1, len(hook.deleteCalls))
		assert.Equal(t, tx, hook.deleteCalls[0].tx)
		assert.Equal(t, "test-project", hook.deleteCalls[0].projectName)
	})

	t.Run("returns wrapped error on hook veto", func(t *testing.T) {
		expectedErr := errors.Default.New("project delete vetoed by plugin")
		hook := &testHookPlugin{
			name:      "test-hook-veto",
			enabled:   true,
			deleteErr: expectedErr,
		}
		t.Cleanup(func() { hook.enabled = false })
		assert.NoError(t, plugin.RegisterPlugin(hook.Name(), hook))

		tx := dalmocks.NewTransaction(t)
		err := runProjectDeleteHooks(tx, "test-project")

		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Contains(t, err.Error(), fmt.Sprintf("error executing delete hook for plugin %s", hook.Name()))
	})
}

func TestDeleteProject_RollsBackOnDeleteHookVeto(t *testing.T) {
	hookErr := errors.Default.New("hook rejection")
	hook := &testHookPlugin{
		name:      "test-hook-rollback-veto",
		enabled:   true,
		deleteErr: hookErr,
	}
	t.Cleanup(func() { hook.enabled = false })
	assert.NoError(t, plugin.RegisterPlugin(hook.Name(), hook))

	mockDB := dalmocks.NewDal(t)
	tx := dalmocks.NewTransaction(t)
	notFound := errors.NotFound.New("blueprint not found")

	mockDB.On("First", mock.Anything, mock.Anything).Return(nil).Once()
	mockDB.On("First", mock.Anything, mock.Anything).Return(notFound).Once()
	mockDB.On("IsErrorNotFound", mock.Anything).Return(true).Twice()
	mockDB.On("Begin").Return(tx).Once()
	tx.On("Rollback").Return(nil).Once()

	previousDB, previousManager := db, bpManager
	db = mockDB
	bpManager = blueprintservices.NewBlueprintManager(mockDB)
	t.Cleanup(func() {
		db = previousDB
		bpManager = previousManager
	})

	err := DeleteProject("project-veto")

	assert.Error(t, err)
	assert.ErrorIs(t, err, hookErr)
	assert.Contains(t, err.Error(), fmt.Sprintf("error executing delete hook for plugin %s", hook.Name()))
	assert.Equal(t, 1, len(hook.deleteCalls))
	assert.Equal(t, "project-veto", hook.deleteCalls[0].projectName)
}

func TestDeleteProject_SuccessfulDeletionInSingleTransaction(t *testing.T) {
	hook := &testHookPlugin{
		name:    "test-hook-success",
		enabled: true,
	}
	t.Cleanup(func() { hook.enabled = false })
	assert.NoError(t, plugin.RegisterPlugin(hook.Name(), hook))

	mockDB := dalmocks.NewDal(t)
	tx := dalmocks.NewTransaction(t)

	// 1. Verify project exists
	mockDB.On("First", mock.MatchedBy(func(target interface{}) bool {
		_, ok := target.(*models.Project)
		return ok
	}), mock.Anything).Return(nil).Once()

	// 2. Blueprint lookup before transaction
	mockDB.On("First", mock.MatchedBy(func(target interface{}) bool {
		bp, ok := target.(*models.Blueprint)
		if ok {
			bp.ID = 42
		}
		return ok
	}), mock.Anything).Return(nil).Once()
	mockDB.On("Pluck", "name", mock.Anything, mock.Anything).Return(nil).Once()
	mockDB.On("All", mock.MatchedBy(func(target interface{}) bool {
		_, ok := target.(*[]*models.BlueprintConnection)
		return ok
	}), mock.Anything).Return(nil).Once()

	// 3. Pipeline lookup for unfinished pipelines check
	mockDB.On("Count", mock.Anything).Return(int64(0), nil).Once()
	mockDB.On("All", mock.MatchedBy(func(target interface{}) bool {
		_, ok := target.(*[]*models.Pipeline)
		return ok
	}), mock.Anything).Return(nil).Once()

	// 4. Begin transaction
	mockDB.On("Begin").Return(tx).Once()

	// 5. Blueprint lookup inside transaction for deletion
	tx.On("First", mock.MatchedBy(func(target interface{}) bool {
		bp, ok := target.(*models.Blueprint)
		if ok {
			bp.ID = 42
		}
		return ok
	}), mock.Anything).Return(nil).Once()

	// 6. DeleteBlueprintInTransaction calls tx.Delete 4 times
	// 7. Core project deletes call tx.Delete 5 times
	// Total tx.Delete calls = 9
	tx.On("Delete", mock.Anything, mock.Anything).Return(nil).Times(9)

	// 8. Single commit for everything
	tx.On("Commit").Return(nil).Once()

	previousDB, previousManager := db, bpManager
	db = mockDB
	bpManager = blueprintservices.NewBlueprintManager(mockDB)
	t.Cleanup(func() {
		db = previousDB
		bpManager = previousManager
	})

	err := DeleteProject("project-success")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(hook.deleteCalls))
	assert.Equal(t, tx, hook.deleteCalls[0].tx)
	assert.Equal(t, "project-success", hook.deleteCalls[0].projectName)
}
