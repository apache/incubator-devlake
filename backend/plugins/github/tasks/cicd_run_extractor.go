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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractRunsMeta)
}

var ExtractRunsMeta = plugin.SubTaskMeta{
	Name:             "Extract Workflow Runs",
	EntryPoint:       ExtractRuns,
	EnabledByDefault: true,
	Description:      "Extract raw run data into tool layer table github_runs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{RAW_RUN_TABLE},
	ProductTables:    []string{models.GithubRun{}.TableName()},
}

func ExtractRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_RUN_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			githubRun := &models.GithubRun{}
			err := errors.Convert(json.Unmarshal(row.Data, githubRun))
			if err != nil {
				return nil, err
			}
			githubRun.RepoId = repoId
			githubRun.ConnectionId = data.Options.ConnectionId
			githubRun.Type = data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, githubRun.Name)
			githubRun.Environment = data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, githubRun.Name, githubRun.HeadBranch)
			return []interface{}{githubRun}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
