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
type JiraSearchPaginationHandler func(res *http.Response) (bool, error)

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

// FetchWithoutPaginationHeaders uses pagination in a different way than FetchPages.
// We set the pagination params to what we want, and then we just keep making requests
// until the response array is empty. This is why we need to check the "next" variable
// on the handler, and the handler that is passed in needs to return a boolean to tell
// us whether or not to continue making requests. This is why we created JiraSearchPaginationHandler.
func (jiraApiClient *JiraApiClient) FetchWithoutPaginationHeaders(
	path string,
	query *url.Values,
	handler JiraSearchPaginationHandler,
) error {
	if query == nil {
		query = &url.Values{}
	}
	// these are the names from the jira search api for pagination
	// eg: https://merico.atlassian.net/rest/api/2/user/assignable/search?project=EE&maxResults=100&startAt=1
	startAt, maxResults := 1, 100

	query.Set("maxResults", fmt.Sprintf("%v", maxResults))
	var next bool = true
	for next {
		nextStartTmp := startAt
		// get page
		query.Set("startAt", strconv.Itoa(nextStartTmp))
		res, err := jiraApiClient.Get(path, query, nil)
		if res.StatusCode == 401 {
			fmt.Println("User does not have access to project users")
		}
		if err != nil {
			return err
		}

		// call page handler
		next, err = handler(res)
		if !next {
			// Done user collection
			return nil
		}
		if err != nil {
			logger.Error("Error: ", err)
			return err
		}
		startAt += maxResults
	}
	return nil
}
