package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/merico-dev/lake/logger"
)

// This is for multiple token functionality so we can loop through an array of tokens.
var tokenIndex int = 0

type ApiClientBeforeRequest func(req *http.Request) error
type ApiClientAfterResponse func(res *http.Response) error

type ApiClient struct {
	client        *http.Client
	endpoint      string
	headers       map[string]string
	maxRetry      int
	beforeRequest ApiClientBeforeRequest
	afterReponse  ApiClientAfterResponse
	tokens        []string
}

func NewApiClient(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
	tokens []string,
) *ApiClient {
	apiClient := &ApiClient{}
	apiClient.Setup(
		endpoint,
		headers,
		timeout,
		maxRetry,
		tokens,
	)
	return apiClient
}

func (apiClient *ApiClient) Setup(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
	tokens []string,
) {
	apiClient.client = &http.Client{Timeout: timeout}
	apiClient.SetEndpoint(endpoint)
	apiClient.SetHeaders(headers)
	apiClient.SetMaxRetry(maxRetry)
	apiClient.SetTokens(tokens)
}

func (apiClient *ApiClient) SetTokens(tokens []string) {
	apiClient.tokens = tokens
}

// Getting tokens allows the plugin to configure its worker scheduler according to how many tokens they have.
func (apiClient *ApiClient) GetTokens() []string {
	return apiClient.tokens
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
	// GitHub is currently the only plugin using the multitoken feature, so it's safe to
	// assume setting the Authorization header using Bearer will work. However, this
	// may not be the case for other APIs, and we will need to pass in "authType" or some
	// other field to specify how authorization should be handled.
	if len(apiClient.tokens) > 0 && apiClient.tokens[0] != "" {
		// override the current auth with a new auth
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiClient.tokens[tokenIndex]))
		tokenIndex = (tokenIndex + 1) % len(apiClient.tokens)
	}
	// before send
	if apiClient.beforeRequest != nil {
		err = apiClient.beforeRequest(req)
		if err != nil {
			return nil, err
		}
	}

	var res *http.Response
	retry := 0
	for {
		logger.Print(fmt.Sprintf("[api-client][retry %v] %v %v", retry, method, uri))
		res, err = apiClient.client.Do(req)
		if err != nil {
			if retry < apiClient.maxRetry-1 {
				retry += 1
				continue
			}
			return nil, err
		} else {
			break
		}
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

func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Print(fmt.Sprintf("UnmarshalResponse failed: %v\n%v\n\n", res.Request.URL.String(), string(resBody)))
		return err
	}
	return json.Unmarshal(resBody, &v)
}
