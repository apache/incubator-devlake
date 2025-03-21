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
	"context"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/argo/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	connection, err := dsHelper.ConnSrv.FindByPk(connectionId)
	if err != nil {
		return nil, nil, err
	}
	scopeDetails, err := dsHelper.ScopeSrv.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}

	_, err = helper.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, nil, err
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
	scopeDetails []*srvhelper.ScopeDetail[models.Project, models.ArgoScopeConfig],
	connection *models.ArgoConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		argo, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}
		// construct task options for Jira
		task, err := helper.MakePipelinePlanTask(
			"argo",
			subtaskMetas,
			scopeConfig.Entities,
			tasks.BitbucketOptions{
				ConnectionId: argo.ConnectionId,
				FullName:     argo.ScopeName(),
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
	scopeDetails []*srvhelper.ScopeDetail[models.Project, models.ArgoScopeConfig],
	connection *models.ArgoConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	idgen := didgen.NewDomainIdGenerator(&models.Project{})
	for _, scopeDetail := range scopeDetails {
		scope, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		id := idgen.Generate(connection.ID, scope.ProjectId)

		// add cicd_scope to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopes = append(scopes, &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: id,
				},
				Name: scope.ScopeId(),
			})
		}
	}
	return scopes, nil
}
