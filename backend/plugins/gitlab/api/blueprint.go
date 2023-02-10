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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func MakePipelinePlan(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	var err errors.Error
	connection := new(models.GitlabConnection)
	err = connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClientFromConnection(
		context.TODO(),
		basicRes,
		connection,
	)
	if err != nil {
		return nil, err
	}

	plan, err := makePipelinePlan(subtaskMetas, scope, apiClient, connection)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func makePipelinePlan(
	subtaskMetas []plugin.SubTaskMeta,
	scope []*plugin.BlueprintScopeV100,
	apiClient aha.ApiClientAbstract,
	connection *models.GitlabConnection,
) (plugin.PipelinePlan, errors.Error) {
	var err errors.Error
	var repo *tasks.GitlabApiProject
	plan := make(plugin.PipelinePlan, len(scope))
	for i, scopeElem := range scope {
		// handle taskOptions and transformationRules, by dumping them to taskOptions
		transformationRules := make(map[string]interface{})
		if len(scopeElem.Transformation) > 0 {
			err = errors.Convert(json.Unmarshal(scopeElem.Transformation, &transformationRules))
			if err != nil {
				return nil, err
			}
		}
		// construct task options for github
		options := make(map[string]interface{})
		err = errors.Convert(json.Unmarshal(scopeElem.Options, &options))
		if err != nil {
			return nil, err
		}
		options["connectionId"] = connection.ID
		options["transformationRules"] = transformationRules
		// make sure task options is valid
		op, err := tasks.DecodeAndValidateTaskOptions(options)
		if err != nil {
			return nil, err
		}

		memorizedGetApiRepo := func() (*tasks.GitlabApiProject, errors.Error) {
			if repo == nil {
				repo, err = getApiRepo(op, apiClient)
			}
			return repo, err
		}

		// refdiff
		if refdiffRules, ok := transformationRules["refdiff"]; ok && refdiffRules != nil {
			// add a new task to next stage
			j := i + 1
			if j == len(plan) {
				plan = append(plan, nil)
			}
			repo, err = memorizedGetApiRepo()
			if err != nil {
				return nil, err
			}
			ops := refdiffRules.(map[string]interface{})
			ops["repoId"] = didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(connection.ID, repo.GitlabId)
			plan[j] = plugin.PipelineStage{
				{
					Plugin:  "refdiff",
					Options: ops,
				},
			}
			// remove it from github transformationRules
			delete(transformationRules, "refdiff")
		}

		// construct subtasks
		subtasks, err := api.MakePipelinePlanSubtasks(subtaskMetas, scopeElem.Entities)
		if err != nil {
			return nil, err
		}
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "gitlab",
			Subtasks: subtasks,
			Options:  options,
		})
		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scopeElem.Entities, plugin.DOMAIN_TYPE_CODE) {
			// here is the tricky part, we have to obtain the repo id beforehand
			repo, err = memorizedGetApiRepo()
			if err != nil {
				return nil, err
			}
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
		// This is just to add a dora subtask, then we can add another two subtasks at the end of plans
		// The only purpose is to adapt old blueprints
		// DEPRECATED, will be removed in v0.17
		// dora
		if productionPattern, ok := transformationRules["productionPattern"]; ok && productionPattern != nil {
			j := i + 1
			// add a new task to next stage
			if j == len(plan) {
				plan = append(plan, nil)
			}
			if plan[j] != nil {
				j++
			}
			if j == len(plan) {
				plan = append(plan, nil)
			}
			repo, err = memorizedGetApiRepo()
			if err != nil {
				return nil, err
			}
			plan[j] = plugin.PipelineStage{
				{
					Plugin: "dora",
					// This is just to add a dora subtask, then we can add another two subtasks at the end of plans
					// The only purpose is to adapt old blueprints
					Subtasks: []string{"EnrichTaskEnv"},
					Options:  map[string]interface{}{},
				},
			}
		}
		plan[i] = stage
		repo = nil
	}
	return plan, nil
}

func getApiRepo(
	op *tasks.GitlabOptions,
	apiClient aha.ApiClientAbstract,
) (*tasks.GitlabApiProject, errors.Error) {
	apiRepo := &tasks.GitlabApiProject{}
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
	err = errors.Convert(json.Unmarshal(body, apiRepo))
	if err != nil {
		return nil, err
	}
	return apiRepo, nil
}
