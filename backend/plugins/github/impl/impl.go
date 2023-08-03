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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
} = (*Github)(nil)

var sortedSubtaskMetas []plugin.SubTaskMeta

type Github struct{}

func init() {
	var err error
	// check subtask meta loop and gen subtask list when plugin init
	sortedSubtaskMetas, err = sorter.NewTableSorter(tasks.SubTaskMetaList).Sort()
	if err != nil {
		panic(err)
	}
}

func (p Github) Connection() dal.Tabler {
	return &models.GithubConnection{}
}

func (p Github) Scope() plugin.ToolLayerScope {
	return &models.GithubRepo{}
}

func (p Github) ScopeConfig() dal.Tabler {
	return &models.GithubScopeConfig{}
}

func (p Github) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Github) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.GithubConnection{},
		&models.GithubAccount{},
		&models.GithubAccountOrg{},
		&models.GithubCommit{},
		&models.GithubCommitStat{},
		&models.GithubIssue{},
		&models.GithubIssueComment{},
		&models.GithubIssueEvent{},
		&models.GithubIssueLabel{},
		&models.GithubJob{},
		&models.GithubMilestone{},
		&models.GithubPrComment{},
		&models.GithubPrCommit{},
		&models.GithubPrIssue{},
		&models.GithubPrLabel{},
		&models.GithubPrReview{},
		&models.GithubPullRequest{},
		&models.GithubRepo{},
		&models.GithubRepoAccount{},
		&models.GithubRepoCommit{},
		&models.GithubReviewer{},
		&models.GithubRun{},
		&models.GithubIssueAssignee{},
		&models.GithubScopeConfig{},
	}
}

func (p Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (p Github) Name() string {
	return "github"
}

func (p Github) SubTaskMetas() []plugin.SubTaskMeta {
	return sortedSubtaskMetas
}

func (p Github) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.GithubConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get github connection by the given connection ID")
	}
	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get github API client instance")
	}
	err = EnrichOptions(taskCtx, op, apiClient.ApiClient)
	if err != nil {
		return nil, err
	}

	regexEnricher := helper.NewRegexEnricher()
	if err = regexEnricher.TryAdd(devops.DEPLOYMENT, op.ScopeConfig.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err = regexEnricher.TryAdd(devops.PRODUCTION, op.ScopeConfig.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}

	taskData := &tasks.GithubTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}

	if op.TimeAfter != "" {
		var timeAfter time.Time
		timeAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `timeAfter`")
		}
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data updated timeAfter %s", timeAfter)
	}
	return taskData, nil
}

func (p Github) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/github"
}

func (p Github) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Github) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.UpdateScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.CreateScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:id": {
			"PATCH":  api.UpdateScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (p Github) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Github) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.GithubTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx plugin.TaskContext,
	op *tasks.GithubOptions,
	apiClient *helper.ApiClient) errors.Error {
	var githubRepo models.GithubRepo
	// validate the op and set name=owner/repo if this is from advanced mode or bpV100
	err := tasks.ValidateTaskOptions(op)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	// for advanced mode or others which we only have name, for bp v200, we have githubId
	err = taskCtx.GetDal().First(&githubRepo, dal.Where(
		"connection_id = ? AND( full_name = ? OR github_id = ?)",
		op.ConnectionId, op.Name, op.GithubId))
	if err == nil {
		op.Name = githubRepo.FullName
		op.GithubId = githubRepo.GithubId
		if op.ScopeConfigId == 0 {
			op.ScopeConfigId = githubRepo.ScopeConfigId
		}
	} else {
		if taskCtx.GetDal().IsErrorNotFound(err) && op.Name != "" {
			var repo *tasks.GithubApiRepo
			repo, err = api.MemorizedGetApiRepo(repo, op, apiClient)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Current repo: %s", repo.FullName))
			scope := convertApiRepoToScope(repo, op.ConnectionId)
			err = taskCtx.GetDal().CreateIfNotExist(scope)
			if err != nil {
				return err
			}
			op.GithubId = repo.GithubId
		} else {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", op.Name))
		}
	}
	// Set GithubScopeConfig if it's nil, this has lower priority
	if op.ScopeConfig == nil && op.ScopeConfigId != 0 {
		var scopeConfig models.GithubScopeConfig
		db := taskCtx.GetDal()
		err = db.First(&scopeConfig, dal.Where("id = ?", githubRepo.ScopeConfigId))
		if err != nil && !db.IsErrorNotFound(err) {
			return errors.BadInput.Wrap(err, "fail to get scopeConfig")
		}
		op.ScopeConfig = &scopeConfig
	}
	if op.ScopeConfig == nil && op.ScopeConfigId == 0 {
		op.ScopeConfig = new(models.GithubScopeConfig)
	}
	return err
}

func convertApiRepoToScope(repo *tasks.GithubApiRepo, connectionId uint64) *models.GithubRepo {
	var scope models.GithubRepo
	scope.ConnectionId = connectionId
	scope.GithubId = repo.GithubId
	scope.CreatedDate = repo.CreatedAt.ToNullableTime()
	scope.Language = repo.Language
	scope.Description = repo.Description
	scope.HTMLUrl = repo.HTMLUrl
	scope.Name = repo.Name
	scope.FullName = repo.FullName
	scope.CloneUrl = repo.CloneUrl
	return &scope
}
