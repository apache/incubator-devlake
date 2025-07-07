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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
	"github.com/apache/incubator-devlake/plugins/testmo/tasks"
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
	scopeDetails []*srvhelper.ScopeDetail[models.TestmoProject, models.TestmoScopeConfig],
	connection *models.TestmoConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}

		scope, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		// construct task options for testmo
		task, err := api.MakePipelinePlanTask(
			"testmo",
			subtaskMetas,
			scopeConfig.Entities,
			tasks.TestmoOptions{
				ConnectionId: connection.ID,
				ProjectId:    scope.Id,
				ScopeConfig:  scopeConfig,
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
	scopeDetails []*srvhelper.ScopeDetail[models.TestmoProject, models.TestmoScopeConfig],
	connection *models.TestmoConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0, len(scopeDetails))

	for _, scopeDetail := range scopeDetails {
		scope := scopeDetail.Scope
		// Return the TestmoProject itself as it implements plugin.Scope
		scopes = append(scopes, scope)
	}

	return scopes, nil
}
