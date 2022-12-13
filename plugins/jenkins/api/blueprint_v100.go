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
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/plugins/jenkins/models"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

func MakePipelinePlanV100(subtaskMetas []core.SubTaskMeta, connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	var err errors.Error
	connection := new(models.JenkinsConnection)
	err = connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, err
	}
	plan, err := makePipelinePlanV100(subtaskMetas, scope, connection)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func makePipelinePlanV100(subtaskMetas []core.SubTaskMeta, scope []*core.BlueprintScopeV100, connection *models.JenkinsConnection) (core.PipelinePlan, errors.Error) {
	var err errors.Error
	plans := make(core.PipelinePlan, 0, len(scope))
	for _, scopeElem := range scope {
		// handle taskOptions and transformationRules, by dumping them to taskOptions
		transformationRules := make(map[string]interface{})
		if len(scopeElem.Transformation) > 0 {
			err = errors.Convert(json.Unmarshal(scopeElem.Transformation, &transformationRules))
			if err != nil {
				return nil, err
			}
		}

		err = GetAllJobs(context.Background(), connection, connection.Endpoint, 100, func(job *models.Job) errors.Error {
			taskOptions := make(map[string]interface{})
			err = errors.Convert(json.Unmarshal(scopeElem.Options, &taskOptions))
			if err != nil {
				return err
			}
			taskOptions["connectionId"] = connection.ID
			taskOptions["transformationRules"] = transformationRules
			taskOptions["jobFullName"] = job.FullName

			op, err := tasks.DecodeTaskOptions(taskOptions)
			if err != nil {
				return err
			}
			_, err = tasks.ValidateTaskOptions(op)
			if err != nil {
				return err
			}

			// subtasks
			subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeElem.Entities)
			if err != nil {
				return err
			}
			stage := core.PipelineStage{
				{
					Plugin:   "jenkins",
					Subtasks: subtasks,
					Options:  taskOptions,
				},
			}

			plans = append(plans, stage)

			return nil
		})
		if err != nil {
			return nil, err
		}

		// Dora
		if productionPattern, ok := transformationRules["productionPattern"]; ok && productionPattern != nil {

			doraOption := make(map[string]interface{})
			doraOption["prefix"] = "jenkins"
			doraRules := make(map[string]interface{})
			doraRules["productionPattern"] = productionPattern
			doraOption["transformationRules"] = doraRules

			stageDora := core.PipelineStage{
				{
					Plugin:   "dora",
					Subtasks: []string{"EnrichTaskEnv"},
					Options:  doraOption,
				},
			}

			plans = append(plans, stageDora)
		}
	}
	return plans, nil
}

// request all jobs
func GetAllJobs(ctx context.Context, connection *models.JenkinsConnection, baseUrl string, pageSize int, callback func(job *models.Job) errors.Error) errors.Error {
	apiClient, err := helper.NewApiClient(
		ctx,
		baseUrl,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", connection.GetEncodedToken()),
		},
		10*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return err
	}

	for i := 0; ; i += pageSize {
		var data struct {
			Jobs []json.RawMessage `json:"jobs"`
		}

		// set query
		query := url.Values{}
		treeValue := fmt.Sprintf("jobs[name,class,url,color,base,jobs,upstreamProjects[name]]{%d,%d}", i, i+pageSize)
		query.Set("tree", treeValue)

		res, err := apiClient.Get("/api/json", query, nil)
		if err != nil {
			return err
		}

		err = helper.UnmarshalResponse(res, &data)
		if err != nil {
			return err
		}

		// break with empty data
		if len(data.Jobs) == 0 {
			break
		}

		for _, rawJobs := range data.Jobs {
			job := &models.Job{}
			err1 := json.Unmarshal(rawJobs, job)
			if err1 != nil {
				return errors.Convert(err1)
			}

			if job.Jobs != nil {
				GetAllJobs(ctx, connection, baseUrl+"job/"+job.Name+"/", pageSize, callback)
			} else {
				err = callback(job)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
