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
	"fmt"
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/webhook/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"github.com/apache/incubator-devlake/plugins/webhook/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/webhook/tasks"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// make sure interface is implemented
var _ core.PluginMeta = (*Webhook)(nil)
var _ core.PluginInit = (*Webhook)(nil)
var _ core.PluginTask = (*Webhook)(nil)
var _ core.PluginApi = (*Webhook)(nil)
var _ core.PluginBlueprintV100 = (*Webhook)(nil)
var _ core.CloseablePluginTask = (*Webhook)(nil)

type Webhook struct{}

func (plugin Webhook) Description() string {
	return "collect some Webhook data"
}

func (plugin Webhook) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Webhook) SubTaskMetas() []core.SubTaskMeta {
	// TODO add your sub task here
	return []core.SubTaskMeta{}
}

func (plugin Webhook) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.WebhookConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("unable to get Webhook connection by the given connection ID: %v", err)
	}

	return &tasks.WebhookTaskData{
		Options: op,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin Webhook) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/webhook"
}

func (plugin Webhook) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Webhook) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		":connectionId/cicd_pipeline": {
			"POST": api.PostCicdPipeline,
		},
		":connectionId/issue": {
			"POST": api.PostIssue,
		},
		":connectionId/issue/:boardKey/:issueId/close": {
			"POST": api.CloseIssue,
		},
	}
}

func (plugin Webhook) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Webhook) Close(taskCtx core.TaskContext) error {
	_, ok := taskCtx.GetData().(*tasks.WebhookTaskData)
	if !ok {
		return fmt.Errorf("GetData failed when try to close %+v", taskCtx)
	}
	return nil
}
