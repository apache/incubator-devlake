package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

// this struct should be moved to `jira_api_common.go`
type JiraApiParams struct {
	SourceId uint64
	BoardId  uint64
}

func CollectApiIssues(
	apiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	since time.Time,
) error {
	// --------- skip this part begin --------
	// user didn't specify a time range to sync, try load from database
	if since.IsZero() {
		var latestUpdated models.JiraIssue
		err := lakeModels.Db.Where("source_id = ?", source.ID).Order("updated DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			logger.Error("jira collect issues:  get last sync time failed", err)
			return err
		}
		since = latestUpdated.Updated
	}
	// build jql
	jql := "ORDER BY updated ASC"
	if !since.IsZero() {
		// prepend a time range criteria if `since` was specified, either by user or from database
		jql = fmt.Sprintf("updated >= '%v' %v", since.Format("2006/01/02 15:04"), jql)
	}
	query := &url.Values{}
	query.Set("jql", jql)
	// --------- skip this part end --------

	const SIZE = 100

	collector, err := core.NewApiCollector(core.ApiCollectorArgs{
		ApiClient: apiClient,
		PageSize:  SIZE,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "agile/1.0/board/{{ Params.BoardId }}/issue?startAt={{ Pager.StartAt }}&maxResults={{ Pager.Size }}",
		/*
			Some api might do pagination by http headers
		*/
		//Header: func(collector *core.ApiCollector) http.Header {
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
			SourceId: source.ID,
			BoardId:  boardId,
		},
		/*
			Reduce records come out from api
		*/
		Query: query,
		/*
			ApiCollector will try to `UnmarshalResponse` into this variable, and handle failure
		*/
		ResponseData: &JiraApiIssuesResponse{},
		/*
			Accept response and Unmarshalled body, return raw data either a single object or list
		*/
		OnData: func(res *http.Response, body interface{}) (interface{}, error) {
			issuesBody := body.(*JiraApiIssuesResponse)
			return issuesBody.Issues, nil
		},
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
		GetTotalPages: func(res *http.Response, body interface{}) (int, error) {
			issuesBody := body.(*JiraApiIssuesResponse)
			return issuesBody.Total / SIZE, nil
		},
		/*
			Table store raw data
		*/
		Table: "jira_api_issues",
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
