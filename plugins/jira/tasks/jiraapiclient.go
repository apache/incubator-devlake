package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/utils"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiClient struct {
	core.ApiClient
}

func NewJiraApiClient(endpoint string, auth string) *JiraApiClient {
	jiraApiClient := &JiraApiClient{}
	jiraApiClient.Setup(
		endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", auth),
		},
		10*time.Second,
		3,
	)
	return jiraApiClient
}

func NewJiraApiClientBySourceId(sourceId uint64) (*JiraApiClient, error) {
	jiraSource := &models.JiraSource{}
	err := lakeModels.Db.First(jiraSource, sourceId).Error
	if err != nil {
		return nil, err
	}
	return NewJiraApiClient(jiraSource.Endpoint, jiraSource.BasicAuthEncoded), nil
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
	for key, value := range *query {
		(*pageQuery)[key] = value
	}
	pageQuery.Set("maxResults", "0")
	// make a call to the api just to get the paging details
	res, err := jiraApiClient.Get(path, pageQuery, nil)
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
			for key, value := range *query {
				(*detailQuery)[key] = value
			}
			detailQuery.Set("maxResults", strconv.Itoa(pageSize))
			detailQuery.Set("startAt", strconv.Itoa(nextStartTmp))
			res, err := jiraApiClient.Get(path, detailQuery, nil)
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
