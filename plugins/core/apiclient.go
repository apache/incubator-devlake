package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
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
	return &ApiClient{
		client:        &http.Client{Timeout: timeout},
		endpoint:      endpoint,
		headers:       headers,
		maxRetry:      maxRetry,
		beforeRequest: nil,
		afterReponse:  nil,
	}
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

// UnmarshalResponse - Uses the built in json decoder to
// decode the response Body into the provided object
// This saves us having to conver the io.Reader into
// a []byte.
func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(v)
}
