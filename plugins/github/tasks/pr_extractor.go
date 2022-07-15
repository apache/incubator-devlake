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
	"runtime/debug"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPullRequestsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequests",
	EntryPoint:       ExtractApiPullRequests,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequests data into tool layer table github_pull_requests",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

type GithubApiPullRequest struct {
	GithubId int `json:"id"`
	Number   int
	State    string
	Title    string
	Body     json.RawMessage
	HtmlUrl  string `json:"html_url"`
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignee        *GithubAccountResponse
	User            *GithubAccountResponse
	ClosedAt        *helper.Iso8601Time `json:"closed_at"`
	MergedAt        *helper.Iso8601Time `json:"merged_at"`
	GithubCreatedAt helper.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt helper.Iso8601Time  `json:"updated_at"`
	MergeCommitSha  string              `json:"merge_commit_sha"`
	Head            struct {
		Ref string
		Sha string
	}
	Base struct {
		Ref string
		Sha string
	}
}

func ExtractApiPullRequests(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	config := data.Options.TransformationRules
	var labelTypeRegex *regexp.Regexp
	var labelComponentRegex *regexp.Regexp
	var prType = config.PrType
	var err error
	if len(prType) > 0 {
		labelTypeRegex, err = regexp.Compile(prType)
		if err != nil {
			return fmt.Errorf("regexp Compile prType failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var prComponent = config.PrComponent
	if len(prComponent) > 0 {
		labelComponentRegex, err = regexp.Compile(prComponent)
		if err != nil {
			return fmt.Errorf("regexp Compile prComponent failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}

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
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			rawL := &GithubApiPullRequest{}
			err := json.Unmarshal(row.Data, rawL)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)
			if rawL.GithubId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			githubPr, err := convertGithubPullRequest(rawL, data.Options.ConnectionId, data.Repo.GithubId)
			if err != nil {
				return nil, err
			}
			if rawL.User != nil {
				githubUser, err := convertAccount(rawL.User, data.Repo.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubUser)
				githubPr.AuthorName = githubUser.Login
				githubPr.AuthorId = githubUser.AccountId
			}
			for _, label := range rawL.Labels {
				results = append(results, &models.GithubPullRequestLabel{
					ConnectionId: data.Options.ConnectionId,
					PullId:       githubPr.GithubId,
					LabelName:    label.Name,
				})
				// if pr.Type has not been set and prType is set in .env, process the below
				if labelTypeRegex != nil {
					groups := labelTypeRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubPr.Type = groups[1]
					}
				}

				// if pr.Component has not been set and prComponent is set in .env, process
				if labelComponentRegex != nil {
					groups := labelComponentRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubPr.Component = groups[1]
					}
				}
			}
			results = append(results, githubPr)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGithubPullRequest(pull *GithubApiPullRequest, connId uint64, repoId int) (*models.GithubPullRequest, error) {
	githubPull := &models.GithubPullRequest{
		ConnectionId:    connId,
		GithubId:        pull.GithubId,
		RepoId:          repoId,
		Number:          pull.Number,
		State:           pull.State,
		Title:           pull.Title,
		Url:             pull.HtmlUrl,
		GithubCreatedAt: pull.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: pull.GithubUpdatedAt.ToTime(),
		ClosedAt:        helper.Iso8601TimeToTime(pull.ClosedAt),
		MergedAt:        helper.Iso8601TimeToTime(pull.MergedAt),
		MergeCommitSha:  pull.MergeCommitSha,
		Body:            string(pull.Body),
		BaseRef:         pull.Base.Ref,
		BaseCommitSha:   pull.Base.Sha,
		HeadRef:         pull.Head.Ref,
		HeadCommitSha:   pull.Head.Sha,
	}
	return githubPull, nil
}
