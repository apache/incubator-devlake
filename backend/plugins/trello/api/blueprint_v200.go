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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"time"

	"github.com/apache/incubator-devlake/plugins/trello/models"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func MakePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, scope []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	scopes, err := makeScopeV200(connectionId, scope)
	if err != nil {
		return nil, nil, err
	}

	plan := make(plugin.PipelinePlan, len(scope))
	plan, err = makePipelinePlanV200(subtaskMetas, plan, scope, connectionId, syncPolicy)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeScopeV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, len(scopes))

	for _, scope := range scopes {
		trelloBoard := &models.TrelloBoard{}
		// get board from db
		err := basicRes.GetDal().First(trelloBoard,
			dal.Where(`connection_id = ? and board_id = ?`,
				connectionId, scope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find board %s", scope.Id))
		}
		// add board to scopes
		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_TICKET) {
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.TrelloConnection{}).Generate(trelloBoard.ConnectionId, trelloBoard.BoardId),
				},
				Name: trelloBoard.Name,
			}
			sc = append(sc, domainBoard)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, plan plugin.PipelinePlan, scopes []*plugin.BlueprintScopeV200, connectionId uint64, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, errors.Error) {
	for i, scope := range scopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}

		// construct task options for trello
		options := make(map[string]interface{})
		options["connectionId"] = connectionId
		options["scopeId"] = scope.Id
		if syncPolicy.TimeAfter != nil {
			options["timeAfter"] = syncPolicy.TimeAfter.Format(time.RFC3339)
		}

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scope.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "trello",
			Subtasks: subtasks,
			Options:  options,
		})

		plan[i] = stage
	}
	return plan, nil
}
