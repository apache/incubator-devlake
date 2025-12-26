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
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
	plugin.CloseablePluginTask
	plugin.DataSourcePluginBlueprintV200
	plugin.PluginSource
} = (*Bitbucket)(nil)

type Bitbucket struct{}

func (p Bitbucket) Connection() dal.Tabler {
	return &models.BitbucketConnection{}
}

func (p Bitbucket) Scope() plugin.ToolLayerScope {
	return &models.BitbucketRepo{}
}

func (p Bitbucket) ScopeConfig() dal.Tabler {
	return &models.BitbucketScopeConfig{}
}

func (p Bitbucket) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Bitbucket) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.BitbucketConnection{},
		&models.BitbucketAccount{},
		&models.BitbucketCommit{},
		&models.BitbucketPullRequest{},
		&models.BitbucketIssue{},
		&models.BitbucketPrComment{},
		&models.BitbucketIssueComment{},
		&models.BitbucketPipeline{},
		&models.BitbucketRepo{},
		&models.BitbucketRepoCommit{},
		&models.BitbucketDeployment{},
		&models.BitbucketPipelineStep{},
		&models.BitbucketPrCommit{},
		&models.BitbucketPrReviewer{},
		&models.BitbucketScopeConfig{},
	}
}

func (p Bitbucket) Description() string {
	return "To collect and enrich data from Bitbucket"
}

func (p Bitbucket) Name() string {
	return "bitbucket"
}

func (p Bitbucket) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,

		tasks.CollectApiPrCommentsMeta,
		tasks.ExtractApiPrCommentsMeta,

		tasks.CollectApiPrCommitsMeta,
		tasks.ExtractApiPrCommitsMeta,

		tasks.CollectApiCommitsMeta,
		tasks.ExtractApiCommitsMeta,

		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,

		tasks.CollectApiIssueCommentsMeta,
		tasks.ExtractApiIssueCommentsMeta,

		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,

		tasks.CollectApiDeploymentsMeta,
		tasks.ExtractApiDeploymentsMeta,

		// must run after deployment to match
		tasks.CollectPipelineStepsMeta,
		tasks.ExtractPipelineStepsMeta,

		tasks.ConvertRepoMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertPullRequestsMeta,
		tasks.ConvertPrReviewersMeta,
		tasks.ConvertPrCommentsMeta,
		tasks.ConvertPrCommitsMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueCommentsMeta,
		tasks.ConvertPipelineMeta,
		tasks.ConvertPipelineStepMeta,
		tasks.ConvertiDeploymentMeta,
	}
}

func (p Bitbucket) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.BitbucketConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket connection by the given connection ID")
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket API client instance")
	}
	err = EnrichOptions(taskCtx, op, apiClient.ApiClient)
	if err != nil {
		return nil, err
	}

	regexEnricher := helper.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}
	taskData := &tasks.BitbucketTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}

	return taskData, nil
}

func (p Bitbucket) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bitbucket/" // the "/" fixes an issue where records from "bitbucket_server" are counted as "bitbucket" records and vice versa
}

func (p Bitbucket) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Bitbucket) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Bitbucket) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
		"connections/:connectionId/scopes/*scopeId": {
			// Behind 'GetScopeDispatcher', there are two paths so far:
			// GetScopeLatestSyncState "connections/:connectionId/scopes/:scopeId/latest-sync-state"
			// GetScope "connections/:connectionId/scopes/:scopeId"
			// Because there may be slash in scopeId, so we handle it manually.
			"GET":    api.GetScopeDispatcher,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopes,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.PostScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId": {
			"PATCH":  api.PatchScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

func (p Bitbucket) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.BitbucketTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx plugin.TaskContext,
	op *tasks.BitbucketOptions,
	apiClient *helper.ApiClient) errors.Error {
	var repo models.BitbucketRepo
	// validate the op and set name=owner/repo if this is from advanced mode or bpV100
	err := tasks.ValidateTaskOptions(op)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	err = taskCtx.GetDal().First(&repo, dal.Where(
		"connection_id = ? AND bitbucket_id = ?",
		op.ConnectionId, op.FullName))
	if err == nil {
		if op.ScopeConfigId == 0 {
			op.ScopeConfigId = repo.ScopeConfigId
		}
	} else {
		if taskCtx.GetDal().IsErrorNotFound(err) && op.FullName != "" {
			var repo *models.BitbucketApiRepo
			repo, err = tasks.GetApiRepo(op, apiClient)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Current repo: %s", repo.FullName))
			scope := repo.ConvertApiScope()
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(scope)
			if err != nil {
				return err
			}
		} else {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", op.FullName))
		}
	}
	// Set GithubScopeConfig if it's nil, this has lower priority
	if op.BitbucketScopeConfig == nil && op.ScopeConfigId != 0 {
		var scopeConfig models.BitbucketScopeConfig
		db := taskCtx.GetDal()
		err = db.First(&scopeConfig, dal.Where("id = ?", repo.ScopeConfigId))
		if err != nil && !db.IsErrorNotFound(err) {
			return errors.BadInput.Wrap(err, "fail to get scopeConfig")
		}
		op.BitbucketScopeConfig = &scopeConfig
	}
	if op.BitbucketScopeConfig == nil && op.ScopeConfigId == 0 {
		op.BitbucketScopeConfig = new(models.BitbucketScopeConfig)
	}
	return err
}
