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

const RAW_ISSUE_TABLE = "github_api_issues"

// this struct should be moved to `gitub_api_common.go`
type GithubApiParams struct {
	Owner string
	Repo  string
}

var CollectApiIssuesMeta = core.SubTaskMeta{
	Name:             "collectApiIssues",
	EntryPoint:       CollectApiIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Github api",
}

func CollectApiIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	if since == nil {
		var latestUpdated models.GithubIssue
		err := db.Model(&latestUpdated).
			Where("repo_id = ?", data.Repo.GithubId).
			Order("github_updated_at DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return fmt.Errorf("failed to get latest github issue record: %w", err)
		}
		if latestUpdated.GithubId > 0 {
			since = &latestUpdated.GithubUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("direction", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		/*
			Some api might do pagination by http headers
		*/
		//Header: func(pager *core.Pager) http.Header {
		//},
		/*
			Sometimes, we need to collect data based on previous collected data, like jira changelog, it requires
			issue_id as part of the url.
			We can mimic `stdin` design, to accept a `Input` function which produces a `Iterator`, collector
			should iterate all records, and do data-fetching for each on, either in parallel or sequential order
			UrlTemplate: "api/3/issue/{{ Input.ID }}/changelog"
		*/
		//Input: databaseIssuesIterator,
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
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
