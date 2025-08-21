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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ExtractApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "Extract Commits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: false,
	Required:         false,
	Description:      "Extract raw commit data into tool layer table bitbucket_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

type CommitsResponse struct {
	//Type string    `json:"type"`
	Sha  string    `json:"hash"`
	Date time.Time `json:"date"`
	//Author *models.BitbucketAccount
	Message string `json:"message"`
	Links   struct {
		Self     struct{ Href string }
		Html     struct{ Href string }
		Diff     struct{ Href string }
		Approve  struct{ Href string }
		Comments struct{ Href string }
	} `json:"links"`
	Parents []struct {
		Type  string
		Hash  string
		Links struct {
			Self struct{ Href string }
			Html struct{ Href string }
		}
	} `json:"parents"`
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
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

			bitbucketCommit := &models.BitbucketCommit{
				Sha:           commit.Sha,
				Message:       commit.Message,
				AuthoredDate:  commit.Date,
				Url:           commit.Links.Self.Href,
				CommittedDate: commit.Date,
			}

			//if commit.Author != nil {
			//	bitbucketCommit.AuthorName = commit.Author.User.FullDisplayName
			//	bitbucketCommit.AuthorId = commit.Author.User.AccountId
			//	results = append(results, commit.Author)
			//}

			bitbucketRepoCommit := &models.BitbucketRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       data.Options.FullName,
				CommitSha:    commit.Sha,
			}

			results = append(results, bitbucketCommit)
			results = append(results, bitbucketRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
