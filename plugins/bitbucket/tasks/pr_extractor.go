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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"
)

var ExtractApiPullRequestsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequests",
	EntryPoint:       ExtractApiPullRequests,
	EnabledByDefault: true,
	Required:         true,
	Description:      "Extract raw PullRequests data into tool layer table bitbucket_pull_requests",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

type BitbucketApiPullRequest struct {
	BitbucketId  int    `json:"id"`
	CommentCount int    `json:"comment_count"`
	TaskCount    int    `json:"task_count"`
	Type         string `json:"type"`
	State        string `json:"state"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	MergeCommit  *struct {
		Type  string
		Hash  string `json:"hash"`
		Links *struct {
			Self struct{ Href string }
			Html struct{ Href string }
		}
	} `json:"merge_commit"`
	Links *struct {
		Self struct{ Href string }
		Html struct{ Href string }
	}
	ClosedBy           *BitbucketAccountResponse `json:"closed_by"`
	Author             *BitbucketAccountResponse `json:"author"`
	BitbucketCreatedAt time.Time                 `json:"created_on"`
	BitbucketUpdatedAt time.Time                 `json:"updated_on"`
	BaseRef            *struct {
		Branch struct {
			Name string
		} `json:"branch"`
		Commit struct {
			Type string
			Hash string
		} `json:"commit"`
		Repo *BitbucketApiRepo `json:"repository"`
	} `json:"destination"`
	HeadRef *struct {
		Branch struct {
			Name string
		} `json:"branch"`
		Commit struct {
			Type string
			Hash string
		} `json:"commit"`
		Repo *BitbucketApiRepo `json:"repository"`
	} `json:"source"`
	Reviewers    []BitbucketAccountResponse `json:"reviewers"`
	Participants []BitbucketAccountResponse `json:"participants"`
}

func ExtractApiPullRequests(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*BitbucketTaskData)
	var err errors.Error
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			rawL := &BitbucketApiPullRequest{}
			err := errors.Convert(json.Unmarshal(row.Data, rawL))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)
			if rawL.BitbucketId == 0 {
				return nil, nil
			}

			bitbucketPr, err := convertBitbucketPullRequest(rawL, data.Options.ConnectionId, data.Repo.BitbucketId)
			if err != nil {
				return nil, err
			}
			if rawL.Author != nil {
				bitbucketUser, err := convertAccount(rawL.Author, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, bitbucketUser)
				bitbucketPr.AuthorName = bitbucketUser.DisplayName
				bitbucketPr.AuthorId = bitbucketUser.AccountId
			}
			if rawL.MergeCommit != nil {
				bitbucketPr.MergeCommitSha = rawL.MergeCommit.Hash
			}
			results = append(results, bitbucketPr)

			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
func convertBitbucketPullRequest(pull *BitbucketApiPullRequest, connId uint64, repoId string) (*models.BitbucketPullRequest, errors.Error) {
	bitbucketPull := &models.BitbucketPullRequest{
		ConnectionId:       connId,
		BitbucketId:        pull.BitbucketId,
		RepoId:             repoId,
		BaseRepoId:         pull.BaseRef.Repo.FullName,
		HeadRepoId:         pull.HeadRef.Repo.FullName,
		State:              pull.State,
		Title:              pull.Title,
		Description:        pull.Description,
		Url:                pull.Links.Html.Href,
		Type:               pull.Type,
		CommentCount:       pull.CommentCount,
		BitbucketCreatedAt: pull.BitbucketCreatedAt,
		BitbucketUpdatedAt: pull.BitbucketUpdatedAt,
		BaseRef:            pull.BaseRef.Branch.Name,
		BaseCommitSha:      pull.BaseRef.Commit.Hash,
		HeadRef:            pull.HeadRef.Branch.Name,
		HeadCommitSha:      pull.HeadRef.Commit.Hash,
	}
	return bitbucketPull, nil
}
