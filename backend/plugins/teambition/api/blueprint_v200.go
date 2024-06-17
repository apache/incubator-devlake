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
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"github.com/apache/incubator-devlake/plugins/teambition/tasks"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
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
	scopeDetails, err := dsHelper.ScopeSrv.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}
	plan, err := makePipelinePlanV200(subtaskMetas, scopeDetails, connection)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(scopeDetails, connection)
	return plan, scopes, err
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopeDetails []*srvhelper.ScopeDetail[models.TeambitionProject, srvhelper.NoScopeConfig],
	connection *models.TeambitionConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}

		scope := scopeDetail.Scope
		// construct task options for circleci
		task, err := helper.MakePipelinePlanTask(
			"teambition",
			subtaskMetas,
			nil,
			tasks.TeambitionOptions{
				ConnectionId: connection.ID,
				ProjectId:    scope.Id,
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
	scopeDetails []*srvhelper.ScopeDetail[models.TeambitionProject, srvhelper.NoScopeConfig],
	connection *models.TeambitionConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0, len(scopeDetails))

	idgen := didgen.NewDomainIdGenerator(&models.TeambitionProject{})
	for _, scopeDetail := range scopeDetails {
		scope, _ := scopeDetail.Scope, scopeDetail.ScopeConfig
		id := idgen.Generate(connection.ID, scope.Id)

		// add cicd_scope to scopes
		// if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CICD) {
		scopes = append(scopes, devops.NewCicdScope(id, scope.Name))
		// }
	}

	return scopes, nil
}
