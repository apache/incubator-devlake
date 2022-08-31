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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractJobsMeta = core.SubTaskMeta{
	Name:             "extractJobs",
	EntryPoint:       ExtractJobs,
	EnabledByDefault: true,
	Description:      "Extract raw run data into tool layer table github_jobs",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractJobs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_JOB_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			githubJob := &models.GithubJob{}
			err := json.Unmarshal(row.Data, githubJob)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1)
			githubJobResult := &models.GithubJob{
				ConnectionId:  data.Options.ConnectionId,
				RepoId:        repoId,
				ID:            githubJob.ID,
				RunID:         githubJob.RunID,
				RunURL:        githubJob.RunURL,
				NodeID:        githubJob.NodeID,
				HeadSha:       githubJob.HeadSha,
				URL:           githubJob.URL,
				HTMLURL:       githubJob.HTMLURL,
				Status:        githubJob.Status,
				Conclusion:    githubJob.Conclusion,
				StartedAt:     githubJob.StartedAt,
				CompletedAt:   githubJob.CompletedAt,
				Name:          githubJob.Name,
				Steps:         githubJob.Steps,
				CheckRunURL:   githubJob.CheckRunURL,
				Labels:        githubJob.Labels,
				RunnerID:      githubJob.RunID,
				RunnerName:    githubJob.RunnerName,
				RunnerGroupID: githubJob.RunnerGroupID,
				// TODO
				//Type: "CI/CD",
			}
			results = append(results, githubJobResult)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
