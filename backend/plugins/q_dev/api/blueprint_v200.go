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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// load connection and scope from the db
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
	scopes, err := makeScopesV200(scopeDetails, connection)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopeDetails []*srvhelper.ScopeDetail[models.QDevS3Slice, srvhelper.NoScopeConfig],
	connection *models.QDevConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		s3Slice := scopeDetail.Scope
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}

		// construct task options for q_dev
		op := &tasks.QDevOptions{
			ConnectionId: s3Slice.ConnectionId,
			S3Prefix:     s3Slice.Prefix,
			ScopeId:      s3Slice.Id,
		}

		// Pass empty entities array to enable all subtasks
		task, err := helper.MakePipelinePlanTask("q_dev", subtaskMetas, []string{}, op)
		if err != nil {
			return nil, err
		}
		stage = append(stage, task)
		plan[i] = stage
	}
	return plan, nil
}

func makeScopesV200(
	scopeDetails []*srvhelper.ScopeDetail[models.QDevS3Slice, srvhelper.NoScopeConfig],
	connection *models.QDevConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	// For Q Developer metrics, we don't need to create domain layer scopes
	// The data is collected and stored directly in the tool layer
	return scopes, nil
}
