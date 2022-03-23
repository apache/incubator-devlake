package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/plugins/helper"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
)

const RAW_COMMENTS_TABLE = "github_api_comments"

// this struct should be moved to `gitub_api_common.go`

func CollectApiComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	// actually, for github pull, since doesn't make any sense, github pull api doesn't support it
	if since == nil {
		var latestUpdatedIssueComt models.GithubIssueComment
		err := db.Model(&latestUpdatedIssueComt).
			Joins("left join github_issues on github_issues.github_id = github_issue_comments.issue_id").
			Where("github_issues.repo_id = ?", data.Repo.GithubId).
			Order("github_updated_at DESC").Limit(1).Find(&latestUpdatedIssueComt).Error
		if err != nil {
			return fmt.Errorf("failed to get latest github issue record: %w", err)
		}
		var latestUpdatedPrComt models.GithubPullRequestComment
		err = db.Model(&latestUpdatedPrComt).
			Joins("left join github_pull_requests on github_pull_requests.github_id = github_pull_request_comments.pull_request_id").
			Where("github_pull_requests.repo_id = ?", data.Repo.GithubId).
			Order("github_updated_at DESC").Limit(1).Find(&latestUpdatedPrComt).Error
		if err != nil {
			return fmt.Errorf("failed to get latest github issue record: %w", err)
		}
		if latestUpdatedIssueComt.GithubId > 0 && latestUpdatedPrComt.GithubId > 0 {
			if latestUpdatedIssueComt.GithubUpdatedAt.Before(latestUpdatedPrComt.GithubUpdatedAt) {
				since = &latestUpdatedPrComt.GithubUpdatedAt
			} else {
				since = &latestUpdatedIssueComt.GithubUpdatedAt
			}
			incremental = true
		} else if latestUpdatedIssueComt.GithubId > 0 {
			since = &latestUpdatedIssueComt.GithubUpdatedAt
			incremental = true
		} else if latestUpdatedPrComt.GithubId > 0 {
			since = &latestUpdatedPrComt.GithubUpdatedAt
			incremental = true
		}

	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/comments",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectApiCommentsMeta = core.SubTaskMeta{
	Name:             "collectApiComments",
	EntryPoint:       CollectApiComments,
	EnabledByDefault: true,
	Description:      "Collect comments data from Github api",
}
