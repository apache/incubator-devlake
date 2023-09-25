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

	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func MakePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	scope []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	var err errors.Error
	connection := new(models.BambooConnection)
	err1 := connectionHelper.FirstById(connection, connectionId)
	if err1 != nil {
		return nil, nil, errors.Default.Wrap(err1, fmt.Sprintf("error on get connection by id[%d]", connectionId))
	}

	sc, err := makeScopeV200(connectionId, scope)
	if err != nil {
		return nil, nil, err
	}

	pp, err := makePipelinePlanV200(subtaskMetas, scope, connection)
	if err != nil {
		return nil, nil, err
	}

	return pp, sc, nil
}

func makeScopeV200(connectionId uint64, scopes []*coreModels.BlueprintScope) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, len(scopes))

	for _, scope := range scopes {
		id := didgen.NewDomainIdGenerator(&models.BambooPlan{}).Generate(connectionId, scope.ScopeId)

		// get project from db
		project, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connectionId, scope.ScopeId)
		if err != nil {
			return nil, err
		}

		// add cicd_scope to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopeCICD := devops.NewCicdScope(id, project.Name)

			sc = append(sc, scopeCICD)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopes []*coreModels.BlueprintScope,
	connection *models.BambooConnection,
) (coreModels.PipelinePlan, errors.Error) {
	plans := make(coreModels.PipelinePlan, 0, len(scopes))
	for _, scope := range scopes {
		var stage coreModels.PipelineStage
		var err errors.Error
		// get project
		_, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connection.ID, scope.ScopeId)
		if err != nil {
			return nil, err
		}

		// bamboo main part
		options := make(map[string]interface{})
		options["connectionId"] = connection.ID
		options["planKey"] = scope.ScopeId
		options["scopeConfigId"] = scopeConfig.ID

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeConfig.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &coreModels.PipelineTask{
			Plugin:   "bamboo",
			Subtasks: subtasks,
			Options:  options,
		})

		plans = append(plans, stage)
	}
	return plans, nil
}
