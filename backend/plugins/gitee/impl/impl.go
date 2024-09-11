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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/gitee/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitee/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.PluginMigration
	plugin.CloseablePluginTask
} = (*Gitee)(nil)

var _ plugin.PluginSource = (*Gitee)(nil)

type Gitee struct{}

func (p Gitee) Connection() dal.Tabler {
	return &models.GiteeConnection{}
}

func (p Gitee) Scope() plugin.ToolLayerScope {
	return &models.GiteeRepo{}
}

func (p Gitee) ScopeConfig() dal.Tabler {
	return nil
}

func (p Gitee) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Gitee) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.GiteeConnection{},
		&models.GiteeAccount{},
		&models.GiteeCommit{},
		&models.GiteeCommitStat{},
		&models.GiteeIssue{},
		&models.GiteeIssueComment{},
		&models.GiteeIssueLabel{},
		&models.GiteePullRequest{},
		&models.GiteePullRequestComment{},
		&models.GiteePullRequestCommit{},
		&models.GiteePullRequestIssue{},
		&models.GiteePullRequestLabel{},
		&models.GiteeRepo{},
		&models.GiteeRepoCommit{},
		&models.GiteeReviewer{},
	}
}

func (p Gitee) Description() string {
	return "To collect and enrich data from Gitee"
}

func (p Gitee) Name() string {
	return "gitee"
}

func (p Gitee) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.CollectCommitsMeta,
		tasks.ExtractCommitsMeta,
		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,
		tasks.CollectApiIssueCommentsMeta,
		tasks.ExtractApiIssueCommentsMeta,
		tasks.CollectApiPullRequestCommitsMeta,
		tasks.ExtractApiPullRequestCommitsMeta,
		tasks.CollectApiPullRequestReviewsMeta,
		tasks.ExtractApiPullRequestReviewsMeta,
		tasks.CollectApiCommitStatsMeta,
		tasks.ExtractApiCommitStatsMeta,
		tasks.EnrichPullRequestIssuesMeta,
		tasks.ConvertRepoMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.ConvertPullRequestCommitsMeta,
		tasks.ConvertPullRequestsMeta,
		tasks.ConvertPullRequestLabelsMeta,
		tasks.ConvertPullRequestIssuesMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertIssueCommentsMeta,
		tasks.ConvertPullRequestCommentsMeta,
		tasks.ConvertPullRequestsMeta,
	}
}

func (p Gitee) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.GiteeOptions
	var err errors.Error
	err = helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}

	if op.Owner == "" {
		return nil, errors.BadInput.New("owner is required for Gitee execution")
	}

	if op.Repo == "" {
		return nil, errors.BadInput.New("repo is required for Gitee execution")
	}

	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}

	connection := &models.GiteeConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)

	if err != nil {
		return nil, err
	}

	err = connectionHelper.FirstById(connection, op.ConnectionId)

	if err != nil {
		return nil, err
	}

	var apiClient *helper.ApiAsyncClient
	syncPolicy := taskCtx.SyncPolicy()
	if !syncPolicy.SkipCollectors {
		newApiClient, err := tasks.NewGiteeApiClient(taskCtx, connection)
		if err != nil {
			return nil, err
		}
		apiClient = newApiClient
	}

	return &tasks.GiteeTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (p Gitee) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitee"
}

func (p Gitee) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Gitee) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
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
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
	}
}

func (p Gitee) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.GiteeTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data != nil && data.ApiClient != nil {
		data.ApiClient.Release()
	}
	return nil
}
