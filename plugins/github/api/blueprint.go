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
	"github.com/apache/incubator-devlake/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
)

func MakePipelinePlan(subtaskMetas []core.SubTaskMeta, connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	var err errors.Error
	connection := new(models.GithubConnection)
	err = connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, err
	}
	token := strings.Split(connection.Token, ",")[0]

	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		10*time.Second,
		connection.Proxy,
		basicRes,
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

func makePipelinePlan(subtaskMetas []core.SubTaskMeta, scope []*core.BlueprintScopeV100, apiClient helper.ApiClientGetter, connection *models.GithubConnection) (core.PipelinePlan, errors.Error) {
	var err errors.Error
	var repo *tasks.GithubApiRepo
	plan := make(core.PipelinePlan, len(scope))
	for i, scopeElem := range scope {
		// handle taskOptions and transformationRules, by dumping them to taskOptions
		transformationRules := make(map[string]interface{})
		if len(scopeElem.Transformation) > 0 {
			err = errors.Convert(json.Unmarshal(scopeElem.Transformation, &transformationRules))
			if err != nil {
				return nil, err
			}
		}
		// refdiff
		if refdiffRules, ok := transformationRules["refdiff"]; ok && refdiffRules != nil {
			// add a new task to next stage
			j := i + 1
			if j == len(plan) {
				plan = append(plan, nil)
			}
			plan[j] = core.PipelineStage{
				{
					Plugin:  "refdiff",
					Options: refdiffRules.(map[string]interface{}),
				},
			}
			// remove it from github transformationRules
			delete(transformationRules, "refdiff")
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
		memorizedGetApiRepo := func() (*tasks.GithubApiRepo, errors.Error) {
			if repo == nil {
				repo, err = getApiRepo(op, apiClient)
			}
			return repo, err
		}

		stage := plan[i]
		if stage == nil {
			stage = core.PipelineStage{}
		}

		// construct github(graphql) task
		if connection.EnableGraphql {
			// FIXME this need fix when 2 plugins merged
			plugin, err := core.GetPlugin(`github_graphql`)
			if err != nil {
				return nil, err
			}
			if pluginGq, ok := plugin.(core.PluginTask); ok {
				subtasks, err := helper.MakePipelinePlanSubtasks(pluginGq.SubTaskMetas(), scopeElem.Entities)
				if err != nil {
					return nil, err
				}
				stage = append(stage, &core.PipelineTask{
					Plugin:   "github_graphql",
					Subtasks: subtasks,
					Options:  options,
				})
			} else {
				return nil, errors.BadInput.New("plugin github_graphql does not support SubTaskMetas")
			}
		} else {
			subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeElem.Entities)
			if err != nil {
				return nil, err
			}
			stage = append(stage, &core.PipelineTask{
				Plugin:   "github",
				Subtasks: subtasks,
				Options:  options,
			})
		}
		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scopeElem.Entities, core.DOMAIN_TYPE_CODE) {
			// here is the tricky part, we have to obtain the repo id beforehand
			token := strings.Split(connection.Token, ",")[0]
			repo, err = memorizedGetApiRepo()
			if err != nil {
				return nil, err
			}
			cloneUrl, err := errors.Convert01(url.Parse(repo.CloneUrl))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword("git", token)
			stage = append(stage, &core.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, repo.GithubId),
					"proxy":  connection.Proxy,
				},
			})
		}
		// dora
		if productionPattern, ok := transformationRules["productionPattern"]; ok && productionPattern != nil {
			j := i + 1
			if j == len(plan) {
				plan = append(plan, nil)
			}
			// add a new task to next stage
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

			plan[j] = core.PipelineStage{
				{
					Plugin:   "dora",
					Subtasks: []string{"EnrichTaskEnv"},
					Options: map[string]interface{}{
						"repoId": didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, repo.GithubId),
						"transformationRules": map[string]interface{}{
							"productionPattern": productionPattern,
						},
					},
				},
			}
			// remove it from github transformationRules
			delete(transformationRules, "productionPattern")
		}
		plan[i] = stage
		repo = nil
	}
	return plan, nil
}

func getApiRepo(op *tasks.GithubOptions, apiClient helper.ApiClientGetter) (*tasks.GithubApiRepo, errors.Error) {
	apiRepo := &tasks.GithubApiRepo{}
	res, err := apiClient.Get(fmt.Sprintf("repos/%s/%s", op.Owner, op.Repo), nil, nil)
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
