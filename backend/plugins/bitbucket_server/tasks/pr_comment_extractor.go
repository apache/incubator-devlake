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

var ExtractApiPrCommentsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestsComments",
	EntryPoint:       ExtractApiPullRequestsComments,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw pull requests comments data into tool layer table BitbucketPrComments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type BitbucketPrCommentsResponse struct {
	BitbucketId int   `json:"id"`
	CreatedOn   int64 `json:"createdDate"`

	User *BitbucketUserResponse `json:"user"`

	Action        string  `json:"action"`
	CommentAction *string `json:"commentAction"`
	Comment       *struct {
		BitbucketId int    `json:"id"`
		Text        string `json:"text"`
		CreatedAt   int64  `json:"createdDate"`
		UpdatedAt   *int64 `json:"updatedDate"`
		Severity    string `json:"severity"`
		State       string `json:"state"`
	} `json:"comment"`
}

func ExtractApiPullRequestsComments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_ACTIVITIES_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			prComment := &BitbucketPrCommentsResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, prComment))
			if err != nil {
				return nil, err
			} else if prComment.Action != "COMMENTED" {
				return []interface{}{}, nil
			}

			input := &struct {
				BitbucketId int
			}{}
			err = errors.Convert(json.Unmarshal(row.Input, input))
			if err != nil {
				return nil, err
			}

			toolprComment, err := convertPullRequestComment(prComment)
			if err != nil {
				return nil, err
			}
			toolprComment.ConnectionId = data.Options.ConnectionId
			toolprComment.RepoId = data.Options.FullName
			toolprComment.PullRequestId = input.BitbucketId

			results := make([]interface{}, 0, 2)

			if prComment.User != nil {
				bitbucketUser, err := convertUser(prComment.User, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				toolprComment.AuthorId = bitbucketUser.BitbucketId
				toolprComment.AuthorName = bitbucketUser.DisplayName
				results = append(results, bitbucketUser)
			}
			results = append(results, toolprComment)

			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestComment(prComment *BitbucketPrCommentsResponse) (*models.BitbucketServerPrComment, errors.Error) {
	bitbucketPrComment := &models.BitbucketServerPrComment{
		BitbucketId: prComment.BitbucketId,
		AuthorName:  prComment.User.DisplayName,
		CreatedAt:   time.UnixMilli(prComment.CreatedOn),
		Body:        prComment.Comment.Text,
	}
	if prComment.Comment.UpdatedAt != nil {
		updatedAt := time.UnixMilli(int64(*prComment.Comment.UpdatedAt))
		bitbucketPrComment.UpdatedAt = &updatedAt
	}
	return bitbucketPrComment, nil
}
