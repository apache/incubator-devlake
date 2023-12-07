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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ExtractApiPrCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestCommits",
	EntryPoint:       ExtractApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestCommits data into tool layer table bitbucket_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type ApiPrCommitsResponse struct {
	BitbucketId        string                `json:"id"`
	DisplayId          string                `json:"displayId"`
	Author             BitbucketUserResponse `json:"author"`
	Message            string                `json:"message"`
	AuthorTimestamp    int64                 `json:"authorTimestamp"`
	CommitterTimestamp int64                 `json:"committerTimestamp"`
}

func ExtractApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMITS_TABLE)
	repoId := data.Options.FullName
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiPullRequestCommit := &ApiPrCommitsResponse{}
			if strings.HasPrefix(string(row.Data), "Not Found") {
				return nil, nil
			}
			err := errors.Convert(json.Unmarshal(row.Data, apiPullRequestCommit))
			if err != nil {
				return nil, err
			}
			pull := &BitbucketServerInput{}
			err = errors.Convert(json.Unmarshal(row.Input, pull))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 3)
			bitbucketRepoCommit := &models.BitbucketServerRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       repoId,
				CommitSha:    apiPullRequestCommit.BitbucketId,
			}
			results = append(results, bitbucketRepoCommit)

			bitbucketCommit, err := convertPullRequestCommit(apiPullRequestCommit, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			bitbucketCommit.ConnectionId = data.Options.ConnectionId
			bitbucketCommit.RepoId = repoId
			results = append(results, bitbucketCommit)

			bitbucketPullRequestCommit := &models.BitbucketServerPrCommit{
				ConnectionId:       data.Options.ConnectionId,
				RepoId:             repoId,
				PullRequestId:      pull.BitbucketId,
				CommitSha:          apiPullRequestCommit.BitbucketId,
				CommitAuthorName:   apiPullRequestCommit.Author.DisplayName,
				CommitAuthorEmail:  apiPullRequestCommit.Author.EmailAddress,
				CommitAuthoredDate: time.UnixMilli(apiPullRequestCommit.AuthorTimestamp),
			}
			if err != nil {
				return nil, err
			}
			results = append(results, bitbucketPullRequestCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestCommit(prCommit *ApiPrCommitsResponse, connId uint64) (*models.BitbucketServerCommit, errors.Error) {
	bitbucketCommit := &models.BitbucketServerCommit{
		CommitSha:     prCommit.BitbucketId,
		Message:       prCommit.Message,
		AuthorId:      prCommit.Author.BitbucketId,
		AuthorName:    prCommit.Author.Name,
		AuthorEmail:   prCommit.Author.EmailAddress,
		AuthoredDate:  time.UnixMilli(prCommit.AuthorTimestamp),
		CommittedDate: time.UnixMilli(prCommit.CommitterTimestamp),
	}
	return bitbucketCommit, nil
}
