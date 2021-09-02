package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/utils"
)

type ApiClientBeforeRequest func(req *http.Request) error
type ApiClientAfterResponse func(res *http.Response) error

type ApiClient struct {
	client        *http.Client
	endpoint      string
	headers       map[string]string
	maxRetry      int
	beforeRequest ApiClientBeforeRequest
	afterReponse  ApiClientAfterResponse
}

func NewApiClient(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
) *ApiClient {
	apiClient := &ApiClient{}
	apiClient.Setup(
		endpoint,
		headers,
		timeout,
		maxRetry,
	)
	return apiClient
}

func (apiClient *ApiClient) Setup(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
) {
	apiClient.client = &http.Client{Timeout: timeout}
	apiClient.SetEndpoint(endpoint)
	apiClient.SetHeaders(headers)
	apiClient.SetMaxRetry(maxRetry)
}

func (apiClient *ApiClient) SetEndpoint(endpoint string) {
	apiClient.endpoint = endpoint
}

func (ApiClient *ApiClient) SetTimeout(timeout time.Duration) {
	ApiClient.client.Timeout = timeout
}

func (ApiClient *ApiClient) SetMaxRetry(maxRetry int) {
	ApiClient.maxRetry = maxRetry
}

func (apiClient *ApiClient) SetHeaders(headers map[string]string) {
	apiClient.headers = headers
}

func (apiClient *ApiClient) SetProxy(proxyUrl string) error {
	pu, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	apiClient.client.Transport = &http.Transport{Proxy: http.ProxyURL(pu)}
	return nil
}

func (apiClient *ApiClient) Do(
	method string,
	path string,
	query *url.Values,
	body *map[string]interface{},
	headers *map[string]string,
) (*http.Response, error) {
	uri := apiClient.endpoint + path

	// append query
	if query != nil {
		queryString := query.Encode()
		if queryString != "" {
			uri += "?" + queryString
		}
	}

	// process body
	var reqBody io.Reader
	if body != nil {
		reqJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(reqJson)
	}
	req, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// populate headers
	if apiClient.headers != nil {
		for name, value := range apiClient.headers {
			req.Header.Set(name, value)
		}
	}
	if headers != nil {
		for name, value := range *headers {
			req.Header.Set(name, value)
		}
	}

	// before send
	if apiClient.beforeRequest != nil {
		err = apiClient.beforeRequest(req)
		if err != nil {
			return nil, err
		}
	}

	logger.Print(fmt.Sprintf("[api-client] %v %v", method, uri))
	res, err := apiClient.client.Do(req)
	if err != nil {
		return nil, err
	}

	// after recieve
	if apiClient.afterReponse != nil {
		err = apiClient.afterReponse(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (apiClient *ApiClient) Get(
	path string,
	query *url.Values,
	headers *map[string]string,
) (*http.Response, error) {
	return apiClient.Do("GET", path, query, nil, headers)
}

type Pagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

type PaginationHandler func(res *http.Response) error

func (ApiClient *ApiClient) FetchWithPagination(path string, query *url.Values, handler PaginationHandler) error {
	if query == nil {
		query = &url.Values{}
	}
	nextStart, total, pageSize := 0, 1, 100

	pageQuery := &url.Values{}
	*pageQuery = *query
	pageQuery.Set("maxResults", "0")
	res, err := ApiClient.Get(path, query, nil)
	if err != nil {
		return err
	}
	apiResponse := &Pagination{}
	err = UnmarshalResponse(res, apiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	total = apiResponse.Total

	scheduler, err := utils.NewWorkerScheduler(10, 50)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	for nextStart < total {
		nextStartTmp := nextStart
		err = scheduler.Submit(func() error {
			// fetch page
			detailQuery := &url.Values{}
			*detailQuery = *query
			detailQuery.Set("maxResults", strconv.Itoa(pageSize))
			detailQuery.Set("startAt", strconv.Itoa(nextStartTmp))
			res, err := ApiClient.Get(path, query, nil)
			if err != nil {
				return err
			}

			// call page handler
			err = handler(res)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			logger.Info("api client page loaded", map[string]interface{}{
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
	scheduler.WaitUntilFinish()
	return nil
}

func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(resBody, &v)
}
