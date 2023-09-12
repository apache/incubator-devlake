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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
	syncPolicy *coreModels.SyncPolicy,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// get the connection info for url
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(coreModels.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connection, syncPolicy)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(bpScopes, connection)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan coreModels.PipelinePlan,
	bpScopes []*coreModels.BlueprintScope,
	connection *models.PagerDutyConnection,
	syncPolicy *coreModels.SyncPolicy,
) (coreModels.PipelinePlan, errors.Error) {
	for i, bpScope := range bpScopes {
		// get board and scope config from db
		service, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connection.ID, bpScope.ScopeId)
		if err != nil {
			return nil, err
		}
		// construct task options for pagerduty
		op := &tasks.PagerDutyOptions{
			ConnectionId: service.ConnectionId,
			ServiceId:    service.Id,
			ServiceName:  service.Name,
		}
		if syncPolicy.TimeAfter != nil {
			op.TimeAfter = syncPolicy.TimeAfter.Format(time.RFC3339)
		}
		var options map[string]any
		options, err = tasks.EncodeTaskOptions(op)
		if err != nil {
			return nil, err
		}
		var subtasks []string
		subtasks, err = api.MakePipelinePlanSubtasks(subtaskMetas, scopeConfig.Entities)
		if err != nil {
			return nil, err
		}
		stage := []*coreModels.PipelineTask{
			{
				Plugin:   "pagerduty",
				Subtasks: subtasks,
				Options:  options,
			},
		}
		plan[i] = stage
	}
	return plan, nil
}

func makeScopesV200(bpScopes []*coreModels.BlueprintScope, connection *models.PagerDutyConnection) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, bpScope := range bpScopes {
		// get board and scope config from db
		service, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connection.ID, bpScope.ScopeId)
		if err != nil {
			return nil, err
		}
		// add board to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_TICKET) {
			scopeTicket := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.Service{}).Generate(connection.ID, service.Id),
				},
				Name: service.Name,
			}
			scopes = append(scopes, scopeTicket)
		}
	}
	return scopes, nil
}
