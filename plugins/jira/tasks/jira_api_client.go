package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/utils"
)

type JiraApiClient struct {
	core.ApiClient
}

func NewJiraApiClient(
	endpoint string,
	auth string,
	proxy string,
	scheduler *utils.WorkerScheduler,
	logger core.Logger,
) *JiraApiClient {
	jiraApiClient := &JiraApiClient{}
	jiraApiClient.Setup(
		endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", auth),
		},
		20*time.Second,
		3,
		scheduler,
	)
	jiraApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})
	if proxy != "" {
		err := jiraApiClient.SetProxy(proxy)
		if err != nil {
			panic(err)
		}
	}
	if proxy != "" {
		err := jiraApiClient.SetProxy(proxy)
		if err != nil {
			panic(err)
		}
	}
	jiraApiClient.SetLogger(logger)
	return jiraApiClient
}

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

type JiraPaginationHandler func(res *http.Response) error
type JiraSearchPaginationHandler func(res *http.Response) (int, error)

func (jiraApiClient *JiraApiClient) FetchPages(path string, query *url.Values, handler JiraPaginationHandler) error {
	if query == nil {
		query = &url.Values{}
	}
	nextStart, pageSize := 0, 100

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
		return nil
	}
	total := jiraApiResponse.Total

	for nextStart < total {
		nextStartTmp := nextStart
		queryCopy := url.Values{}
		for key, value := range *query {
			queryCopy[key] = value
		}
		queryCopy.Set("maxResults", strconv.Itoa(pageSize))
		queryCopy.Set("startAt", strconv.Itoa(nextStartTmp))
		err = jiraApiClient.GetAsync(path, &queryCopy, nil, handler)
		if err != nil {
			return err
		}

		nextStart += pageSize
	}
	jiraApiClient.WaitAsync()
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
	// this method is sequential, data would be collected page by page
	// because all collectors using this method do not contains many records.
	if query == nil {
		query = &url.Values{}
	}
	// these are the names from the jira search api for pagination
	// eg: https://merico.atlassian.net/rest/api/2/user/assignable/search?project=EE&maxResults=100&startAt=1
	startAt, maxResults := 0, 100

	query.Set("maxResults", fmt.Sprintf("%v", maxResults))
	// some jira api like 'agile sprints' maxResults upper limit may less that 100.
	// it's not safe to assume all api accept maxResults=100 as a valid parameter.
	// should let handler return the actual received length and use it to increase `startAt`,
	// and because the size of next page always smaller than current one, we could merge
	// `maxResults` and `length` into one variable
	for maxResults > 0 {
		// get page
		query.Set("startAt", strconv.Itoa(startAt))
		res, err := jiraApiClient.Get(path, query, nil)
		if err != nil {
			return err
		}
		if res.StatusCode == 401 {
			res.Body.Close()
			return fmt.Errorf("User does not have access to project users")
		}

		// call page handler
		maxResults, err = handler(res)
		res.Body.Close()
		if err != nil {
			return err
		}
		startAt += maxResults
	}
	return nil
}

func (jiraApiClient *JiraApiClient) GetJiraServerInfo() (*models.JiraServerInfo, int, error) {
	res, err := jiraApiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, 0, err
	}
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, res.StatusCode, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	serverInfo := &models.JiraServerInfo{}
	err = core.UnmarshalResponse(res, serverInfo)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return serverInfo, res.StatusCode, nil
}

func (jiraApiClient *JiraApiClient) GetMyselfInfo() (*models.ApiMyselfResponse, error) {
	res, err := jiraApiClient.Get("api/3/myself", nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	myselfFromApi := &models.ApiMyselfResponse{}
	err = core.UnmarshalResponse(res, myselfFromApi)
	if err != nil {
		return nil, err
	}
	return myselfFromApi, nil
}
