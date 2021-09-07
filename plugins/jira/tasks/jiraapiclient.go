package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/utils"
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

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

type JiraPaginationHandler func(res *http.Response) error

func (jiraApiClient *JiraApiClient) FetchPages(scheduler *utils.WorkerScheduler, path string, query *url.Values, handler JiraPaginationHandler) error {
	if query == nil {
		query = &url.Values{}
	}
	nextStart, total, pageSize := 0, 1, 100

	// 获取issue总数
	// get issue count
	pageQuery := &url.Values{}
	*pageQuery = *query
	pageQuery.Set("maxResults", "0")
	// make a call to the api just to get the paging details
	res, err := jiraApiClient.Get(path, query, nil)
	if err != nil {
		return err
	}
	jiraApiResponse := &JiraPagination{}
	err = core.UnmarshalResponse(res, jiraApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}
	total = jiraApiResponse.Total

	for nextStart < total {
		nextStartTmp := nextStart
		err = scheduler.Submit(func() error {
			// fetch page
			detailQuery := &url.Values{}
			*detailQuery = *query
			detailQuery.Set("maxResults", strconv.Itoa(pageSize))
			detailQuery.Set("startAt", strconv.Itoa(nextStartTmp))
			res, err := jiraApiClient.Get(path, query, nil)
			if err != nil {
				return err
			}

			// call page handler
			err = handler(res)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			logger.Info("jira api client page loaded", map[string]interface{}{
				"path":      path,
				"nextStart": nextStartTmp,
				"total":     total,
			})
			return nil
		})
		if err != nil {
			return err
		}
		nextStart += pageSize
	}
	return nil
}
