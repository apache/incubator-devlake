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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// get the connection info for url
	connection := &models.ZentaoConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(coreModels.PipelinePlan, len(bpScopes))
	plan, scopes, err := makePipelinePlanV200(subtaskMetas, plan, bpScopes, connection)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan coreModels.PipelinePlan,
	bpScopes []*coreModels.BlueprintScope,
	connection *models.ZentaoConnection,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	domainScopes := make([]plugin.Scope, 0)
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}
		// construct task options
		op := &tasks.ZentaoOptions{
			ConnectionId: connection.ID,
		}

		var entities []string

		project, scopeConfig, err := projectScopeHelper.DbHelper().GetScopeAndConfig(connection.ID, bpScope.ScopeId)
		if err != nil {
			return nil, nil, err
		}
		op.ProjectId = project.Id
		entities = scopeConfig.Entities

		if utils.StringsContains(entities, plugin.DOMAIN_TYPE_TICKET) {
			scopeTicket := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.ZentaoProject{}).Generate(connection.ID, project.Id),
				},
				Name: project.Name,
				Type: project.Type,
			}
			domainScopes = append(domainScopes, scopeTicket)
		}
		/*} else {
			product, scopeConfig, err := productScopeHelper.DbHelper().GetScopeAndConfig(connection.ID, scopeId)
			if err != nil {
				return nil, nil, err
			}
			op.ProductId = product.Id
			entities = scopeConfig.Entities

			if utils.StringsContains(entities, plugin.DOMAIN_TYPE_TICKET) {
				scopeTicket := &ticket.Board{
					DomainEntity: domainlayer.DomainEntity{
						Id: didgen.NewDomainIdGenerator(&models.ZentaoProduct{}).Generate(connection.ID, product.Id),
					},
					Name: product.Name,
					Type: product.Type,
				}
				domainScopes = append(domainScopes, scopeTicket)
			}
		}*/

		options, err := tasks.EncodeTaskOptions(op)
		if err != nil {
			return nil, nil, err
		}

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, entities)
		if err != nil {
			return nil, nil, err
		}
		stage = append(stage, &coreModels.PipelineTask{
			Plugin:   "zentao",
			Subtasks: subtasks,
			Options:  options,
		})
		if err != nil {
			return nil, nil, err
		}

		plan[i] = stage
	}
	return plan, domainScopes, nil
}
