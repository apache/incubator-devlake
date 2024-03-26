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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
)

var _ plugin.SubTaskEntryPoint = ExtractAccounts

var ExtractJobsMeta = plugin.SubTaskMeta{
	Name:             "Extract Jobs",
	EntryPoint:       ExtractJobs,
	EnabledByDefault: true,
	Description:      "Extract raw run data into tool layer table github_jobs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractJobs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: githubTasks.GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_GRAPHQL_JOBS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			apiJob := &GraphqlQueryCheckRunWrapper{}
			err := errors.Convert(json.Unmarshal(row.Data, apiJob))
			if err != nil {
				return nil, err
			}

			nodes := apiJob.Node
			results := make([]interface{}, 0, 1)
			for _, node := range nodes {
				for _, checkRun := range node.CheckSuite.CheckRuns.Nodes {

					paramsBytes, err := json.Marshal(checkRun.Steps.Nodes)
					if err != nil {
						taskCtx.GetLogger().Error(err, `Marshal checkRun.Steps.Nodes fail and ignore`)
					}
					githubJob := &models.GithubJob{
						ConnectionId: data.Options.ConnectionId,
						RunID:        node.CheckSuite.WorkflowRun.DatabaseId,
						RepoId:       data.Options.GithubId,
						ID:           checkRun.DatabaseId,
						NodeID:       checkRun.Id,
						HTMLURL:      checkRun.DetailsUrl,
						Status:       strings.ToUpper(checkRun.Status),
						Conclusion:   strings.ToUpper(checkRun.Conclusion),
						StartedAt:    checkRun.StartedAt,
						CompletedAt:  checkRun.CompletedAt,
						Name:         checkRun.Name,
						Steps:        paramsBytes,
						Type:         data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, checkRun.Name),
						Environment:  data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, checkRun.Name),
						// these columns can not fill by graphql
						//HeadSha:       ``,  // use _tool_github_runs
						//RunURL:        ``,
						//CheckRunURL:   ``,
						//Labels:        ``, // not in use
						//RunnerID:      ``, // not in use
						//RunnerName:    ``, // not in use
						//RunnerGroupID: ``, // not in use
					}
					results = append(results, githubJob)
				}
			}
			return results, nil

		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
