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
	"regexp"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPrReviewCommentsMeta)
}

var ExtractApiPrReviewCommentsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPrReviewComments",
	EntryPoint:       ExtractApiPrReviewComments,
	EnabledByDefault: true,
	Description: "Extract raw comment data  into tool layer table github_pull_request_comments" +
		"and github_issue_comments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{RAW_PR_REVIEW_COMMENTS_TABLE},
	ProductTables: []string{
		models.GithubRepoAccount{}.TableName(),
		models.GithubPrCommit{}.TableName()},
}

func ExtractApiPrReviewComments(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	db := taskCtx.GetDal()
	prUrlPattern := fmt.Sprintf(`https\:\/\/api\.github\.com\/repos\/%s\/pulls\/(\d+)`, data.Options.Name)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_PR_REVIEW_COMMENTS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
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
			err := errors.Convert(json.Unmarshal(row.Data, &prReviewComment))
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
				GithubCreatedAt: prReviewComment.GithubCreatedAt.ToTime(),
				GithubUpdatedAt: prReviewComment.GithubUpdatedAt.ToTime(),
				Type:            "DIFF",
			}

			prUrlRegex, err := errors.Convert01(regexp.Compile(prUrlPattern))
			if err != nil {
				return nil, errors.Default.Wrap(err, "regexp Compile prUrlPattern failed")
			}
			prId, err := enrichGithubPrComment(data, db, prUrlRegex, prReviewComment.PrUrl)
			if err != nil {
				return nil, errors.Default.Wrap(err, "parse prId failed")
			}
			if prId != 0 {
				githubPrComment.PullRequestId = prId
			}

			if prReviewComment.User != nil {
				githubPrComment.AuthorUserId = prReviewComment.User.Id
				githubPrComment.AuthorUsername = prReviewComment.User.Login

				githubAccount, err := convertAccount(prReviewComment.User, data.Options.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubAccount)
			}

			results = append(results, githubPrComment)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func enrichGithubPrComment(data *GithubTaskData, db dal.Dal, prUrlRegex *regexp.Regexp, prUrl string) (int, errors.Error) {
	groups := prUrlRegex.FindStringSubmatch(prUrl)
	if len(groups) > 0 {
		prNumber, err := strconv.Atoi(groups[1])
		if err != nil {
			return 0, errors.Default.Wrap(err, "parse prId failed")
		}
		pr := &models.GithubPullRequest{}
		err1 := db.First(pr, dal.Where("connection_id = ? and number = ? and repo_id = ?", data.Options.ConnectionId, prNumber, data.Options.GithubId))
		if db.IsErrorNotFound(err1) {
			return 0, nil
		} else if err1 != nil {
			return 0, errors.NotFound.Wrap(err1, "github pull request parse failed ")
		}
		return pr.GithubId, nil
	}
	return 0, nil
}
