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
	"fmt"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiRepoMeta = core.SubTaskMeta{
	Name:        "extractApiRepo",
	EntryPoint:  ExtractApiRepositories,
	Required:    true,
	Description: "Extract raw Repositories data into tool layer table github_repos",
	DomainTypes: core.DOMAIN_TYPES,
}

type ApiRepoResponse GithubApiRepo

type GithubApiRepo struct {
	Name        string `json:"name"`
	GithubId    int    `json:"id"`
	HTMLUrl     string `json:"html_url"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Owner       *GithubUserResponse
	Parent      *GithubApiRepo      `json:"parent"`
	CreatedAt   helper.Iso8601Time  `json:"created_at"`
	UpdatedAt   *helper.Iso8601Time `json:"updated_at"`
	CloneUrl    string              `json:"clone_url"`
}

func ExtractApiRepositories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_REPOSITORIES_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiRepoResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			if body.GithubId == 0 {
				return nil, fmt.Errorf("repo %s/%s not found", data.Options.Owner, data.Options.Repo)
			}
			results := make([]interface{}, 0, 1)
			githubRepository := &models.GithubRepo{
				ConnectionId: data.Options.ConnectionId,
				GithubId:     body.GithubId,
				Name:         body.Name,
				HTMLUrl:      body.HTMLUrl,
				Description:  body.Description,
				OwnerId:      body.Owner.Id,
				OwnerLogin:   body.Owner.Login,
				Language:     body.Language,
				CreatedDate:  body.CreatedAt.ToTime(),
				UpdatedDate:  helper.Iso8601TimeToTime(body.UpdatedAt),
			}
			data.Repo = githubRepository

			if body.Parent != nil {
				githubRepository.ParentGithubId = body.Parent.GithubId
				githubRepository.ParentHTMLUrl = body.Parent.HTMLUrl
			}
			results = append(results, githubRepository)

			githubUser, err := convertUser(body.Owner, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, githubUser)

			parentTaskContext := taskCtx.TaskContext()
			if parentTaskContext != nil {
				parentTaskContext.GetData().(*GithubTaskData).Repo = githubRepository
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
