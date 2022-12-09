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
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/utils"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func MakePipelinePlanV200(subtaskMetas []core.SubTaskMeta, connectionId uint64, scope []*core.BlueprintScopeV200) (core.PipelinePlan, []core.Scope, errors.Error) {
	var err errors.Error
	connection := new(models.GitlabConnection)
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

func makeScopeV200(connectionId uint64, scopes []*core.BlueprintScopeV200) ([]core.Scope, errors.Error) {
	sc := make([]core.Scope, 0, 3*len(scopes))

	for _, scope := range scopes {
		intScopeId, err1 := strconv.Atoi(scope.Id)
		if err1 != nil {
			return nil, errors.Default.Wrap(err1, fmt.Sprintf("Failed to strconv.Atoi for scope.Id [%s]", scope.Id))
		}
		id := didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connectionId, intScopeId)

		// get repo from db
		gitlabProject, err := GetRepoByConnectionIdAndscopeId(connectionId, scope.Id)
		if err != nil {
			return nil, err
		}

		if utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_CODE) ||
			utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_CROSS) {
			// if we don't need to collect gitex, we need to add repo to scopes here
			scopeRepo := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: id,
				},
				Name: gitlabProject.Name,
			}
			if gitlabProject.ForkedFromProjectWebUrl != "" {
				scopeRepo.ForkedFrom = gitlabProject.ForkedFromProjectWebUrl
			}
			sc = append(sc, scopeRepo)
		}

		// add cicd_scope to scopes
		if utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_CICD) {
			scopeCICD := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: id,
				},
				Name: gitlabProject.Name,
			}
			sc = append(sc, scopeCICD)
		}

		// add board to scopes
		if utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_TICKET) {
			scopeTicket := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: id,
				},
				Name: gitlabProject.Name,
			}
			sc = append(sc, scopeTicket)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(subtaskMetas []core.SubTaskMeta, scopes []*core.BlueprintScopeV200, connection *models.GitlabConnection) (core.PipelinePlan, errors.Error) {
	plans := make(core.PipelinePlan, 0, 3*len(scopes))
	for _, scope := range scopes {
		var stage core.PipelineStage
		var err errors.Error
		// get repo
		repo, err := GetRepoByConnectionIdAndscopeId(connection.ID, scope.Id)
		if err != nil {
			return nil, err
		}

		// get transformationRuleId
		transformationRules, err := GetTransformationRuleByRepo(repo)
		if err != nil {
			return nil, err
		}

		// get int scopeId
		intScopeId, err1 := strconv.Atoi(scope.Id)
		if err != nil {
			return nil, errors.Default.Wrap(err1, fmt.Sprintf("Failed to strconv.Atoi for scope.Id [%s]", scope.Id))
		}

		// gitlab main part
		options := make(map[string]interface{})
		options["connectionId"] = connection.ID
		options["projectId"] = intScopeId
		options["entities"] = scope.Entities

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scope.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &core.PipelineTask{
			Plugin:   "gitlab",
			Subtasks: subtasks,
			Options:  options,
		})

		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scope.Entities, core.DOMAIN_TYPE_CODE) {
			cloneUrl, err := errors.Convert01(url.Parse(repo.HttpUrlToRepo))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword("git", connection.Token)
			stage = append(stage, &core.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connection.ID, repo.GitlabId),
					"proxy":  connection.Proxy,
				},
			})
		}

		// refdiff part
		if transformationRules.Refdiff != nil {
			task := &core.PipelineTask{
				Plugin:  "refdiff",
				Options: transformationRules.Refdiff,
			}
			stage = append(stage, task)
		}

		plans = append(plans, stage)

	}
	return plans, nil
}

// GetRepoByConnectionIdAndscopeId get tbe repo by the connectionId and the scopeId
func GetRepoByConnectionIdAndscopeId(connectionId uint64, scopeId string) (*models.GitlabProject, errors.Error) {
	repo := &models.GitlabProject{}
	db := BasicRes.GetDal()
	err := db.First(repo, dal.Where("connection_id = ? AND gitlab_id = ?", connectionId, scopeId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find repo by connection [%d] scope [%s]", connectionId, scopeId))
		}
		return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find repo by connection [%d] scope [%s]", connectionId, scopeId))
	}

	return repo, nil
}

// GetTransformationRuleByRepo get the GetTransformationRule by Repo
func GetTransformationRuleByRepo(repo *models.GitlabProject) (*models.GitlabTransformationRule, errors.Error) {
	transformationRules := &models.GitlabTransformationRule{}
	transformationRuleId := repo.TransformationRuleId
	if transformationRuleId != 0 {
		db := BasicRes.GetDal()
		err := db.First(transformationRules, dal.Where("id = ?", transformationRuleId))
		if err != nil {
			if db.IsErrorNotFound(err) {
				return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find transformationRules by transformationRuleId [%d]", transformationRuleId))
			}
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find transformationRules by transformationRuleId [%d]", transformationRuleId))
		}
	} else {
		transformationRules.ID = 0
	}

	return transformationRules, nil
}
