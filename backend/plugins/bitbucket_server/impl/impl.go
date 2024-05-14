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
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/tasks"
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
} = (*BitbucketServer)(nil)

type BitbucketServer struct{}

func (p BitbucketServer) Connection() dal.Tabler {
	return &models.BitbucketServerConnection{}
}

func (p BitbucketServer) Scope() plugin.ToolLayerScope {
	return &models.BitbucketServerRepo{}
}

func (p BitbucketServer) ScopeConfig() dal.Tabler {
	return &models.BitbucketServerScopeConfig{}
}

func (p BitbucketServer) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p BitbucketServer) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.BitbucketServerConnection{},
		&models.BitbucketServerUser{},
		&models.BitbucketServerPullRequest{},
		&models.BitbucketServerPrComment{},
		&models.BitbucketServerRepo{},
		&models.BitbucketServerPrCommit{},
		&models.BitbucketServerScopeConfig{},
	}
}

func (p BitbucketServer) Description() string {
	return "To collect and enrich data from Bitbucket Server"
}

func (p BitbucketServer) Name() string {
	return "bitbucket_server"
}

func (p BitbucketServer) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,

		tasks.CollectApiPrActivitiesMeta,
		tasks.ExtractApiPrActivitiesMeta,

		tasks.CollectApiPrCommitsMeta,
		tasks.ExtractApiPrCommitsMeta,

		tasks.ConvertRepoMeta, // ?
		tasks.ConvertPullRequestsMeta,

		tasks.ConvertPrCommentsMeta,
		tasks.ConvertPrCommitsMeta,

		tasks.ConvertUsersMeta,
	}
}

func (p BitbucketServer) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
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
	connection := &models.BitbucketServerConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket server connection by the given connection ID")
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket server API client instance")
	}
	err = EnrichOptions(taskCtx, op, apiClient.ApiClient)
	if err != nil {
		return nil, err
	}

	regexEnricher := helper.NewRegexEnricher()
	taskData := &tasks.BitbucketServerTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}

	return taskData, nil
}

func (p BitbucketServer) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bitbucket_server/" // the "/" fixes an issue where records from "bitbucket_server" are counted as "bitbucket" records and vice versa
}

func (p BitbucketServer) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p BitbucketServer) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p BitbucketServer) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
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
		"connections/:connectionId/scopes/*scopeId": {
			// Behind 'GetScopeDispatcher', there are two paths so far:
			// GetScopeLatestSyncState "connections/:connectionId/scopes/:scopeId/latest-sync-state"
			// GetScope "connections/:connectionId/scopes/:scopeId"
			// Because there may be slash in scopeId, so we handle it manually.
			"GET":    api.GetScopeDispatcher,
			"PATCH":  api.UpdateScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.CreateScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/*scopeConfigId": {
			"PATCH":  api.UpdateScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

func (p BitbucketServer) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.BitbucketServerTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx plugin.TaskContext,
	op *tasks.BitbucketServerOptions,
	apiClient *helper.ApiClient) errors.Error {
	var repo models.BitbucketServerRepo
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
			var repo *models.BitbucketServerApiRepo
			repo, err = tasks.GetApiRepo(op, apiClient)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Current repo: %s", repo.Slug))
			scope := repo.ConvertApiScope().(*models.BitbucketServerRepo)
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(scope)
			if err != nil {
				return err
			}
		} else {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", op.FullName))
		}
	}
	// Set scope config if it's nil, this has lower priority
	if op.BitbucketServerScopeConfig == nil && op.ScopeConfigId != 0 {
		var scopeConfig models.BitbucketServerScopeConfig
		db := taskCtx.GetDal()
		err = db.First(&scopeConfig, dal.Where("id = ?", repo.ScopeConfigId))
		if err != nil && !db.IsErrorNotFound(err) {
			return errors.BadInput.Wrap(err, "fail to get scopeConfig")
		}
		op.BitbucketServerScopeConfig = &scopeConfig
	}
	if op.BitbucketServerScopeConfig == nil && op.ScopeConfigId == 0 {
		op.BitbucketServerScopeConfig = new(models.BitbucketServerScopeConfig)
	}
	return err
}
