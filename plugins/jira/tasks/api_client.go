package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiClient struct {
	*helper.ApiAsyncClient
}

func NewJiraApiClient(taskCtx core.TaskContext, source *models.JiraSource) (*JiraApiClient, error) {
	// load configuration
	encKey := taskCtx.GetConfig(core.EncodeKeyEnvStr)
	auth, err := core.Decrypt(encKey, source.BasicAuthEncoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt Auth Token: %w", err)
	}

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: source.RateLimit,
	}
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		source.Proxy,
		source.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", auth),
		},
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	jiraApiClient := &JiraApiClient{
		asyncApiClient,
	}

	jiraApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})
	return jiraApiClient, nil
}

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

// Deprecated
type JiraPaginationHandler func(res *http.Response) error

// Deprecated
type JiraSearchPaginationHandler func(res *http.Response) (int, error)

// Deprecated
func (jiraApiClient *JiraApiClient) FetchPages(path string, query url.Values, handler JiraPaginationHandler) error {
	if query == nil {
		query = url.Values{}
	}
	nextStart, pageSize := 0, 100

	// 获取issue总数
	// get issue count
	pageQuery := url.Values{}
	for key, value := range query {
		pageQuery[key] = value
	}
	pageQuery.Set("maxResults", "0")
	// make a call to the api just to get the paging details
	res, err := jiraApiClient.Get(path, pageQuery, nil)
	if err != nil {
		return err
	}
	jiraApiResponse := &JiraPagination{}
	err = helper.UnmarshalResponse(res, jiraApiResponse)
	if err != nil {
		return nil
	}
	total := jiraApiResponse.Total

	for nextStart < total {
		nextStartTmp := nextStart
		queryCopy := url.Values{}
		for key, value := range query {
			queryCopy[key] = value
		}
		queryCopy.Set("maxResults", strconv.Itoa(pageSize))
		queryCopy.Set("startAt", strconv.Itoa(nextStartTmp))
		err = jiraApiClient.GetAsync(path, queryCopy, nil, handler)
		if err != nil {
			return err
		}

		nextStart += pageSize
	}
	jiraApiClient.WaitAsync()
	return nil
}

// Deprecated
// FetchWithoutPaginationHeaders uses pagination in a different way than FetchPages.
// We set the pagination params to what we want, and then we just keep making requests
// until the response array is empty. This is why we need to check the "next" variable
// on the handler, and the handler that is passed in needs to return a boolean to tell
// us whether or not to continue making requests. This is why we created JiraSearchPaginationHandler.
func (jiraApiClient *JiraApiClient) FetchWithoutPaginationHeaders(
	path string,
	query url.Values,
	handler JiraSearchPaginationHandler,
) error {
	// this method is sequential, data would be collected page by page
	// because all collectors using this method do not contains many records.
	if query == nil {
		query = url.Values{}
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
