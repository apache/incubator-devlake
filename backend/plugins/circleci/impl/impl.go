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
	"github.com/apache/incubator-devlake/plugins/circleci/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/apache/incubator-devlake/plugins/circleci/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/circleci/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Circleci)(nil)
var _ plugin.PluginInit = (*Circleci)(nil)
var _ plugin.PluginTask = (*Circleci)(nil)
var _ plugin.PluginApi = (*Circleci)(nil)
var _ plugin.CloseablePluginTask = (*Circleci)(nil)

type Circleci struct{}

// Name implements plugin.PluginMeta.
func (Circleci) Name() string {
	return "circleci"
}

func (p Circleci) Description() string {
	return "collect some Circleci data"
}

func (p Circleci) Init(br context.BasicRes) errors.Error {
	api.Init(br, &p)
	return nil
}

func (p Circleci) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.CircleciConnection{},
		&models.CircleciAccount{},
		&models.CircleciProject{},
		&models.CircleciPipeline{},
		&models.CircleciWorkflow{},
		&models.CircleciJob{},
		&models.CircleciScopeConfig{},
	}
}

func (p Circleci) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
		tasks.ConvertAccountsMeta,
		tasks.CollectProjectsMeta,
		tasks.ConvertProjectsMeta,
		tasks.CollectPipelinesMeta,
		tasks.ExtractPipelinesMeta,
		tasks.CollectWorkflowsMeta,
		tasks.ExtractWorkflowsMeta,
		tasks.CollectJobsMeta,
		tasks.ExtractJobsMeta,
		tasks.ConvertJobsMeta,
		tasks.ConvertWorkflowsMeta,
	}
}

func (p Circleci) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Circleci) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.CircleciConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Circleci connection by the given connection ID")
	}

	apiClient, err := tasks.NewCircleciApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Circleci API client instance")
	}
	project := &models.CircleciProject{}
	err = taskCtx.GetDal().First(project, dal.Where("slug = ?", op.ProjectSlug))
	if err != nil {
		return nil, errors.Default.Wrap(err, "project not found")
	}

	taskData := &tasks.CircleciTaskData{
		Options:   op,
		ApiClient: apiClient,
		Project:   project,
	}
	// fallback to project scope config
	if op.ScopeConfigId == 0 {
		project := &models.CircleciProject{}
		err = taskCtx.GetDal().First(project, dal.Where("slug = ?", op.ProjectSlug))
		if err != nil {
			return nil, err
		}
		op.ScopeConfigId = project.ScopeConfigId
	}
	// fallback to given scope config
	if op.ScopeConfig == nil && op.ScopeConfigId != 0 {
		var scopeConfig models.CircleciScopeConfig
		db := taskCtx.GetDal()
		err = db.First(&scopeConfig, dal.Where("id = ?", op.ScopeConfigId))
		if err != nil && !db.IsErrorNotFound(err) {
			return nil, errors.BadInput.Wrap(err, "fail to load scopeConfig")
		}
		op.ScopeConfig = &scopeConfig
	}
	if op.ScopeConfig == nil && op.ScopeConfigId == 0 {
		op.ScopeConfig = new(models.CircleciScopeConfig)
	}
	regexEnricher := helper.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.ScopeConfig.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ScopeConfig.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}
	taskData.RegexEnricher = regexEnricher
	return taskData, nil
}

// RootPkgPath PkgPath information lost when compiled as plugin(.so)
func (p Circleci) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/circleci"
}

func (p Circleci) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Circleci) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
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

func (p Circleci) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.CircleciTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
