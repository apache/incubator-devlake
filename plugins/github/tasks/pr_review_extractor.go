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

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPullRequestReviewsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequestReviews",
	EntryPoint:       ExtractApiPullRequestReviews,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestReviewers data into tool layer table github_reviewers",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

type PullRequestReview struct {
	GithubId    int `json:"id"`
	User        *GithubAccountResponse
	Body        string
	State       string
	CommitId    string             `json:"commit_id"`
	SubmittedAt helper.Iso8601Time `json:"submitted_at"`
}

func ExtractApiPullRequestReviews(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PR_REVIEW_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiPullRequestReview := &PullRequestReview{}
			if strings.HasPrefix(string(row.Data), "{\"message\": \"Not Found\"") {
				return nil, nil
			}
			err := json.Unmarshal(row.Data, apiPullRequestReview)
			if err != nil {
				return nil, err
			}
			if apiPullRequestReview.State == "PENDING" {
				return nil, nil
			}
			pull := &SimplePr{}
			err = json.Unmarshal(row.Input, pull)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)

			githubReviewer := &models.GithubReviewer{
				ConnectionId:  data.Options.ConnectionId,
				GithubId:      apiPullRequestReview.User.Id,
				Login:         apiPullRequestReview.User.Login,
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
				AuthorUsername: apiPullRequestReview.User.Login,
				AuthorUserId:   apiPullRequestReview.User.Id,
			}

			results = append(results, githubReviewer)
			results = append(results, githubPrReview)
			githubUser, err := convertAccount(apiPullRequestReview.User, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, githubUser)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
