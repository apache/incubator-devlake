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
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

const RAW_PRS_TABLE = "github_graphql_prs"

type GraphqlQueryPrWrapper struct {
	RateLimit struct {
		Cost int
	}
	// now it orderBy UPDATED_AT and use cursor pagination
	// It may miss some PRs updated when collection.
	// Because these missed PRs will be collected on next, But it's not enough.
	// So Next Millstone(0.17) we should change it to filter by CREATE_AT + collect detail
	Repository struct {
		PullRequests struct {
			PageInfo   *api.GraphqlQueryPageInfo
			Prs        []GraphqlQueryPr `graphql:"nodes"`
			TotalCount graphql.Int
		} `graphql:"pullRequests(first: $pageSize, after: $skipCursor, orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GraphqlQueryPr struct {
	DatabaseId int
	Number     int
	State      string
	Title      string
	Body       string
	Url        string
	Labels     struct {
		Nodes []struct {
			Id   string
			Name string
		}
	} `graphql:"labels(first: 100)"`
	Author    *GraphqlInlineAccountQuery
	Assignees struct {
		// FIXME now domain layer just support one assignee
		Assignees []GraphqlInlineAccountQuery `graphql:"nodes"`
	} `graphql:"assignees(first: 1)"`
	ClosedAt    *time.Time
	MergedAt    *time.Time
	UpdatedAt   time.Time
	CreatedAt   time.Time
	MergeCommit *struct {
		Oid string
	}
	HeadRefName string
	HeadRefOid  string
	BaseRefName string
	BaseRefOid  string
	Commits     struct {
		PageInfo   *api.GraphqlQueryPageInfo
		Nodes      []GraphqlQueryCommit `graphql:"nodes"`
		TotalCount graphql.Int
	} `graphql:"commits(first: 100)"`
	Reviews struct {
		TotalCount graphql.Int
		Nodes      []GraphqlQueryReview `graphql:"nodes"`
	} `graphql:"reviews(first: 100)"`
}

type GraphqlQueryReview struct {
	Body       string
	Author     *GraphqlInlineAccountQuery
	State      string `json:"state"`
	DatabaseId int    `json:"databaseId"`
	Commit     struct {
		Oid string
	}
	SubmittedAt *time.Time `json:"submittedAt"`
}

type GraphqlQueryCommit struct {
	Commit struct {
		Oid     string
		Message string
		Author  struct {
			Name  string
			Email string
			Date  time.Time
			User  *GraphqlInlineAccountQuery
		}
		Committer struct {
			Date  time.Time
			Email string
			Name  string
		}
	}
	Url string
}

var CollectPrMeta = plugin.SubTaskMeta{
	Name:             "CollectPr",
	EntryPoint:       CollectPr,
	EnabledByDefault: true,
	Description:      "Collect Pr data from GithubGraphql api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

var _ plugin.SubTaskEntryPoint = CollectPr

func CollectPr(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*tasks.GithubTaskData)
	config := data.Options.GithubTransformationRule
	var labelTypeRegex *regexp.Regexp
	var labelComponentRegex *regexp.Regexp
	var err errors.Error
	if config != nil && len(config.PrType) > 0 {
		labelTypeRegex, err = errors.Convert01(regexp.Compile(config.PrType))
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prType failed")
		}
	}
	if config != nil && len(config.PrComponent) > 0 {
		labelComponentRegex, err = errors.Convert01(regexp.Compile(config.PrComponent))
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prComponent failed")
		}
	}

	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: tasks.GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_PRS_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}

	incremental := collectorWithState.IsIncremental()

	err = collectorWithState.InitGraphQLCollector(api.GraphqlCollectorArgs{
		GraphqlClient: data.GraphqlClient,
		PageSize:      10,
		Incremental:   incremental,
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		BuildQuery: func(reqData *api.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			query := &GraphqlQueryPrWrapper{}
			if reqData == nil {
				return query, map[string]interface{}{}, nil
			}
			ownerName := strings.Split(data.Options.Name, "/")
			variables := map[string]interface{}{
				"pageSize":   graphql.Int(reqData.Pager.Size),
				"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
				"owner":      graphql.String(ownerName[0]),
				"name":       graphql.String(ownerName[1]),
			}
			return query, variables, nil
		},
		GetPageInfo: func(iQuery interface{}, args *api.GraphqlCollectorArgs) (*api.GraphqlQueryPageInfo, error) {
			query := iQuery.(*GraphqlQueryPrWrapper)
			return query.Repository.PullRequests.PageInfo, nil
		},
		ResponseParser: func(iQuery interface{}, variables map[string]interface{}) ([]interface{}, error) {
			query := iQuery.(*GraphqlQueryPrWrapper)
			prs := query.Repository.PullRequests.Prs

			results := make([]interface{}, 0, 1)
			isFinish := false
			for _, rawL := range prs {
				// collect data even though in increment mode because of updating existing data
				if collectorWithState.TimeAfter != nil && !collectorWithState.TimeAfter.Before(rawL.UpdatedAt) {
					isFinish = true
					break
				}
				githubPr, err := convertGithubPullRequest(rawL, data.Options.ConnectionId, data.Options.GithubId)
				if err != nil {
					return nil, err
				}
				if rawL.Author != nil {
					githubUser, err := convertGraphqlPreAccount(*rawL.Author, data.Options.GithubId, data.Options.ConnectionId)
					if err != nil {
						return nil, err
					}
					results = append(results, githubUser)
				}
				for _, label := range rawL.Labels.Nodes {
					results = append(results, &models.GithubPrLabel{
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

				for _, apiPullRequestReview := range rawL.Reviews.Nodes {
					if apiPullRequestReview.State != "PENDING" {
						githubReviewer := &models.GithubReviewer{
							ConnectionId:  data.Options.ConnectionId,
							PullRequestId: githubPr.GithubId,
						}

						githubPrReview := &models.GithubPrReview{
							ConnectionId:   data.Options.ConnectionId,
							GithubId:       apiPullRequestReview.DatabaseId,
							Body:           apiPullRequestReview.Body,
							State:          apiPullRequestReview.State,
							CommitSha:      apiPullRequestReview.Commit.Oid,
							GithubSubmitAt: apiPullRequestReview.SubmittedAt,

							PullRequestId: githubPr.GithubId,
						}

						if apiPullRequestReview.Author != nil {
							githubReviewer.GithubId = apiPullRequestReview.Author.Id
							githubReviewer.Login = apiPullRequestReview.Author.Login

							githubPrReview.AuthorUserId = apiPullRequestReview.Author.Id
							githubPrReview.AuthorUsername = apiPullRequestReview.Author.Login

							githubUser, err := convertGraphqlPreAccount(*apiPullRequestReview.Author, data.Options.GithubId, data.Options.ConnectionId)
							if err != nil {
								return nil, err
							}
							results = append(results, githubUser)
						}

						results = append(results, githubReviewer)
						results = append(results, githubPrReview)
					}
				}

				for _, apiPullRequestCommit := range rawL.Commits.Nodes {
					githubCommit, err := convertPullRequestCommit(apiPullRequestCommit)
					if err != nil {
						return nil, err
					}
					results = append(results, githubCommit)

					githubPullRequestCommit := &models.GithubPrCommit{
						ConnectionId:       data.Options.ConnectionId,
						CommitSha:          apiPullRequestCommit.Commit.Oid,
						PullRequestId:      githubPr.GithubId,
						CommitAuthorName:   githubCommit.AuthorName,
						CommitAuthorEmail:  githubCommit.AuthorEmail,
						CommitAuthoredDate: githubCommit.AuthoredDate,
					}
					if err != nil {
						return nil, err
					}
					results = append(results, githubPullRequestCommit)

					if apiPullRequestCommit.Commit.Author.User != nil {
						githubUser, err := convertGraphqlPreAccount(*apiPullRequestCommit.Commit.Author.User, data.Options.GithubId, data.Options.ConnectionId)
						if err != nil {
							return nil, err
						}
						results = append(results, githubUser)
					}
				}
			}

			if isFinish {
				return results, api.ErrFinishCollect
			} else {
				return results, nil
			}
		},
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

func convertGithubPullRequest(pull GraphqlQueryPr, connId uint64, repoId int) (*models.GithubPullRequest, errors.Error) {
	githubPull := &models.GithubPullRequest{
		ConnectionId:    connId,
		GithubId:        pull.DatabaseId,
		RepoId:          repoId,
		Number:          pull.Number,
		State:           pull.State,
		Title:           pull.Title,
		Url:             pull.Url,
		GithubCreatedAt: pull.CreatedAt,
		GithubUpdatedAt: pull.UpdatedAt,
		ClosedAt:        pull.ClosedAt,
		MergedAt:        pull.MergedAt,
		Body:            pull.Body,
		BaseRef:         pull.BaseRefName,
		BaseCommitSha:   pull.BaseRefOid,
		HeadRef:         pull.HeadRefName,
		HeadCommitSha:   pull.HeadRefOid,
	}
	if pull.MergeCommit != nil {
		githubPull.MergeCommitSha = pull.MergeCommit.Oid
	}
	if pull.Author != nil {
		githubPull.AuthorName = pull.Author.Login
		githubPull.AuthorId = pull.Author.Id
	}
	return githubPull, nil
}

func convertPullRequestCommit(prCommit GraphqlQueryCommit) (*models.GithubCommit, errors.Error) {
	githubCommit := &models.GithubCommit{
		Sha:            prCommit.Commit.Oid,
		Message:        prCommit.Commit.Message,
		AuthorName:     prCommit.Commit.Author.Name,
		AuthorEmail:    prCommit.Commit.Author.Email,
		AuthoredDate:   prCommit.Commit.Author.Date,
		CommitterName:  prCommit.Commit.Committer.Name,
		CommitterEmail: prCommit.Commit.Committer.Email,
		CommittedDate:  prCommit.Commit.Committer.Date,
		Url:            prCommit.Url,
	}
	if prCommit.Commit.Author.User != nil {
		githubCommit.AuthorId = prCommit.Commit.Author.User.Id
	}
	return githubCommit, nil
}
