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
	goerror "errors"
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []core.SubTaskMeta, connectionId uint64, bpScopes []*core.BlueprintScopeV200) (core.PipelinePlan, []core.Scope, errors.Error) {
	// get the connection info for url
	connection := &models.JiraConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(core.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connection)
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
	subtaskMetas []core.SubTaskMeta,
	plan core.PipelinePlan,
	bpScopes []*core.BlueprintScopeV200,
	connection *models.JiraConnection,
) (core.PipelinePlan, errors.Error) {
	db := basicRes.GetDal()
	var err errors.Error
	var stage core.PipelineStage
	for i, bpScope := range bpScopes {
		jiraBoard := &models.JiraBoard{}
		// get repo from db
		err = db.First(jiraBoard, dal.Where(`connection_id = ? and board_id = ?`, connection.ID, bpScope.Id))
		if err != nil && goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find board%s", bpScope.Id))
		}
		transformationRule := &models.JiraTransformationRule{}
		// get transformation rules from db
		err = db.First(transformationRule, dal.Where(`id = ?`, jiraBoard.TransformationRuleId))
		if err != nil && !goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// construct task options for Jira
		var options map[string]interface{}
		err = errors.Convert(mapstructure.Decode(jiraBoard, &options))
		if err != nil {
			return nil, err
		}
		// make sure task options is valid
		_, err = tasks.DecodeAndValidateTaskOptions(options)
		if err != nil {
			return nil, err
		}

		var transformationRuleMap map[string]interface{}
		err = errors.Convert(mapstructure.Decode(transformationRule, &transformationRuleMap))
		if err != nil {
			return nil, err
		}
		options["transformationRules"] = transformationRuleMap
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
		if err != nil {
			return nil, err
		}
		stage = append(stage, &core.PipelineTask{
			Plugin:   "jira",
			Subtasks: subtasks,
			Options:  options,
		})
		plan[i] = stage
	}

	return plan, nil
}

func makeScopesV200(bpScopes []*core.BlueprintScopeV200, connection *models.JiraConnection) ([]core.Scope, errors.Error) {
	scopes := make([]core.Scope, 0)
	for _, bpScope := range bpScopes {
		jiraBoard := &models.JiraBoard{}
		// get repo from db
		err := basicRes.GetDal().First(jiraBoard,
			dal.Where(`connection_id = ? and board_id = ?`,
				connection.ID, bpScope.Id))
		if err != nil && goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find board%d", jiraBoard.BoardId))
		}
		transformationRule := &models.JiraTransformationRule{}
		// get transformation rules from db
		err = basicRes.GetDal().First(transformationRule, dal.Where(`id = ?`, jiraBoard.TransformationRuleId))
		if err != nil && !goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// add board to scopes
		if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_TICKET) {
			jiraBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(jiraBoard.ConnectionId, jiraBoard.BoardId),
				},
				Name: jiraBoard.Name,
			}
			scopes = append(scopes, jiraBoard)
		}
	}
	return scopes, nil
}
