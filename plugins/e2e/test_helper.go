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

package e2e

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/mock"
)

// nolint:unused
type mockPluginHelper struct {
	mock *mocks.TestPlugin
}

// nolint:unused
func newMockPluginHelper() *mockPluginHelper {
	plugin := new(mocks.TestPlugin)
	plugin.On("Description").Return("desc").Maybe()
	plugin.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/e2e").Maybe()
	return &mockPluginHelper{mock: plugin}
}

// nolint:unused
func (m *mockPluginHelper) SubTaskMetas(f func() []core.SubTaskMeta) *mock.Call {
	return m.mock.On("SubTaskMetas").Return(f)
}

// nolint:unused
func (m *mockPluginHelper) PrepareTaskData(f func(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error)) *mock.Call {
	// workaround solution because mockery doesn't support a good 'passthrough' behavior for multi-return functions
	var call *mock.Call
	call = m.mock.On("PrepareTaskData", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		ifc, err := f(args.Get(0).(core.TaskContext), args.Get(1).(map[string]interface{}))
		call.ReturnArguments = []interface{}{ifc, err}
	})
	return call
}

// nolint:unused
func (m *mockPluginHelper) GetPlugin() *mocks.TestPlugin {
	return m.mock
}
