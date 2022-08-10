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
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"gorm.io/gorm"
	"regexp"
	"runtime/debug"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPrReviewCommentsMeta = core.SubTaskMeta{
	Name:             "extractApiPrReviewComments",
	EntryPoint:       ExtractApiPrReviewComments,
	EnabledByDefault: true,
	Description: "Extract raw comment data  into tool layer table github_pull_request_comments" +
		"and github_issue_comments",
	DomainTypes: []string{core.DOMAIN_TYPE_CODE},
}

func ExtractApiPrReviewComments(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubGraphqlTaskData)
	prUrlPattern := fmt.Sprintf(`https\:\/\/api\.github\.com\/repos\/%s\/%s\/pulls\/(\d+)`, data.Options.Owner, data.Options.Repo)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubGraphqlApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PR_REVIEW_COMMENTS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var prReviewComment struct {
				GithubId        int `json:"id"`
				Body            json.RawMessage
				User            *GithubAccountResponse
				PrUrl           string             `json:"pull_request_url"`
				GithubCreatedAt helper.Iso8601Time `json:"created_at"`
				GithubUpdatedAt helper.Iso8601Time `json:"updated_at"`
				CommitId        string             `json:"commit_id"`
				PrReviewId      int                `json:"pull_request_review_id"`
			}
			err := json.Unmarshal(row.Data, &prReviewComment)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1)
			githubPrComment := &models.GithubPrComment{
				ConnectionId:    data.Options.ConnectionId,
				GithubId:        prReviewComment.GithubId,
				Body:            string(prReviewComment.Body),
				CommitSha:       prReviewComment.CommitId,
				ReviewId:        prReviewComment.PrReviewId,
				AuthorUsername:  prReviewComment.User.Login,
				AuthorUserId:    prReviewComment.User.Id,
				GithubCreatedAt: prReviewComment.GithubCreatedAt.ToTime(),
				GithubUpdatedAt: prReviewComment.GithubUpdatedAt.ToTime(),
				Type:            "DIFF",
			}

			prUrlRegex, err := regexp.Compile(prUrlPattern)
			if err != nil {
				return nil, fmt.Errorf("regexp Compile prUrlPattern failed:[%s] stack:[%s]", err.Error(), debug.Stack())
			}
			if prUrlRegex != nil {
				groups := prUrlRegex.FindStringSubmatch(prReviewComment.PrUrl)
				if len(groups) > 0 {
					prNumber, err := strconv.Atoi(groups[1])
					if err != nil {
						return nil, fmt.Errorf("parse prId failed:[%s] stack:[%s]", err.Error(), debug.Stack())
					}
					pr := &models.GithubPullRequest{}
					err = taskCtx.GetDal().First(pr, dal.Where("connection_id = ? and number = ? and repo_id = ?", data.Options.ConnectionId, prNumber, data.Repo.GithubId))
					if err != nil && err != gorm.ErrRecordNotFound {
						return nil, err
					}
					githubPrComment.PullRequestId = pr.GithubId
				}
			}
			results = append(results, githubPrComment)
			githubAccounts, err := convertRestPreAccount(prReviewComment.User, data.Repo.GithubId, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, githubAccounts...)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
