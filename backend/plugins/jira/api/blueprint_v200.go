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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// load connection, scope and scopeConfig from the db
	connection, err := dsHelper.ConnSrv.FindByPk(connectionId)
	if err != nil {
		return nil, nil, err
	}
	scopeDetails, err := dsHelper.ScopeSrv.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}

	// needed for the connection to populate its access tokens
	// if AppKey authentication method is selected
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
	scopeDetails []*srvhelper.ScopeDetail[models.JiraBoard, models.JiraScopeConfig],
	connection *models.JiraConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}

		scope, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		// construct task options for Jira
		task, err := helper.MakePipelinePlanTask(
			"jira",
			subtaskMetas,
			scopeConfig.Entities,
			JiraTaskOptions{
				ConnectionId: scope.ConnectionId,
				BoardId:      scope.BoardId,
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
	scopeDetails []*srvhelper.ScopeDetail[models.JiraBoard, models.JiraScopeConfig],
	connection *models.JiraConnection,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, scopeDetail := range scopeDetails {
		jiraBoard, scopeConfig := scopeDetail.Scope, scopeDetail.ScopeConfig
		// add board to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_TICKET) {
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(jiraBoard.ConnectionId, jiraBoard.BoardId),
				},
				Name: jiraBoard.Name,
				Url:  jiraBoard.Self,
				Type: jiraBoard.Type,
			}
			scopes = append(scopes, domainBoard)
		}
	}
	return scopes, nil
}
