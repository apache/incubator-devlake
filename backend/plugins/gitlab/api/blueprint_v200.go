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
		gitlabProject, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connectionId, scope.Id)
		if err != nil {
			return nil, err
		}

		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CODE) {
			// if we don't need to collect gitex, we need to add repo to scopes here
			scopeRepo := code.NewRepo(id, gitlabProject.Name)

			if gitlabProject.ForkedFromProjectWebUrl != "" {
				scopeRepo.ForkedFrom = gitlabProject.ForkedFromProjectWebUrl
			}
			sc = append(sc, scopeRepo)
		}

		// add cicd_scope to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopeCICD := devops.NewCicdScope(id, gitlabProject.Name)

			sc = append(sc, scopeCICD)
		}

		// add board to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_TICKET) {
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
		gitlabProject, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connection.ID, scope.Id)
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
		options["scopeConfigId"] = scopeConfig.ID
		if syncPolicy.TimeAfter != nil {
			options["timeAfter"] = syncPolicy.TimeAfter.Format(time.RFC3339)
		}

		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeConfig.Entities)
		if err != nil {
			return nil, err
		}

		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "gitlab",
			Subtasks: subtasks,
			Options:  options,
		})

		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CODE) {
			cloneUrl, err := errors.Convert01(url.Parse(gitlabProject.HttpUrlToRepo))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword("git", connection.Token)
			stage = append(stage, &plugin.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connection.ID, gitlabProject.GitlabId),
					"proxy":  connection.Proxy,
				},
			})
		}

		plans = append(plans, stage)

		// refdiff part
		if scopeConfig.Refdiff != nil {
			task := &plugin.PipelineTask{
				Plugin:  "refdiff",
				Options: scopeConfig.Refdiff,
			}
			plans = append(plans, plugin.PipelineStage{task})
		}
	}
	return plans, nil
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
