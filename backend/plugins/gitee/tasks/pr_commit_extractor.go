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
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"strings"
)

var ExtractApiPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestCommits",
	EntryPoint:       ExtractApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestCommits data into tool layer table gitee_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type PrCommitsResponse struct {
	Sha    string `json:"sha"`
	Commit PullRequestCommit
	Url    string
	Author struct {
		Id    int
		Login string
		Name  string
	}
	Committer struct {
		Id    int
		Login string
		Name  string
	}
}

type PullRequestCommit struct {
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
	Message      string
	CommentCount int `json:"comment_count"`
}

func ExtractApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMIT_TABLE)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			apiPullRequestCommit := &PrCommitsResponse{}
			if strings.HasPrefix(string(row.Data), "{\"message\": \"Not Found\"") {
				return nil, nil
			}
			err := errors.Convert(json.Unmarshal(row.Data, apiPullRequestCommit))
			if err != nil {
				return nil, err
			}
			pull := &SimplePr{}
			err = errors.Convert(json.Unmarshal(row.Input, pull))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 2)

			giteeCommit, err := convertPullRequestCommit(apiPullRequestCommit)
			if err != nil {
				return nil, err
			}
			results = append(results, giteeCommit)

			giteePullRequestCommit := &models.GiteePullRequestCommit{
				ConnectionId:       data.Options.ConnectionId,
				CommitSha:          apiPullRequestCommit.Sha,
				PullRequestId:      pull.GiteeId,
				CommitAuthorEmail:  apiPullRequestCommit.Commit.Author.Email,
				CommitAuthorName:   apiPullRequestCommit.Commit.Author.Name,
				CommitAuthoredDate: apiPullRequestCommit.Commit.Author.Date.ToTime(),
			}
			if err != nil {
				return nil, err
			}
			results = append(results, giteePullRequestCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestCommit(prCommit *PrCommitsResponse) (*models.GiteeCommit, errors.Error) {
	giteeCommit := &models.GiteeCommit{
		Sha:            prCommit.Sha,
		Message:        prCommit.Commit.Message,
		AuthorId:       prCommit.Author.Id,
		AuthorName:     prCommit.Commit.Author.Name,
		AuthorEmail:    prCommit.Commit.Author.Email,
		AuthoredDate:   prCommit.Commit.Author.Date.ToTime(),
		CommitterName:  prCommit.Commit.Committer.Name,
		CommitterEmail: prCommit.Commit.Committer.Email,
		CommittedDate:  prCommit.Commit.Committer.Date.ToTime(),
		WebUrl:         prCommit.Url,
	}
	return giteeCommit, nil
}
