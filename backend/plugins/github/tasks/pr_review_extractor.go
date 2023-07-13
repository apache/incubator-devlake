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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPullRequestReviewsMeta)
}

var ExtractApiPullRequestReviewsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestReviews",
	EntryPoint:       ExtractApiPullRequestReviews,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestReviewers data into tool layer table github_reviewers",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{RAW_PR_REVIEW_TABLE},
	ProductTables: []string{
		models.GithubRepoAccount{}.TableName(),
		models.GithubReviewer{}.TableName(),
		models.GithubPrReview{}.TableName()},
}

type PullRequestReview struct {
	GithubId    int `json:"id"`
	User        *GithubAccountResponse
	Body        string
	State       string
	CommitId    string          `json:"commit_id"`
	SubmittedAt api.Iso8601Time `json:"submitted_at"`
}

func ExtractApiPullRequestReviews(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PR_REVIEW_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			apiPullRequestReview := &PullRequestReview{}
			if strings.HasPrefix(string(row.Data), "{\"message\": \"Not Found\"") {
				return nil, nil
			}
			err := errors.Convert(json.Unmarshal(row.Data, apiPullRequestReview))
			if err != nil {
				return nil, err
			}
			if apiPullRequestReview.State == "PENDING" || apiPullRequestReview.User == nil {
				return nil, nil
			}
			pull := &SimplePr{}
			err = errors.Convert(json.Unmarshal(row.Input, pull))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)

			githubReviewer := &models.GithubReviewer{
				ConnectionId:  data.Options.ConnectionId,
				PullRequestId: pull.GithubId,
			}

			githubPrReview := &models.GithubPrReview{
				ConnectionId:   data.Options.ConnectionId,
				GithubId:       apiPullRequestReview.GithubId,
				Body:           apiPullRequestReview.Body,
				State:          apiPullRequestReview.State,
				CommitSha:      apiPullRequestReview.CommitId,
				GithubSubmitAt: apiPullRequestReview.SubmittedAt.ToNullableTime(),
				PullRequestId:  pull.GithubId,
			}

			if apiPullRequestReview.User != nil {
				githubReviewer.GithubId = apiPullRequestReview.User.Id
				githubReviewer.Login = apiPullRequestReview.User.Login

				githubPrReview.AuthorUserId = apiPullRequestReview.User.Id
				githubPrReview.AuthorUsername = apiPullRequestReview.User.Login

				githubUser, err := convertAccount(apiPullRequestReview.User, data.Options.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubUser)
			}

			results = append(results, githubReviewer)
			results = append(results, githubPrReview)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
