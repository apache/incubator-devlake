package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
)

type JiraApiClient struct {
	core.ApiClient
}

var jiraApiClient *JiraApiClient

func GetJiraApiClient() *JiraApiClient {
	if jiraApiClient == nil {
		jiraApiClient = &JiraApiClient{}
		jiraApiClient.Setup(
			config.V.GetString("JIRA_ENDPOINT"),
			map[string]string{
				"Authorization": fmt.Sprintf("Basic %v", config.V.GetString("JIRA_BASIC_AUTH_ENCODED")),
			},
			10*time.Second,
			3,
		)
	}
	return jiraApiClient
}

type JiraPaginationHandler func(res *http.Response) (*JiraPagination, error)

func (jiraApiClient *JiraApiClient) FetchPages(path string, query *url.Values, handler JiraPaginationHandler) error {
	if query == nil {
		query = &url.Values{}
	}
	nextStart, total, query := 0, 1, &url.Values{}
	for nextStart < total {
		// fetch page
		query.Set("maxResults", "100")
		query.Set("startAt", strconv.Itoa(nextStart))
		res, err := jiraApiClient.Get(path, query, nil)
		if err != nil {
			return err
		}

		// call page handler
		pagination, err := handler(res)
		if err != nil {
			return err
		}

		// next page
		nextStart = pagination.StartAt + pagination.MaxResults
		total = pagination.Total
	}
	return nil
}
