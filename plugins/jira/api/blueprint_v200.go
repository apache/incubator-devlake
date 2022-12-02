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
	db := basicRes.GetDal()
	// get the connection info for url
	connection := &models.JiraConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(core.PipelinePlan, 0, len(bpScopes))
	scopes := make([]core.Scope, 0, len(bpScopes))
	for i, bpScope := range bpScopes {
		var jiraBoard *models.JiraBoard
		// get repo from db
		err = db.First(jiraBoard, dal.Where(`id = ?`, bpScope.Id))
		if err != nil {
			return nil, nil, err
		}
		var transformationRule *models.JiraTransformationRule
		// get transformation rules from db
		err = db.First(transformationRule, dal.Where(`id = ?`, jiraBoard.TransformationRuleId))
		if err != nil && goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		var scope []core.Scope
		// make pipeline for each bpScope
		plan[i], scope, err = makeDataSourcePipelinePlanV200(subtaskMetas, bpScope, jiraBoard, transformationRule)
		if err != nil {
			return nil, nil, err
		}
		if len(scope) > 0 {
			scopes = append(scopes, scope...)
		}

	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []core.SubTaskMeta,
	bpScope *core.BlueprintScopeV200,
	jiraBoard *models.JiraBoard,
	transformationRule *models.JiraTransformationRule,
) (core.PipelineStage, []core.Scope, errors.Error) {
	var err errors.Error
	var stage core.PipelineStage
	scopes := make([]core.Scope, 0)

	// construct task options for jenkins
	var options map[string]interface{}
	err = errors.Convert(mapstructure.Decode(jiraBoard, &options))
	if err != nil {
		return nil, nil, err
	}
	// make sure task options is valid
	_, err = tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, nil, err
	}

	var transformationRuleMap map[string]interface{}
	err = errors.Convert(mapstructure.Decode(transformationRule, &transformationRuleMap))
	if err != nil {
		return nil, nil, err
	}
	options["transformationRules"] = transformationRuleMap
	subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
	if err != nil {
		return nil, nil, err
	}
	stage = append(stage, &core.PipelineTask{
		Plugin:   "jira",
		Subtasks: subtasks,
		Options:  options,
	})

	// add cicd_scope to scopes
	if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_TICKET) {
		scopeCICD := &ticket.Board{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(jiraBoard.ConnectionId, jiraBoard.BoardId),
			},
			Name: jiraBoard.Name,
		}
		scopes = append(scopes, scopeCICD)
	}

	return stage, scopes, nil
}
