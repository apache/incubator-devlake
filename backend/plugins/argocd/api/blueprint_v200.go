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

package api

import (
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/apache/incubator-devlake/plugins/argocd/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	connection, err := dsHelper.ConnSrv.FindByPk(connectionId)
	if err != nil {
		return nil, nil, err
	}
	scopeDetails, err := dsHelper.ScopeApi.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		// attempt auto-create missing scopes for blueprint (only name known)
		cfg, _ := CreateDefaultScopeConfig(connectionId)
		for _, bs := range bpScopes {
			if _, findErr := dsHelper.ScopeSrv.ModelSrvHelper.FindByPk(connectionId, bs.ScopeId); findErr != nil {
				app := &models.ArgocdApplication{}
				app.ConnectionId = connectionId
				app.Name = bs.ScopeId
				if cfg != nil {
					app.ScopeConfigId = cfg.ID
				}
				_ = dsHelper.ScopeSrv.ModelSrvHelper.CreateOrUpdate(app)
			}
		}
		scopeDetails, err = dsHelper.ScopeApi.MapScopeDetails(connectionId, bpScopes)
		if err != nil {
			return nil, nil, err
		}
	}

	plan, err := makeDataSourcePipelinePlanV200(subtaskMetas, scopeDetails, connection)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(scopeDetails, connection)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopeDetails []*srvhelper.ScopeDetail[models.ArgocdApplication, models.ArgocdScopeConfig],
	connection *models.ArgocdConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		application := scopeDetail.Scope
		scopeConfig := scopeDetail.ScopeConfig
		if scopeConfig == nil {
			scopeConfig = &models.ArgocdScopeConfig{}
		}
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}
		scopeConfigId := application.ScopeConfigId
		if scopeConfig != nil && scopeConfig.ID != 0 {
			scopeConfigId = scopeConfig.ID
		}
		task, err := api.MakePipelinePlanTask(
			"argocd",
			subtaskMetas,
			scopeConfig.Entities,
			tasks.ArgocdOptions{
				ConnectionId:    connection.ID,
				ApplicationName: application.Name,
				ScopeConfigId:   scopeConfigId,
			},
		)
		if err != nil {
			return nil, err
		}
		stage = append(stage, task)
		plan[i] = stage
	}

	return plan, nil
}

func makeScopesV200(
	scopeDetails []*srvhelper.ScopeDetail[models.ArgocdApplication, models.ArgocdScopeConfig],
	connection *models.ArgocdConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0, len(scopeDetails)*2)
	idGen := didgen.NewDomainIdGenerator(&models.ArgocdApplication{})
	for _, scopeDetail := range scopeDetails {
		application, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		scopes = append(scopes, application)
		if scopeConfig == nil {
			continue
		}
		entities := scopeConfig.Entities
		if len(entities) == 0 {
			entities = plugin.DOMAIN_TYPES
		}
		if utils.StringsContains(entities, plugin.DOMAIN_TYPE_CICD) {
			scopeId := idGen.Generate(connection.ID, application.Name)
			scopes = append(scopes, &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{Id: scopeId},
				Name:         application.Name,
			})
		}
	}
	return scopes, nil
}
