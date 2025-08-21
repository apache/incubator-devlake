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
	"regexp"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ExtractApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequests",
	EntryPoint:       ExtractApiPullRequests,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw PullRequests data into tool layer table bitbucket_pull_requests",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type ApiPrResponse struct {
	BitbucketId int `json:"id"`
	// Type        string `json:"type"`
	State       string `json:"state"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MergeCommit *struct {
		Hash string `json:"hash"`
		// date only return when fields defined
		Date *common.Iso8601Time `json:"date"`
	} `json:"merge_commit"`
	Links *struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Author *struct {
		User *ApiUserResponse `json:"user"`
	} `json:"author"`
	BitbucketCreatedAt int64  `json:"createdDate"`
	BitbucketUpdatedAt int64  `json:"updatedDate"`
	BitbucketClosedAt  *int64 `json:"closedDate"`
	BaseRef            *struct {
		Branch string                         `json:"displayId"`
		Commit string                         `json:"latestCommit"`
		Repo   *models.BitbucketServerApiRepo `json:"repository"`
	} `json:"toRef"`
	HeadRef *struct {
		Branch string                         `json:"displayId"`
		Commit string                         `json:"latestCommit"`
		Repo   *models.BitbucketServerApiRepo `json:"repository"`
	} `json:"fromRef"`
	Properties *struct {
		ResolvedTaskCount int `json:"resolvedTaskCount"`
		CommentCount      int `json:"commentCount"`
		OpenTaskCount     int `json:"openTaskCount"`
	} `json:"properties"`
}

func ExtractApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*BitbucketServerTaskData)
	config := data.Options.BitbucketServerScopeConfig
	var prTypeRegex *regexp.Regexp
	var prComponentRegex *regexp.Regexp
	var err errors.Error

	if config != nil {
		if len(config.PrType) > 0 {
			prTypeRegex, err = errors.Convert01(regexp.Compile(config.PrType))
			if err != nil {
				return errors.Default.Wrap(err, "regexp Compile prType failed")
			}
		}

		if len(config.PrComponent) > 0 {
			prComponentRegex, err = errors.Convert01(regexp.Compile(config.PrComponent))
			if err != nil {
				return errors.Default.Wrap(err, "regexp Compile prComponent failed")
			}
		}
	}

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			rawL := &ApiPrResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, rawL))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)
			if rawL.BitbucketId == 0 {
				return nil, nil
			}

			bitbucketPr, err := convertBitbucketPullRequest(rawL, data.Options.ConnectionId, data.Options.FullName)
			if err != nil {
				return nil, err
			}
			if rawL.Author != nil {
				bitbucketUser, err := convertUser(rawL.Author.User, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, bitbucketUser)
				bitbucketPr.AuthorName = bitbucketUser.DisplayName
				bitbucketPr.AuthorId = bitbucketUser.BitbucketId
			}
			if rawL.MergeCommit != nil {
				bitbucketPr.MergeCommitSha = rawL.MergeCommit.Hash
				bitbucketPr.MergedAt = rawL.MergeCommit.Date.ToNullableTime()
			} else if rawL.State == code.MERGED && rawL.BitbucketClosedAt != nil {
				mergedAt := time.UnixMilli(*rawL.BitbucketClosedAt)
				bitbucketPr.MergedAt = &mergedAt
			}

			if prTypeRegex != nil {
				prTypes := prTypeRegex.FindStringSubmatch(rawL.Title)
				if len(prTypes) > 0 {
					bitbucketPr.Type = prTypes[0]
				}
			}
			if prComponentRegex != nil {
				prComponents := prComponentRegex.FindStringSubmatch(rawL.Description)
				if len(prComponents) > 0 {
					bitbucketPr.Component = prComponents[0]
				}
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
func convertBitbucketPullRequest(pull *ApiPrResponse, connId uint64, repoId string) (*models.BitbucketServerPullRequest, errors.Error) {
	bitbucketPull := &models.BitbucketServerPullRequest{
		ConnectionId: connId,
		BitbucketId:  pull.BitbucketId,
		Number:       pull.BitbucketId,
		RepoId:       repoId,
		BaseRepoId:   repoId, // using pull.BaseRef.Repo.Slug is not suitable because we can have same repo name in multiple projects
		State:        pull.State,
		Title:        pull.Title,
		Description:  pull.Description,
		// Type:               pull.Type,
		BitbucketServerCreatedAt: time.UnixMilli(pull.BitbucketCreatedAt),
		BitbucketServerUpdatedAt: time.UnixMilli(pull.BitbucketUpdatedAt),
	}
	if pull.Links != nil && len(pull.Links.Self) > 0 {
		bitbucketPull.Url = pull.Links.Self[0].Href
	}

	if pull.BaseRef != nil {
		bitbucketPull.BaseRef = pull.BaseRef.Branch
		bitbucketPull.BaseCommitSha = pull.BaseRef.Commit
	}
	if pull.HeadRef != nil {
		if pull.HeadRef.Repo != nil {
			bitbucketPull.HeadRepoId = pull.HeadRef.Repo.Slug
		}
		bitbucketPull.HeadRef = pull.HeadRef.Branch
		bitbucketPull.HeadCommitSha = pull.HeadRef.Commit
	}
	if pull.Properties != nil {
		bitbucketPull.CommentCount = pull.Properties.CommentCount
	}
	if pull.BitbucketClosedAt != nil {
		closedAt := time.UnixMilli(*pull.BitbucketClosedAt)
		bitbucketPull.ClosedAt = &closedAt
	}

	return bitbucketPull, nil
}
