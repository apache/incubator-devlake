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

type ApiPrCommitResponse struct {
	BitbucketId        string          `json:"id"`
	DisplayId          string          `json:"displayId"`
	Author             ApiUserResponse `json:"author"`
	Message            string          `json:"message"`
	AuthorTimestamp    int64           `json:"authorTimestamp"`
	CommitterTimestamp int64           `json:"committerTimestamp"`
}

func ExtractApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMITS_TABLE)
	repoId := data.Options.FullName
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiPullRequestCommit := &ApiPrCommitResponse{}
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

			bitbucketPullRequestCommit := &models.BitbucketServerPrCommit{
				ConnectionId:       data.Options.ConnectionId,
				RepoId:             repoId,
				PullRequestId:      pull.BitbucketId,
				CommitSha:          apiPullRequestCommit.BitbucketId,
				CommitAuthorName:   apiPullRequestCommit.Author.DisplayName,
				CommitAuthorEmail:  apiPullRequestCommit.Author.EmailAddress,
				CommitAuthoredDate: time.UnixMilli(apiPullRequestCommit.AuthorTimestamp),
			}

			return []interface{}{bitbucketPullRequestCommit}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
