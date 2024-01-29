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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ExtractApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw commit data into tool layer table bitbucket_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

type ApiCommitResponse struct {
	BitbucketId        string          `json:"id"`
	DisplayId          string          `json:"displayId"`
	Author             ApiUserResponse `json:"author"`
	Message            string          `json:"message"`
	AuthorTimestamp    int64           `json:"authorTimestamp"`
	CommitterTimestamp int64           `json:"committerTimestamp"`
	Parents            []struct {
		BitbucketID string `json:"id"`
		DisplayID   string `json:"displayId"`
	} `json:"parents"`
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMITS_TABLE)
	repoId := data.Options.FullName

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			commit := &ApiCommitResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, commit))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 4)

			parentCommitShasArr := make([]string, 0, len(commit.Parents))
			for _, parent := range commit.Parents {
				parentCommitShasArr = append(parentCommitShasArr, parent.BitbucketID)
			}
			parentCommitShas := strings.Join(parentCommitShasArr, ",")

			bitbucketCommit := &models.BitbucketServerCommit{
				ConnectionId:     data.Options.ConnectionId,
				RepoId:           repoId,
				CommitSha:        commit.BitbucketId,
				ParentCommitShas: parentCommitShas,
				AuthorName:       commit.Author.Name,
				AuthorEmail:      commit.Author.EmailAddress,
				Message:          commit.Message,
				AuthoredDate:     time.UnixMilli(commit.AuthorTimestamp),
				CommittedDate:    time.UnixMilli(commit.CommitterTimestamp),
			}

			bitbucketRepoCommit := &models.BitbucketServerRepoCommit{
				ConnectionId:     data.Options.ConnectionId,
				RepoId:           data.Options.FullName,
				CommitSha:        commit.BitbucketId,
				ParentCommitShas: parentCommitShas,
			}

			bitbucketUser, err := convertUser(&commit.Author, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, bitbucketUser, bitbucketCommit, bitbucketRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
