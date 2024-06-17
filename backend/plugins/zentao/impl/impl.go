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
	"github.com/apache/incubator-devlake/core/runner"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/plugins/zentao/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
	"github.com/spf13/viper"
)

// make sure interface is implemented

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	//plugin.CompositePluginBlueprintV200
	plugin.PluginModel
	plugin.PluginSource
	plugin.CloseablePluginTask
} = (*Zentao)(nil)

type Zentao struct{}

func (p Zentao) Description() string {
	return "collect some Zentao data"
}

func (p Zentao) Name() string {
	return "zentao"
}

func (p Zentao) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Zentao) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ZentaoAccount{},
		&models.ZentaoBug{},
		&models.ZentaoBugCommit{},
		&models.ZentaoChangelog{},
		&models.ZentaoChangelogDetail{},
		&models.ZentaoDepartment{},
		&models.ZentaoExecution{},
		&models.ZentaoProduct{},
		&models.ZentaoProject{},
		&models.ZentaoRemoteDbAction{},
		&models.ZentaoRemoteDbActionHistory{},
		&models.ZentaoRemoteDbHistory{},
		&models.ZentaoStory{},
		&models.ZentaoStoryCommit{},
		&models.ZentaoStoryRepoCommit{},
		&models.ZentaoTask{},
		&models.ZentaoTaskCommit{},
		&models.ZentaoTaskRepoCommit{},
		&models.ZentaoBugRepoCommit{},
		&models.ZentaoConnection{},
		&models.ZentaoScopeConfig{},
		&models.ZentaoExecutionStory{},
		&models.ZentaoExecutionSummary{},
		&models.ZentaoProductSummary{},
		&models.ZentaoProjectStory{},
	}
}

func (p Zentao) Connection() dal.Tabler {
	return &models.ZentaoConnection{}
}

func (p Zentao) Scope() plugin.ToolLayerScope {
	return &models.ZentaoProject{}
}

func (p Zentao) ScopeConfig() dal.Tabler {
	return &models.ZentaoScopeConfig{}
}

func (p Zentao) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ConvertProjectMeta,

		// both
		tasks.CollectAccountMeta,
		tasks.ExtractAccountMeta,
		tasks.ConvertAccountMeta,

		tasks.CollectDepartmentMeta,
		tasks.ExtractDepartmentMeta,

		//project
		tasks.CollectExecutionSummaryMeta,
		tasks.ExtractExecutionSummaryMeta,

		tasks.CollectExecutionSummaryDevMeta,
		tasks.ExtractExecutionSummaryDevMeta,

		tasks.CollectExecutionMeta,
		tasks.ExtractExecutionMeta,
		tasks.ConvertExecutionMeta,

		tasks.CollectTaskMeta,
		tasks.ExtractTaskMeta,
		tasks.ConvertTaskMeta,

		tasks.CollectTaskCommitsMeta,
		tasks.ExtractTaskCommitsMeta,
		tasks.CollectTaskRepoCommitsMeta,
		tasks.ExtractTaskRepoCommitsMeta,
		tasks.ConvertTaskRepoCommitsMeta,

		// product
		tasks.CollectStoryMeta,
		tasks.ExtractStoryMeta,
		tasks.ConvertStoryMeta,
		tasks.ConvertExecutionStoryMeta,

		tasks.CollectBugMeta,
		tasks.ExtractBugMeta,
		tasks.ConvertBugMeta,

		tasks.CollectStoryCommitsMeta,
		tasks.ExtractStoryCommitsMeta,
		tasks.CollectStoryRepoCommitsMeta,
		tasks.ExtractStoryRepoCommitsMeta,
		tasks.ConvertStoryRepoCommitsMeta,

		tasks.CollectBugCommitsMeta,
		tasks.ExtractBugCommitsMeta,
		tasks.CollectBugRepoCommitsMeta,
		tasks.ExtractBugRepoCommitsMeta,
		tasks.ConvertBugRepoCommitsMeta,

		tasks.DBGetChangelogMeta,
		tasks.ConvertChangelogMeta,
	}
}

func (p Zentao) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not decode Zentao options")
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.ZentaoConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Zentao connection by the given connection ID: %v")
	}

	apiClient, err := tasks.NewZentaoApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Zentao API client instance: %v")
	}

	if op.ScopeConfig == nil && op.ScopeConfigId != 0 {
		err = taskCtx.GetDal().First(&op.ScopeConfig, dal.Where("id = ?", op.ScopeConfigId))
		if err != nil && taskCtx.GetDal().IsErrorNotFound(err) {
			return nil, errors.BadInput.Wrap(err, "fail to load scope config from database")
		}
	}

	data := &tasks.ZentaoTaskData{
		Options:      op,
		ApiClient:    apiClient,
		Stories:      map[int64]struct{}{},
		Tasks:        map[int64]struct{}{},
		Bugs:         map[int64]struct{}{},
		AccountCache: tasks.NewAccountCache(taskCtx.GetDal(), op.ConnectionId),
	}

	if connection.DbUrl != "" {
		if connection.DbLoggingLevel == "" {
			connection.DbLoggingLevel = taskCtx.GetConfig("DB_LOGGING_LEVEL")
		}

		if connection.DbIdleConns == 0 {
			connection.DbIdleConns = taskCtx.GetConfigReader().GetInt("DB_IDLE_CONNS")
		}

		if connection.DbMaxConns == 0 {
			connection.DbMaxConns = taskCtx.GetConfigReader().GetInt("DB_MAX_CONNS")
		}

		v := viper.New()
		v.Set("DB_URL", connection.DbUrl)
		v.Set("DB_LOGGING_LEVEL", connection.DbLoggingLevel)
		v.Set("DB_IDLE_CONNS", connection.DbIdleConns)
		v.Set("DbMaxConns", connection.DbMaxConns)

		rgorm, err := runner.NewGormDb(v, taskCtx.GetLogger())
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("failed to connect to the zentao remote databases %s", connection.DbUrl))
		}

		data.RemoteDb = dalgorm.NewDalgorm(rgorm)
	}

	return data, nil
}

// RootPkgPath information lost when compiled as plugin(.so)
func (p Zentao) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/zentao"
}

func (p Zentao) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Zentao) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/scopes": {
			"PUT": api.PutScopes,
			"GET": api.GetScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteProjectScope,
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
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/proxy/*path": {
			"GET": api.Proxy,
		},
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

func (p Zentao) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Zentao) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ZentaoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
