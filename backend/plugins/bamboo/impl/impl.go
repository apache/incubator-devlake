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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

// make sure interface is implemented
var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginBlueprintV100
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Bamboo)(nil)

type Bamboo struct{}

func (p Bamboo) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Bamboo) Connection() interface{} {
	return &models.BambooConnection{}
}

func (p Bamboo) Scope() interface{} {
	return nil
}

func (p Bamboo) TransformationRule() interface{} {
	return nil
}

func (p Bamboo) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Bamboo) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.BambooConnection{},
		&models.BambooProject{},
		&models.BambooPlan{},
		&models.BambooJob{},
		&models.BambooPlanBuild{},
		&models.BambooPlanBuildVcsRevision{},
		&models.BambooJobBuild{},
	}
}

func (p Bamboo) Description() string {
	return "collect some Bamboo data"
}

func (p Bamboo) SubTaskMetas() []plugin.SubTaskMeta {
	// TODO add your sub task here
	return []plugin.SubTaskMeta{
		tasks.CollectPlanMeta,
		tasks.ExtractPlanMeta,
		tasks.CollectJobMeta,
		tasks.ExtractJobMeta,
		tasks.CollectPlanBuildMeta,
		tasks.ExtractPlanBuildMeta,
		tasks.CollectJobBuildMeta,
		tasks.ExtractJobBuildMeta,
		tasks.CollectDeployMeta,
		tasks.ExtractDeployMeta,
		tasks.CollectDeployBuildMeta,
		tasks.ExtractDeployBuildMeta,

		tasks.ConvertJobBuildsMeta,
		tasks.ConvertPlanBuildsMeta,
		tasks.ConvertPlanVcsMeta,
		tasks.ConvertProjectsMeta,
		tasks.ConvertDeployBuildsMeta,
	}
}

func (p Bamboo) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)

	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.BambooConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Bamboo connection by the given connection ID")
	}

	apiClient, err := tasks.NewBambooApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Bamboo API client instance")
	}

	if op.ProjectKey != "" {
		var scope *models.BambooProject
		// support v100 & advance mode
		// If we still cannot find the record in db, we have to request from remote server and save it to db
		db := taskCtx.GetDal()
		err = db.First(&scope, dal.Where("connection_id = ? AND project_key = ?", op.ConnectionId, op.ProjectKey))
		if err != nil && db.IsErrorNotFound(err) {
			apiProject, err := api.GetApiProject(op.ProjectKey, apiClient)
			if err != nil {
				return nil, err
			}
			logger.Debug(fmt.Sprintf("Current project: %s", apiProject.Key))
			scope = apiProject.ConvertApiScope().(*models.BambooProject)
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(&scope)
			if err != nil {
				return nil, err
			}
		}
		op.TransformationRuleId = scope.TransformationRuleId
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find project: %s", op.ProjectKey))
		}
	}

	if op.BambooTransformationRule == nil && op.TransformationRuleId != 0 {
		var transformationRule models.BambooTransformationRule
		db := taskCtx.GetDal()
		err = db.First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil {
			if db.IsErrorNotFound(err) {
				return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find transformationRules by transformationRuleId [%d]", op.TransformationRuleId))
			}
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find transformationRules by transformationRuleId [%d]", op.TransformationRuleId))
		}
		op.BambooTransformationRule = &transformationRule
	}
	if op.BambooTransformationRule == nil && op.TransformationRuleId == 0 {
		op.BambooTransformationRule = new(models.BambooTransformationRule)
	}
	regexEnricher := helper.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}
	return &tasks.BambooTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Bamboo) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bamboo"
}

func (p Bamboo) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Bamboo) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/transformation_rules": {
			"POST": api.CreateTransformationRule,
			"GET":  api.GetTransformationRuleList,
		},
		"connections/:connectionId/transformation_rules/:id": {
			"PATCH": api.UpdateTransformationRule,
			"GET":   api.GetTransformationRule,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
	}
}

func (p Bamboo) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return nil, errors.Default.New("Bamboo does not support blueprint v100")
}

func (p Bamboo) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.BambooTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
