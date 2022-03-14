package core

import (
	"bytes"
	"context"
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
	ctx           context.Context
	scheduler     *utils.WorkerScheduler
	logger        Logger
}

func NewApiClient(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
	scheduler *utils.WorkerScheduler,
) *ApiClient {
	apiClient := &ApiClient{}
	apiClient.Setup(
		endpoint,
		headers,
		timeout,
		maxRetry,
		scheduler,
	)
	return apiClient
}

func (apiClient *ApiClient) Setup(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	maxRetry int,
	scheduler *utils.WorkerScheduler,

) {
	apiClient.client = &http.Client{Timeout: timeout}
	apiClient.SetEndpoint(endpoint)
	apiClient.SetHeaders(headers)
	apiClient.SetMaxRetry(maxRetry)
	apiClient.setScheduler(scheduler)
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

func (apiClient *ApiClient) SetAfterFunction(callback ApiClientAfterResponse) {
	apiClient.afterReponse = callback
}

func (apiClient *ApiClient) SetContext(ctx context.Context) {
	apiClient.ctx = ctx
}

func (apiClient *ApiClient) setScheduler(scheduler *utils.WorkerScheduler) {
	apiClient.scheduler = scheduler
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

func (apiClient *ApiClient) SetLogger(logger Logger) {
	apiClient.logger = logger
}

func (apiClient *ApiClient) logDebug(format string, a ...interface{}) {
	if apiClient.logger != nil {
		apiClient.logger.Debug(format, a...)
	}
}

func (apiClient *ApiClient) logError(format string, a ...interface{}) {
	if apiClient.logger != nil {
		apiClient.logger.Error(format, a...)
	}
}

func (apiClient *ApiClient) Do(
	method string,
	path string,
	query url.Values,
	body interface{},
	headers http.Header,
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
	var req *http.Request
	if apiClient.ctx != nil {
		req, err = http.NewRequestWithContext(apiClient.ctx, method, *uri, reqBody)
	} else {
		req, err = http.NewRequest(method, *uri, reqBody)
	}
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
		for name, values := range headers {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
	}

	var res *http.Response
	// canceling check
	if apiClient.ctx != nil {
		select {
		case <-apiClient.ctx.Done():
			return nil, apiClient.ctx.Err()
		default:
		}
	}
	// before send
	if apiClient.beforeRequest != nil {
		err = apiClient.beforeRequest(req)
		if err != nil {
			return nil, err
		}
	}
	apiClient.logDebug("[api-client] %v %v", method, *uri)
	res, err = apiClient.client.Do(req)
	if err != nil {
		apiClient.logError("[api-client] failed to request %s with error:\n%w", req.URL.String(), err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// canceling check
	if apiClient.ctx != nil {
		select {
		case <-apiClient.ctx.Done():
			return nil, apiClient.ctx.Err()
		default:
		}
	}

	// after recieve
	if apiClient.afterReponse != nil {
		err = apiClient.afterReponse(res)
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

func (apiClient *ApiClient) Get(
	path string,
	query url.Values,
	headers http.Header,
) (*http.Response, error) {
	return apiClient.Do("GET", path, query, nil, headers)
}

func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%w %s", err, res.Request.URL.String())
	}
	err = json.Unmarshal(resBody, &v)
	if err != nil {
		return fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
	}
	return nil
}

func GetURIStringPointer(baseUrl string, relativePath string, query url.Values) (*string, error) {
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
	if query != nil {
		queryString := u.Query()
		for key, value := range query {
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

func (apiClient *ApiClient) DoAsync(
	method string,
	path string,
	query url.Values,
	body interface{},
	header http.Header,
	handler func(*http.Response) error,
	retry int,
) error {
	return apiClient.scheduler.Submit(func() error {
		var err error
		var res *http.Response
		res, err = apiClient.Do(method, path, query, body, header)
		if err == nil {
			err = handler(res)
		}
		if err != nil && retry < apiClient.maxRetry {
			apiClient.logError("retry #%d for %s", retry, err.Error())
			err = apiClient.DoAsync(method, path, query, body, header, handler, retry+1)
		}
		return err
	})
}

func (apiClient *ApiClient) GetAsync(path string, query url.Values, header http.Header, handler func(*http.Response) error) error {
	return apiClient.DoAsync(http.MethodGet, path, query, nil, header, handler, 0)
}

func (apiClient *ApiClient) GetSync(path string, query url.Values, header http.Header, handler func(*http.Response) (bool, error)) (bool, error) {
	var err error
	var res *http.Response
	res, err = apiClient.Do(http.MethodGet, path, query, nil, header)
	if err != nil {
		return false, err
	}
	return handler(res)
}

func (apiClient *ApiClient) WaitAsync() {
	println("wait until finish")
	apiClient.scheduler.WaitUntilFinish()
}

type AsyncApiClient interface {
	GetAsync(path string, query url.Values, header http.Header, handler func(*http.Response) error) error
	WaitAsync()
}

type SyncApiClient interface {
	GetSync(path string, query url.Values, header http.Header, handler func(*http.Response) (bool, error)) (bool, error)
}

var _ AsyncApiClient = (*ApiClient)(nil)
