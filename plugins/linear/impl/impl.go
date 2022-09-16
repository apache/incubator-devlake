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

package impl

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/linear/api"
	"github.com/apache/incubator-devlake/plugins/linear/models"
	"github.com/apache/incubator-devlake/plugins/linear/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/linear/tasks"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// make sure interface is implemented
var _ core.PluginMeta = (*Linear)(nil)
var _ core.PluginInit = (*Linear)(nil)
var _ core.PluginTask = (*Linear)(nil)
var _ core.PluginApi = (*Linear)(nil)
var _ core.PluginBlueprintV100 = (*Linear)(nil)
var _ core.CloseablePluginTask = (*Linear)(nil)

type Linear struct{}

func (plugin Linear) Description() string {
	return "collect some Linear data"
}

func (plugin Linear) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) errors.Error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Linear) SubTaskMetas() []core.SubTaskMeta {
	// TODO add your sub task here
	return []core.SubTaskMeta{
		tasks.CollectIssuesMeta,
	}
}

func (plugin Linear) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.LinearConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Linear connection by the given connection ID")
	}

	asyncGraphqlClient, err := tasks.NewLinearGraphqlClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	return &tasks.LinearTaskData{
		Options:       op,
		GraphqlClient: asyncGraphqlClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin Linear) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/linear"
}

func (plugin Linear) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Linear) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (plugin Linear) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Linear) Close(taskCtx core.TaskContext) errors.Error {
	return nil
}
