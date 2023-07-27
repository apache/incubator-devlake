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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiCommitsMeta)
}

var ExtractApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: false,
	Description:      "Extract raw commit data into tool layer table github_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	DependencyTables: []string{RAW_COMMIT_TABLE},
	ProductTables: []string{
		models.GithubCommit{}.TableName(),
		models.GithubRepoCommit{}.TableName()},
}

type CommitsResponse struct {
	Sha       string `json:"sha"`
	Commit    Commit
	Url       string
	Author    *models.GithubAccount
	Committer *models.GithubAccount
}

type Commit struct {
	Author struct {
		Name  string
		Email string
		Date  api.Iso8601Time
	}
	Committer struct {
		Name  string
		Email string
		Date  api.Iso8601Time
	}
	Message string
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraCommits by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			/*
				Table store raw data
			*/
			Table: RAW_COMMIT_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			commit := &CommitsResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, commit))
			if err != nil {
				return nil, err
			}
			if commit.Sha == "" {
				return nil, nil
			}

			results := make([]interface{}, 0, 4)

			githubCommit := &models.GithubCommit{
				Sha:            commit.Sha,
				Message:        commit.Commit.Message,
				AuthorName:     commit.Commit.Author.Name,
				AuthorEmail:    commit.Commit.Author.Email,
				AuthoredDate:   commit.Commit.Author.Date.ToTime(),
				CommitterName:  commit.Commit.Committer.Name,
				CommitterEmail: commit.Commit.Committer.Email,
				CommittedDate:  commit.Commit.Committer.Date.ToTime(),
				Url:            commit.Url,
			}
			if commit.Author != nil {
				githubCommit.AuthorId = commit.Author.Id
				results = append(results, commit.Author)
			}
			if commit.Committer != nil {
				githubCommit.CommitterId = commit.Committer.Id
				results = append(results, commit.Committer)
			}

			githubRepoCommit := &models.GithubRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       data.Options.GithubId,
				CommitSha:    commit.Sha,
			}

			results = append(results, githubCommit)
			results = append(results, githubRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
