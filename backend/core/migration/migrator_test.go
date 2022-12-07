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

package migration

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/impls/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHasPendingScripts(t *testing.T) {
	// simulate db reaction
	mockDal := new(mockdal.Dal)
	mockDal.On("AutoMigrate", mock.Anything, mock.Anything).Return(nil).Once()
	mockDal.On("All", mock.Anything, mock.Anything).Return(func(i interface{}, _ ...dal.Clause) errors.Error {
		precords := i.(*[]MigrationHistory)
		*precords = []MigrationHistory{
			{ScriptName: "A", ScriptVersion: 1, Comment: "UniTest", CreatedAt: time.Now()},
			{ScriptName: "B", ScriptVersion: 2, Comment: "UniTest", CreatedAt: time.Now()},
			{ScriptName: "C", ScriptVersion: 3, Comment: "UniTest", CreatedAt: time.Now()},
		}
		return nil
	}).Once()
	mockDal.On("Create", &MigrationHistory{
		ScriptName:    "E",
		ScriptVersion: 4,
		Comment:       "UnitTest",
	}, mock.Anything).Return(nil).Once()
	mockDal.On("Create", &MigrationHistory{
		ScriptName:    "D",
		ScriptVersion: 5,
		Comment:       "UnitTest",
	}, mock.Anything).Return(nil).Once()

	// migrator initialization
	basicRes := context.NewDefaultBasicRes(viper.New(), unithelper.DummyLogger(), mockDal)
	migrator, err := NewMigrator(basicRes)
	assert.Nil(t, err)

	// assuming we have 2 new scripts
	scriptD := new(mockplugin.MigrationScript)
	scriptD.On("Up", mock.Anything).Return(nil).Once()
	scriptD.On("Version").Return(uint64(5))
	scriptD.On("Name").Return("D")
	scriptE := new(mockplugin.MigrationScript)
	scriptE.On("Up", mock.Anything).Return(nil).Once()
	scriptE.On("Version").Return(uint64(4))
	scriptE.On("Name").Return("E")
	migrator.Register([]plugin.MigrationScript{scriptD, scriptE}, "UnitTest")

	// we should have pending scripts
	assert.True(t, migrator.HasPendingScripts())

	// lets try migrating
	assert.Nil(t, migrator.Execute())

	// should not be any pending scripts anymore
	assert.False(t, migrator.HasPendingScripts())

	// make sure all method got called
	mockDal.AssertExpectations(t)
}
