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
	"github.com/apache/devlake/core/errors"
	coreModels "github.com/apache/devlake/core/models"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/helpers/srvhelper"
	"github.com/apache/devlake/plugins/kiro/models"
	"github.com/apache/devlake/plugins/kiro/tasks"
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

	plan, err := makeDataSourcePipelinePlanV200(subtaskMetas, scopeDetails, connection)
	if err != nil {
		return nil, nil, err
	}

	// No domain layer scopes: this plugin writes only to _tool_kiro_* tables.
	// Cross-tool AI modelling in the domain layer is separate work.
	return plan, []plugin.Scope{}, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopeDetails []*srvhelper.ScopeDetail[models.KiroS3Slice, srvhelper.NoScopeConfig],
	connection *models.KiroConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		slice := scopeDetail.Scope

		op := &tasks.KiroOptions{
			ConnectionId: slice.ConnectionId,
			ScopeId:      slice.Id,
			AccountId:    slice.AccountId,
			Year:         slice.Year,
			Month:        slice.Month,
		}

		// An empty entity list enables every subtask; the three streams are
		// always collected together because they describe the same activity.
		task, err := helper.MakePipelinePlanTask("kiro", subtaskMetas, []string{}, op)
		if err != nil {
			return nil, err
		}

		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}
		plan[i] = append(stage, task)
	}
	return plan, nil
}
