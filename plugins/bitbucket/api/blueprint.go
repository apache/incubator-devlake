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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
)

func MakePipelinePlan(subtaskMetas []core.SubTaskMeta, connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	var err errors.Error
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
		if refdiffRules, ok := transformationRules["refdiff"]; ok {
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
			// remove it from bitbucket transformationRules
			delete(transformationRules, "refdiff")
		}
		// construct task options for bitbucket
		options := make(map[string]interface{})
		err = errors.Convert(json.Unmarshal(scopeElem.Options, &options))
		if err != nil {
			return nil, err
		}
		options["connectionId"] = connectionId
		options["transformationRules"] = transformationRules
		// make sure task options is valid
		op, err := tasks.DecodeAndValidateTaskOptions(options)
		if err != nil {
			return nil, err
		}
		// construct subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeElem.Entities)
		if err != nil {
			return nil, err
		}
		stage := plan[i]
		if stage == nil {
			stage = core.PipelineStage{}
		}
		stage = append(stage, &core.PipelineTask{
			Plugin:   "bitbucket",
			Subtasks: subtasks,
			Options:  options,
		})
		// collect git data by gitextractor if CODE was requested
		if utils.StringsContains(scopeElem.Entities, core.DOMAIN_TYPE_CODE) {
			// here is the tricky part, we have to obtain the repo id beforehand
			connection := new(models.BitbucketConnection)
			err = connectionHelper.FirstById(connection, connectionId)
			if err != nil {
				return nil, err
			}
			token := strings.Split(connection.GetEncodedToken(), ",")[0]
			apiClient, err := helper.NewApiClient(
				context.TODO(),
				connection.Endpoint,
				map[string]string{
					"Authorization": fmt.Sprintf("Basic %s", token),
				},
				10*time.Second,
				connection.Proxy,
				basicRes,
			)
			if err != nil {
				return nil, err
			}
			res, err := apiClient.Get(fmt.Sprintf("repositories/%s/%s", op.Owner, op.Repo), nil, nil)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				return nil, errors.Default.New(fmt.Sprintf(
					"unexpected status code when requesting repo detail %d %s",
					res.StatusCode, res.Request.URL.String(),
				))
			}
			body, err := errors.Convert01(io.ReadAll(res.Body))
			if err != nil {
				return nil, err
			}
			apiRepo := new(tasks.BitbucketApiRepo)
			err = errors.Convert(json.Unmarshal(body, apiRepo))
			if err != nil {
				return nil, err
			}
			cloneUrl, err := errors.Convert01(url.Parse(apiRepo.Links.Clone[0].Href))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword("git", token)
			stage = append(stage, &core.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connectionId, apiRepo.BitbucketId),
				},
			})
		}
		plan[i] = stage
	}
	return plan, nil
}
