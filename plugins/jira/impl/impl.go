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
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Jira)(nil)
var _ core.PluginInit = (*Jira)(nil)
var _ core.PluginTask = (*Jira)(nil)
var _ core.PluginApi = (*Jira)(nil)
var _ core.Migratable = (*Jira)(nil)
var _ core.PluginBlueprintV100 = (*Jira)(nil)

type Jira struct{}

func (plugin Jira) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectStatusMeta,
		tasks.ExtractStatusMeta,

		tasks.CollectProjectsMeta,
		tasks.ExtractProjectsMeta,

		tasks.CollectBoardMeta,
		tasks.ExtractBoardMeta,

		tasks.CollectIssuesMeta,
		tasks.ExtractIssuesMeta,

		tasks.CollectChangelogsMeta,
		tasks.ExtractChangelogsMeta,

		tasks.CollectUsersMeta,

		tasks.CollectWorklogsMeta,
		tasks.ExtractWorklogsMeta,

		tasks.CollectRemotelinksMeta,
		tasks.ExtractRemotelinksMeta,

		tasks.CollectSprintsMeta,
		tasks.ExtractSprintsMeta,

		tasks.ConvertBoardMeta,

		tasks.ConvertIssuesMeta,

		tasks.ConvertWorklogsMeta,

		tasks.ConvertChangelogsMeta,

		tasks.ConvertSprintsMeta,
		tasks.ConvertSprintIssuesMeta,

		tasks.ConvertIssueCommitsMeta,
		tasks.ConvertIssueRepoCommitsMeta,

		tasks.ExtractUsersMeta,
		tasks.ConvertUsersMeta,
	}
}

func (plugin Jira) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.JiraOptions
	var err error
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}
	connection := &models.JiraConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	var since time.Time
	if op.Since != "" {
		since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return nil, fmt.Errorf("invalid value for `since`: %w", err)
		}
	}
	jiraApiClient, err := tasks.NewJiraApiClient(taskCtx, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to create jira api client: %w", err)
	}
	info, code, err := tasks.GetJiraServerInfo(jiraApiClient)
	if err != nil || code != http.StatusOK || info == nil {
		return nil, fmt.Errorf("fail to get server info: error:[%s] code:[%d]", err, code)
	}
	taskData := &tasks.JiraTaskData{
		Options:        &op,
		ApiClient:      jiraApiClient,
		JiraServerInfo: *info,
	}
	if !since.IsZero() {
		taskData.Since = &since
		logger.Debug("collect data updated since %s", since)
	}
	return taskData, nil
}

func (plugin Jira) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Jira) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jira"
}

func (plugin Jira) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Jira) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"echo": {
			"POST": func(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
				return &core.ApiResourceOutput{Body: input.Body}, nil
			},
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}
