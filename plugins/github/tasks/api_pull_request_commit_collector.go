package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/merico-dev/lake/plugins/helper"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
)

const RAW_PULL_REQUEST_COMMIT_TABLE = "github_api_pull_request_commits"

// this struct should be moved to `gitub_api_common.go`

var CollectApiPullRequestCommitsMeta = core.SubTaskMeta{
	Name:             "collectApiPullRequestCommits",
	EntryPoint:       CollectApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Collect PullRequestCommits data from Github api",
}

type SimplePr struct {
	Number   int
	GithubId int
}

func CollectApiPullRequestCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	incremental := false

	cursor, err := db.Model(&models.GithubPullRequest{}).
		Where("repo_id = ?", data.Repo.GithubId).
		Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(SimplePr{}))
	if err != nil {
		return err
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
			Table: RAW_PULL_REQUEST_COMMIT_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		Input:       iterator,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/pulls/{{ .Input.Number }}/commits",

		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
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
