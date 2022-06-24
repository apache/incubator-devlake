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
	"time"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/tapd/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Tapd)(nil)
var _ core.PluginInit = (*Tapd)(nil)
var _ core.PluginTask = (*Tapd)(nil)
var _ core.PluginApi = (*Tapd)(nil)
var _ core.Migratable = (*Tapd)(nil)

type Tapd struct{}

func (plugin Tapd) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Tapd) Description() string {
	return "To collect and enrich data from Tapd"
}

func (plugin Tapd) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectCompanyMeta,
		tasks.ExtractCompanyMeta,
		tasks.CollectSubWorkspaceMeta,
		tasks.ExtractSubWorkspaceMeta,
		tasks.CollectStoryCustomFieldsMeta,
		tasks.ExtractStoryCustomFieldsMeta,
		tasks.CollectTaskCustomFieldsMeta,
		tasks.ExtractTaskCustomFieldsMeta,
		tasks.CollectBugCustomFieldsMeta,
		tasks.ExtractBugCustomFieldsMeta,
		tasks.CollectStoryCategoriesMeta,
		tasks.ExtractStoryCategoriesMeta,
		tasks.CollectStoryStatusMeta,
		tasks.ExtractStoryStatusMeta,
		tasks.CollectBugStatusMeta,
		tasks.ExtractBugStatusMeta,
		tasks.CollectUserMeta,
		tasks.ExtractUserMeta,
		tasks.CollectIterationMeta,
		tasks.ExtractIterationMeta,
		tasks.CollectStoryMeta,
		tasks.CollectBugMeta,
		tasks.CollectTaskMeta,
		tasks.ExtractStoryMeta,
		tasks.ExtractBugMeta,
		tasks.ExtractTaskMeta,
		tasks.CollectBugChangelogMeta,
		tasks.ExtractBugChangelogMeta,
		tasks.CollectStoryChangelogMeta,
		tasks.ExtractStoryChangelogMeta,
		tasks.CollectTaskChangelogMeta,
		tasks.ExtractTaskChangelogMeta,
		tasks.CollectWorklogMeta,
		tasks.ExtractWorklogMeta,
		tasks.CollectBugCommitMeta,
		tasks.ExtractBugCommitMeta,
		tasks.CollectStoryCommitMeta,
		tasks.ExtractStoryCommitMeta,
		tasks.CollectTaskCommitMeta,
		tasks.ExtractTaskCommitMeta,
		tasks.CollectStoryBugMeta,
		tasks.ExtractStoryBugsMeta,
		tasks.ConvertSubWorkspaceMeta,
		tasks.ConvertUserMeta,
		tasks.ConvertIterationMeta,
		tasks.ConvertStoryMeta,
		tasks.ConvertBugMeta,
		tasks.ConvertTaskMeta,
		tasks.ConvertWorklogMeta,
		tasks.ConvertBugChangelogMeta,
		tasks.ConvertStoryChangelogMeta,
		tasks.ConvertTaskChangelogMeta,
		tasks.ConvertBugCommitMeta,
		tasks.ConvertStoryCommitMeta,
		tasks.ConvertTaskCommitMeta,
		tasks.ConvertStoryLabelsMeta,
		tasks.ConvertTaskLabelsMeta,
		tasks.ConvertBugLabelsMeta,
	}
}

func (plugin Tapd) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.TapdOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}
	connection := &models.TapdConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
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
	if connection.RateLimit == 0 {
		connection.RateLimit = 6000
	}
	tapdApiClient, err := tasks.NewTapdApiClient(taskCtx, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to create tapd api client: %w", err)
	}
	taskData := &tasks.TapdTaskData{
		Options:    &op,
		ApiClient:  tapdApiClient,
		Connection: connection,
	}
	if !since.IsZero() {
		taskData.Since = &since
	}
	tasks.WorkspaceIdGen = didgen.NewDomainIdGenerator(&models.TapdWorkspace{})
	tasks.IssueIdGen = didgen.NewDomainIdGenerator(&models.TapdIssue{})
	tasks.IterIdGen = didgen.NewDomainIdGenerator(&models.TapdIteration{})

	return taskData, nil
}

func (plugin Tapd) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/tapd"
}

func (plugin Tapd) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Tapd) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
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
	}
}
