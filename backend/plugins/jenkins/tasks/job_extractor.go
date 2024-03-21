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

package tasks

import (
	"encoding/json"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

// this struct should be moved to `gitub_api_common.go`

var ExtractApiJobsMeta = plugin.SubTaskMeta{
	Name:             "extractApiJobs",
	EntryPoint:       ExtractApiJobs,
	EnabledByDefault: true,
	Description:      "Extract raw jobs data into tool layer table jenkins_jobs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiJobs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	logger := taskCtx.GetLogger()

	if data.Options.Class != WORKFLOW_MULTI_BRANCH_PROJECT {
		logger.Debug("class must be %s, got %s", WORKFLOW_MULTI_BRANCH_PROJECT, data.Options.Class)
		return nil
	}

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &models.Job{}

			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)
			job := body.ToJenkinsJob()
			job.ConnectionId = data.Options.ConnectionId
			job.ScopeConfigId = data.Options.ScopeConfigId
			job.Path = strings.Replace(body.URL, data.Connection.Endpoint, "", 1)

			results = append(results, job)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
