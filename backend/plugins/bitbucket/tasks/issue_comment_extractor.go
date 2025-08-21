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

type BitbucketIssueCommentsResponse struct {
	Type        string     `json:"type"`
	BitbucketId int        `json:"id"`
	CreatedOn   time.Time  `json:"created_on"`
	UpdatedOn   *time.Time `json:"updated_on"`
	Content     struct {
		Type string `json:"type"`
		Raw  string `json:"raw"`
	} `json:"content"`
	User  *BitbucketAccountResponse `json:"user"`
	Issue struct {
		//Type string `json:"type"`
		Id int `json:"id"`
		//Repository *BitbucketApiRepo `json:"repository"`
		//Title string `json:"title"`
	}
}

var ExtractApiIssueCommentsMeta = plugin.SubTaskMeta{
	Name:             "Extract Issue Comments",
	EntryPoint:       ExtractApiIssueComments,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw issue comments data into tool layer table BitbucketIssueComments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractApiIssueComments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_COMMENTS_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			issueComment := &BitbucketIssueCommentsResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, issueComment))
			if err != nil {
				return nil, err
			}

			toolIssueComment, err := convertIssueComment(issueComment)
			toolIssueComment.ConnectionId = data.Options.ConnectionId
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 2)

			if issueComment.User != nil {
				bitbucketUser, err := convertAccount(issueComment.User, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				toolIssueComment.AuthorId = bitbucketUser.AccountId
				toolIssueComment.AuthorName = bitbucketUser.DisplayName
				results = append(results, bitbucketUser)
			}

			results = append(results, toolIssueComment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertIssueComment(issueComment *BitbucketIssueCommentsResponse) (*models.BitbucketIssueComment, errors.Error) {
	bitbucketIssueComment := &models.BitbucketIssueComment{
		BitbucketId:        issueComment.BitbucketId,
		AuthorId:           issueComment.User.AccountId,
		IssueId:            issueComment.Issue.Id,
		AuthorName:         issueComment.User.DisplayName,
		BitbucketCreatedAt: issueComment.CreatedOn,
		BitbucketUpdatedAt: issueComment.UpdatedOn,
		Type:               issueComment.Type,
		Body:               issueComment.Content.Raw,
	}
	return bitbucketIssueComment, nil
}
