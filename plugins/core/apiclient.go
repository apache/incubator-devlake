package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/merico-dev/lake/logger"
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
func (apiClient *ApiClient) GetEndpoint() string {
	return apiClient.endpoint
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
func (apiClient *ApiClient) GetHeaders() map[string]string {
	return apiClient.headers
}

func (apiClient *ApiClient) SetBeforeFunction(callback ApiClientBeforeRequest) {
	apiClient.beforeRequest = callback
}

func (apiClient *ApiClient) SetProxy(proxyUrl string) error {
	pu, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	if pu.Scheme == "http" || pu.Scheme == "socks5" {
		apiClient.client.Transport = &http.Transport{Proxy: http.ProxyURL(pu)}
	}
	return nil
}

func (apiClient *ApiClient) Do(
	method string,
	path string,
	query *url.Values,
	body *map[string]interface{},
	headers *map[string]string,
) (*http.Response, error) {
	uri, err := GetURIStringPointer(apiClient.endpoint, path, query)
	if err != nil {
		return nil, err
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
	req, err := http.NewRequest(method, *uri, reqBody)
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

	var res *http.Response
	retry := 0
	for {
		logger.Print(fmt.Sprintf("[api-client][retry %v] %v %v", retry, method, *uri))
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

func GetURIStringPointer(baseUrl string, relativePath string, queryParams *url.Values) (*string, error) {
	// If the base URL doesn't end with a slash, and has a relative path attached
	// the values will be removed by the Go package, therefore we need to add a missing slash.
	AddMissingSlashToURL(&baseUrl)
	base, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	// If the relative path starts with a '/', we need to remove it or the values will be removed by the Go package.
	relativePath = RemoveStartingSlashFromPath(relativePath)
	u, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}
	if queryParams != nil {
		queryString := u.Query()
		for key, value := range *queryParams {
			queryString.Set(key, strings.Join(value, ""))
		}
		u.RawQuery = queryString.Encode()
	}
	uri := base.ResolveReference(u).String()
	return &uri, nil
}

func AddMissingSlashToURL(baseUrl *string) {
	pattern := `\/$`
	isMatch, _ := regexp.Match(pattern, []byte(*baseUrl))
	if !isMatch {
		*baseUrl += "/"
	}
}
func RemoveStartingSlashFromPath(relativePath string) string {
	pattern := `^\/`
	byteArrayOfPath := []byte(relativePath)
	isMatch, _ := regexp.Match(pattern, byteArrayOfPath)
	if isMatch {
		// Remove the slash.
		// This is basically the trimFirstRune function found: https://stackoverflow.com/questions/48798588/how-do-you-remove-the-first-character-of-a-string/48798712
		_, i := utf8.DecodeRuneInString(relativePath)
		return relativePath[i:]
	}
	return relativePath
}
