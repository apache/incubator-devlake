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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

var _ plugin.PluginMeta = (*Tapd)(nil)
var _ plugin.PluginInit = (*Tapd)(nil)
var _ plugin.PluginTask = (*Tapd)(nil)
var _ plugin.PluginApi = (*Tapd)(nil)
var _ plugin.PluginModel = (*Tapd)(nil)
var _ plugin.PluginMigration = (*Tapd)(nil)
var _ plugin.CloseablePluginTask = (*Tapd)(nil)

type Tapd struct{}

func (p Tapd) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Tapd) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.TapdAccount{},
		&models.TapdBug{},
		&models.TapdBugChangelog{},
		&models.TapdBugChangelogItem{},
		&models.TapdBugCommit{},
		&models.TapdBugCustomFields{},
		&models.TapdBugLabel{},
		&models.TapdBugStatus{},
		&models.TapdConnection{},
		&models.TapdIssue{},
		&models.TapdIteration{},
		&models.TapdIterationBug{},
		&models.TapdIterationStory{},
		&models.TapdIterationTask{},
		&models.TapdStory{},
		&models.TapdStoryBug{},
		&models.TapdStoryCategory{},
		&models.TapdStoryChangelog{},
		&models.TapdStoryChangelogItem{},
		&models.TapdStoryCommit{},
		&models.TapdStoryCustomFields{},
		&models.TapdStoryLabel{},
		&models.TapdStoryStatus{},
		&models.TapdSubWorkspace{},
		&models.TapdTask{},
		&models.TapdTaskChangelog{},
		&models.TapdTaskChangelogItem{},
		&models.TapdTaskCommit{},
		&models.TapdTaskCustomFields{},
		&models.TapdTaskLabel{},
		&models.TapdWorkSpaceBug{},
		&models.TapdWorkSpaceStory{},
		&models.TapdWorkSpaceTask{},
		&models.TapdWorklog{},
		&models.TapdWorkspace{},
		&models.TapdWorkspaceIteration{},
	}
}

func (p Tapd) Description() string {
	return "To collect and enrich data from Tapd"
}

func (p Tapd) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectCompanyMeta,
		tasks.ExtractCompanyMeta,
		tasks.CollectSubWorkspaceMeta,
		tasks.ExtractSubWorkspaceMeta,
		tasks.CollectWorkitemTypesMeta,
		tasks.ExtractWorkitemTypesMeta,
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
		tasks.CollectStoryStatusLastStepMeta,
		tasks.EnrichStoryStatusLastStepMeta,
		tasks.CollectBugStatusMeta,
		tasks.ExtractBugStatusMeta,
		tasks.CollectBugStatusLastStepMeta,
		tasks.EnrichBugStatusLastStepMeta,
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
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
		tasks.ConvertAccountsMeta,
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

func (p Tapd) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	var op tasks.TapdOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
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
	if connection.RateLimitPerHour == 0 {
		connection.RateLimitPerHour = 3600
	}
	tapdApiClient, err := tasks.NewTapdApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to create tapd api client")
	}
	cstZone, err1 := time.LoadLocation("Asia/Shanghai")
	if err1 != nil {
		return nil, errors.Default.Wrap(err1, "fail to get CST Location")
	}
	op.CstZone = cstZone
	taskData := &tasks.TapdTaskData{
		Options:    &op,
		ApiClient:  tapdApiClient,
		Connection: connection,
	}
	var timeAfter time.Time
	if op.TimeAfter != "" {
		timeAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `timeAfter`")
		}
	}
	if !timeAfter.IsZero() {
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data updated timeAfter %s", timeAfter)
	}
	return taskData, nil
}

func (p Tapd) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/tapd"
}

func (p Tapd) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Tapd) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
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
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (p Tapd) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.TapdTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
