package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

const RAW_ISSUE_TABLE = "jira_api_issues"

// this struct should be moved to `jira_api_common.go`
type JiraApiParams struct {
	SourceId uint64
	BoardId  uint64
}

type JiraApiRawIssuesResponse struct {
	JiraPagination
	Issues []json.RawMessage `json:"issues"`
}

func CollectApiIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JiraTaskData)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	if since == nil {
		var latestUpdated models.JiraIssue
		err := db.Where("source_id = ?", data.Source.ID).Order("updated DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return fmt.Errorf("failed to get latest jira issue record: %w", err)
		}
		if latestUpdated.IssueId > 0 {
			*since = latestUpdated.Updated
			incremental = true
		}
	}
	// build jql
	// IMPORTANT: we have to keep paginated data in a consistence order to avoid data-missing, if we sort issues by
	//  `updated`, issue will be jumping between pages if it got updated during the collection process
	jql := "ORDER BY created ASC"
	if since != nil {
		// prepend a time range criteria if `since` was specified, either by user or from database
		jql = fmt.Sprintf("updated >= '%v' %v", since.Format("2006/01/02 15:04"), jql)
	}

	const SIZE = 100

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		Ctx:         taskCtx,
		ApiClient:   data.ApiClient,
		PageSize:    SIZE,
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
		UrlTemplate: "agile/1.0/board/{{ .Params.BoardId }}/issue",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(pager *helper.Pager) (*url.Values, error) {
			query := &url.Values{}
			query.Set("jql", jql)
			query.Set("startAt", fmt.Sprintf("%v", pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", pager.Size))
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
			This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
			set of data to be process, for example, we process JiraIssues by Board
		*/
		Params: JiraApiParams{
			SourceId: data.Source.ID,
			BoardId:  data.Options.BoardId,
		},
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
		GetTotalPages: func(res *http.Response) (int, error) {
			body := &JiraPagination{}
			err := core.UnmarshalResponse(res, body)
			if err != nil {
				return 0, err
			}
			return body.Total / SIZE, nil
		},
		/*
			Table store raw data
		*/
		Table: RAW_ISSUE_TABLE,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
