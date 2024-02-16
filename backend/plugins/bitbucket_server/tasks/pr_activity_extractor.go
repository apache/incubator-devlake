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
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ExtractApiPrActivitiesMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestsActivities",
	EntryPoint:       ExtractApiPullRequestActivities,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw pull requests activity data into tool layer table(s)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type ApiPrActivityResponse struct {
	BitbucketId int   `json:"id"`
	CreatedOn   int64 `json:"createdDate"`

	User *ApiUserResponse `json:"user"`

	Action string `json:"action"`

	CommentAction *string `json:"commentAction"`
	Comment       *struct {
		BitbucketId int    `json:"id"`
		Text        string `json:"text"`
		CreatedAt   int64  `json:"createdDate"`
		UpdatedAt   *int64 `json:"updatedDate"`
		Severity    string `json:"severity"`
		State       string `json:"state"`
	} `json:"comment"`

	Commit *struct {
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
	} `json:"commit"`
}

func ExtractApiPullRequestActivities(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_ACTIVITIES_TABLE)
	db := taskCtx.GetDal()

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			prActivity := &ApiPrActivityResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, prActivity))
			if err != nil {
				return nil, err
			}

			prId := &struct {
				BitbucketId int
			}{}
			err = errors.Convert(json.Unmarshal(row.Input, prId))
			if err != nil {
				return nil, err
			}

			if prActivity.Action == "COMMENTED" && prActivity.Comment != nil {
				toolprComment, err := convertPullRequestComment(prActivity)
				if err != nil {
					return nil, err
				}
				toolprComment.ConnectionId = data.Options.ConnectionId
				toolprComment.RepoId = data.Options.FullName
				toolprComment.PullRequestId = prId.BitbucketId

				results := make([]interface{}, 0, 2)

				if prActivity.User != nil {
					bitbucketUser, err := convertUser(prActivity.User, data.Options.ConnectionId)
					if err != nil {
						return nil, err
					}
					toolprComment.AuthorId = bitbucketUser.BitbucketId
					toolprComment.AuthorName = bitbucketUser.DisplayName
					results = append(results, bitbucketUser)
				}
				results = append(results, toolprComment)

				return results, nil
			} else if prActivity.Action == "MERGED" && prActivity.Commit != nil {
				params := BitbucketServerApiParams{}
				err = errors.Convert(json.Unmarshal([]byte(row.Params), &params))
				if err != nil {
					return nil, err
				}

				err := db.UpdateColumn(
					&models.BitbucketServerPullRequest{
						ConnectionId: params.ConnectionId,
						RepoId:       params.FullName,
						BitbucketId:  prId.BitbucketId,
					},
					"merge_commit_sha",
					prActivity.Commit.BitbucketId,
				)
				if err != nil {
					return nil, err
				}
			}

			return []interface{}{}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestComment(prActivity *ApiPrActivityResponse) (*models.BitbucketServerPrComment, errors.Error) {
	bitbucketPrComment := &models.BitbucketServerPrComment{
		BitbucketId: prActivity.BitbucketId,
		AuthorName:  prActivity.User.DisplayName,
		CreatedAt:   time.UnixMilli(prActivity.CreatedOn),
		Body:        prActivity.Comment.Text,
	}
	if prActivity.Comment.UpdatedAt != nil {
		updatedAt := time.UnixMilli(*prActivity.Comment.UpdatedAt)
		bitbucketPrComment.UpdatedAt = &updatedAt
	}
	return bitbucketPrComment, nil
}
