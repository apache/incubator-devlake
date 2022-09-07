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
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table GiteaCommit,GiteaAccount and GiteaRepoCommit",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_CROSS},
}

type GiteaCommit struct {
	Author struct {
		Date  helper.Iso8601Time `json:"date"`
		Email string             `json:"email"`
		Name  string             `json:"name"`
	}
	Committer struct {
		Date  helper.Iso8601Time `json:"date"`
		Email string             `json:"email"`
		Name  string             `json:"name"`
	}
	Message string `json:"message"`
}

type GiteaApiCommitResponse struct {
	Author    *models.GiteaAccount `json:"author"`
	Commit    GiteaCommit          `json:"commit"`
	Committer *models.GiteaAccount `json:"committer"`
	HtmlUrl   string               `json:"html_url"`
	Sha       string               `json:"sha"`
	Url       string               `json:"url"`
}

func ExtractApiCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			results := make([]interface{}, 0, 4)

			commit := &GiteaApiCommitResponse{}

			err := json.Unmarshal(row.Data, commit)

			if err != nil {
				return nil, err
			}

			if commit.Sha == "" {
				return nil, nil
			}

			giteaCommit, err := ConvertCommit(commit)

			if err != nil {
				return nil, err
			}

			if commit.Author != nil {
				giteaCommit.AuthorId = commit.Author.Id
				results = append(results, commit.Author)
			}
			if commit.Committer != nil {
				giteaCommit.CommitterId = commit.Committer.Id
				results = append(results, commit.Committer)
			}

			giteaRepoCommit := &models.GiteaRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       data.Repo.GiteaId,
				CommitSha:    commit.Sha,
			}
			results = append(results, giteaCommit)
			results = append(results, giteaRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// ConvertCommit Convert the API response to our DB model instance
func ConvertCommit(commit *GiteaApiCommitResponse) (*models.GiteaCommit, error) {
	giteaCommit := &models.GiteaCommit{
		Sha:            commit.Sha,
		Message:        commit.Commit.Message,
		AuthorName:     commit.Commit.Author.Name,
		AuthorEmail:    commit.Commit.Author.Email,
		AuthoredDate:   commit.Commit.Author.Date.ToTime(),
		CommitterName:  commit.Commit.Author.Name,
		CommitterEmail: commit.Commit.Author.Email,
		CommittedDate:  commit.Commit.Author.Date.ToTime(),
		WebUrl:         commit.Url,
	}
	return giteaCommit, nil
}
