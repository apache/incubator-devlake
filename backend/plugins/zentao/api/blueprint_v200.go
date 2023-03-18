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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, bpScopes []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	connectionHelper := helper.NewConnectionHelper(basicRes, validator.New())
	// get the connection info for url
	connection := &models.ZentaoConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, scopes, err := makePipelinePlanV200(subtaskMetas, plan, bpScopes, connection, syncPolicy)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan plugin.PipelinePlan,
	bpScopes []*plugin.BlueprintScopeV200,
	connection *models.ZentaoConnection,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	var err errors.Error
	domainScopes := make([]plugin.Scope, 0)
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		// construct task options
		op := &tasks.ZentaoOptions{
			ConnectionId: connection.ID,
		}

		scopeType := strings.Split(bpScope.Id, `/`)[0]
		scopeId := strings.Split(bpScope.Id, `/`)[1]
		if scopeType == `project` {
			scope := &models.ZentaoProject{}
			// get repo from db
			err = basicRes.GetDal().First(scope, dal.Where(`connection_id = ? AND id = ?`, connection.ID, scopeId))
			if err != nil {
				return nil, nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find zentao project %s", bpScope.Id))
			}
			op.ProjectId = scope.Id

			if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_TICKET) {
				scopeTicket := &ticket.Board{
					DomainEntity: domainlayer.DomainEntity{
						Id: didgen.NewDomainIdGenerator(&models.ZentaoProject{}).Generate(connection.ID, scope.Id),
					},
					Name: scope.Name,
					Type: scope.Type,
				}
				domainScopes = append(domainScopes, scopeTicket)
			}
		} else {
			scope := &models.ZentaoProduct{}
			// get repo from db
			err = basicRes.GetDal().First(scope, dal.Where(`connection_id = ? AND id = ?`, connection.ID, scopeId))
			if err != nil {
				return nil, nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find zentao product %s", bpScope.Id))
			}
			op.ProductId = scope.Id

			if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_TICKET) {
				scopeTicket := &ticket.Board{
					DomainEntity: domainlayer.DomainEntity{
						Id: didgen.NewDomainIdGenerator(&models.ZentaoProduct{}).Generate(connection.ID, scope.Id),
					},
					Name: scope.Name,
					Type: scope.Type,
				}
				domainScopes = append(domainScopes, scopeTicket)
			}
		}

		if syncPolicy.TimeAfter != nil {
			op.TimeAfter = syncPolicy.TimeAfter.Format(time.RFC3339)
		}
		options, err := tasks.EncodeTaskOptions(op)
		if err != nil {
			return nil, nil, err
		}

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
		if err != nil {
			return nil, nil, err
		}
		stage = append(stage, &plugin.PipelineTask{
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
