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

	"github.com/apache/incubator-devlake/plugins/bamboo/models"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func MakePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	scope []*plugin.BlueprintScopeV200,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
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

	pp, err := makePipelinePlanV200(subtaskMetas, scope, connection, syncPolicy)
	if err != nil {
		return nil, nil, err
	}

	return pp, sc, nil
}

func makeScopeV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, len(scopes))

	for _, scope := range scopes {
		id := didgen.NewDomainIdGenerator(&models.BambooProject{}).Generate(connectionId, scope.Id)

		// get project from db
		project, err := GetProjectByConnectionIdAndscopeId(connectionId, scope.Id)
		if err != nil {
			return nil, err
		}

		scopeConfig, err := GetScopeConfigByproject(project)
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
	scopes []*plugin.BlueprintScopeV200,
	connection *models.BambooConnection, syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	plans := make(plugin.PipelinePlan, 0, len(scopes))
	for _, scope := range scopes {
		var stage plugin.PipelineStage
		var err errors.Error
		// get project
		project, err := GetProjectByConnectionIdAndscopeId(connection.ID, scope.Id)
		if err != nil {
			return nil, err
		}

		scopeConfig, err := GetScopeConfigByproject(project)
		if err != nil {
			return nil, err
		}

		// bamboo main part
		options := make(map[string]interface{})
		options["connectionId"] = connection.ID
		options["projectKey"] = scope.Id
		options["scopeConfigId"] = scopeConfig.ID

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeConfig.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "bamboo",
			Subtasks: subtasks,
			Options:  options,
		})

		plans = append(plans, stage)
	}
	return plans, nil
}

// GetProjectByConnectionIdAndscopeId get tbe project by the connectionId and the scopeId
func GetProjectByConnectionIdAndscopeId(connectionId uint64, scopeId string) (*models.BambooProject, errors.Error) {
	key := scopeId
	project := &models.BambooProject{}
	db := basicRes.GetDal()
	err := db.First(project, dal.Where("connection_id = ? AND project_key = ?", connectionId, key))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find project by connection [%d] scope [%s]", connectionId, scopeId))
		}
		return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find project by connection [%d] scope [%s]", connectionId, scopeId))
	}

	return project, nil
}

// GetScopeConfigByproject get the BambooScopeConfig by project
func GetScopeConfigByproject(project *models.BambooProject) (*models.BambooScopeConfig, errors.Error) {
	scopeConfig := &models.BambooScopeConfig{}
	scopeConfigId := project.ScopeConfigId
	if scopeConfigId != 0 {
		db := basicRes.GetDal()
		err := db.First(scopeConfig, dal.Where("id = ?", scopeConfigId))
		if err != nil {
			if db.IsErrorNotFound(err) {
				return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find ScopeConfig by ScopeConfig [%d]", scopeConfigId))
			}
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find ScopeConfig by ScopeConfig [%d]", scopeConfigId))
		}
	} else {
		scopeConfig.ID = 0
	}

	return scopeConfig, nil
}
