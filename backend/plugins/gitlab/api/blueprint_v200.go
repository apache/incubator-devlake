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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func MakePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	scope []*plugin.BlueprintScopeV200,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
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

	pp, err := makePipelinePlanV200(subtaskMetas, scope, connection, syncPolicy)
	if err != nil {
		return nil, nil, err
	}

	return pp, sc, nil
}

func makeScopeV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, 3*len(scopes))

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

		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_CODE) {
			// if we don't need to collect gitex, we need to add repo to scopes here
			scopeRepo := code.NewRepo(id, gitlabProject.Name)

			if gitlabProject.ForkedFromProjectWebUrl != "" {
				scopeRepo.ForkedFrom = gitlabProject.ForkedFromProjectWebUrl
			}
			sc = append(sc, scopeRepo)
		}

		// add cicd_scope to scopes
		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopeCICD := devops.NewCicdScope(id, gitlabProject.Name)

			sc = append(sc, scopeCICD)
		}

		// add board to scopes
		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_TICKET) {
			scopeTicket := ticket.NewBoard(id, gitlabProject.Name)

			sc = append(sc, scopeTicket)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	scopes []*plugin.BlueprintScopeV200,
	connection *models.GitlabConnection, syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	plans := make(plugin.PipelinePlan, 0, 3*len(scopes))
	for _, scope := range scopes {
		var stage plugin.PipelineStage
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
		options["transformationRuleId"] = transformationRules.ID
		if syncPolicy.TimeAfter != nil {
			options["timeAfter"] = syncPolicy.TimeAfter.Format(time.RFC3339)
		}

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scope.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "gitlab",
			Subtasks: subtasks,
			Options:  options,
		})

		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scope.Entities, plugin.DOMAIN_TYPE_CODE) {
			cloneUrl, err := errors.Convert01(url.Parse(repo.HttpUrlToRepo))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword("git", connection.Token)
			stage = append(stage, &plugin.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connection.ID, repo.GitlabId),
					"proxy":  connection.Proxy,
				},
			})
		}

		plans = append(plans, stage)

		// refdiff part
		if transformationRules.Refdiff != nil {
			task := &plugin.PipelineTask{
				Plugin:  "refdiff",
				Options: transformationRules.Refdiff,
			}
			plans = append(plans, plugin.PipelineStage{task})
		}
	}
	return plans, nil
}

// GetRepoByConnectionIdAndscopeId get tbe repo by the connectionId and the scopeId
func GetRepoByConnectionIdAndscopeId(connectionId uint64, scopeId string) (*models.GitlabProject, errors.Error) {
	gitlabId, e := strconv.Atoi(scopeId)
	if e != nil {
		return nil, errors.Default.Wrap(e, fmt.Sprintf("scopeId %s is not integer", scopeId))
	}
	repo := &models.GitlabProject{}
	db := basicRes.GetDal()
	err := db.First(repo, dal.Where("connection_id = ? AND gitlab_id = ?", connectionId, gitlabId))
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
		db := basicRes.GetDal()
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

func GetApiProject(
	op *tasks.GitlabOptions,
	apiClient aha.ApiClientAbstract,
) (*models.GitlabApiProject, errors.Error) {
	repoRes := &models.GitlabApiProject{}
	res, err := apiClient.Get(fmt.Sprintf("projects/%d", op.ProjectId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting repo detail from %s", res.Request.URL.String()))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}
	err = errors.Convert(json.Unmarshal(body, repoRes))
	if err != nil {
		return nil, err
	}
	return repoRes, nil
}
