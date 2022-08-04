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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
)

var EnrichPipelinesMeta = core.SubTaskMeta{
	Name:             "enrichPipelines",
	EntryPoint:       EnrichPipelines,
	EnabledByDefault: true,
	Description:      "Create tool layer table github_pipelines from github_runs",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichPipelines(taskCtx core.SubTaskContext) (err error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	entity := &githubModels.GithubPipeline{}
	cursor, err := db.Cursor(
		dal.Select("head_sha, head_branch, status, conclusion, github_created_at, github_updated_at, run_attempt, run_started_at"),
		dal.From(&githubModels.GithubRun{}),
		dal.Orderby("head_sha"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		var item githubModels.GithubRun
		err = db.Fetch(cursor, &item)
		if err != nil {
			return err
		}

		if item.HeadSha != entity.Commit {
			entity.ConnectionId = data.Options.ConnectionId
			entity.Commit = item.HeadSha
			entity.Branch = item.HeadBranch // to do
			entity.StartedDate = item.GithubCreatedAt
			entity.FinishedDate = item.GithubUpdatedAt
			entity.Status = item.Status
			if entity.Status == "completed" {
				entity.Duration = float64(item.GithubUpdatedAt.Sub(*item.GithubCreatedAt).Seconds())
			}
			entity.Results = item.Conclusion
			// TODO
			entity.Type = "CI/CD"
		} else {
			if item.GithubCreatedAt.Before(*entity.StartedDate) {
				entity.StartedDate = item.GithubCreatedAt
			}
			if item.GithubUpdatedAt.After(*entity.FinishedDate) {
				entity.FinishedDate = item.GithubCreatedAt
			}
			if item.Status != "completed" {
				entity.Status = item.Status
			} else {
				entity.Duration = float64(item.GithubUpdatedAt.Sub(*item.GithubCreatedAt).Seconds())
			}
			if item.Conclusion != "success" {
				entity.Results = item.Conclusion
			}

		}
		err := db.CreateOrUpdate(entity)
		if err != nil {
			return err
		}
	}

	return err

}
