/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper/common"
	"github.com/apache/incubator-devlake/utils"
)

// ErrIgnoreAndContinue is a error which should be ignored
var ErrIgnoreAndContinue = errors.Default.New("ignore and continue")

// ApiClient is designed for simple api requests
type ApiClient struct {
	client        *http.Client
	endpoint      string
	headers       map[string]string
	beforeRequest common.ApiClientBeforeRequest
	afterResponse common.ApiClientAfterResponse
	ctx           context.Context
	logger        core.Logger
}

// NewApiClient FIXME ...
func NewApiClient(
	ctx context.Context,
	endpoint string,
	headers map[string]string,
	timeout time.Duration,
	proxy string,
	br core.BasicRes,
) (*ApiClient, error) {
	parsedUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, fmt.Sprintf("Invalid URL: %s", endpoint), errors.AsUserMessage())
	}
	if parsedUrl.Scheme == "" {
		return nil, errors.BadInput.New("Invalid URL scheme", errors.AsUserMessage())
	}
	err = utils.CheckDNS(parsedUrl.Hostname())
	if err != nil {
		return nil, errors.Default.Wrap(err, "Failed to resolve DNS", errors.AsUserMessage())
	}
	port, err := utils.ResolvePort(parsedUrl.Port(), parsedUrl.Scheme)
	if err != nil {
		return nil, errors.Default.New("Failed to resolve Port", errors.AsUserMessage())
	}
	err = utils.CheckNetwork(parsedUrl.Hostname(), port, time.Duration(2)*time.Second)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Failed to connect", errors.AsUserMessage())
	}
	apiClient := &ApiClient{}
	apiClient.Setup(
		endpoint,
		headers,
		timeout,
	)
	// create the Transport
	apiClient.client.Transport = &http.Transport{}

	// set insecureSkipVerify
	insecureSkipVerify, err := utils.StrToBoolOr(br.GetConfig("IN_SECURE_SKIP_VERIFY"), false)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to parse IN_SECURE_SKIP_VERIFY", errors.UserMessage("Malformed config on server"))
	}
	if insecureSkipVerify {
		apiClient.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if proxy != "" {
		err = apiClient.SetProxy(proxy)
		if err != nil {
			return nil, err
		}
	}
	apiClient.SetContext(ctx)
	return apiClient, nil
}

// Setup FIXME ...
func (apiClient *ApiClient) Setup(
	endpoint string,
	headers map[string]string,
	timeout time.Duration,

) {
	apiClient.client = &http.Client{Timeout: timeout}
	apiClient.SetEndpoint(endpoint)
	apiClient.SetHeaders(headers)
}

// SetEndpoint FIXME ...
func (apiClient *ApiClient) SetEndpoint(endpoint string) {
	apiClient.endpoint = endpoint
}

// GetEndpoint FIXME ...
func (apiClient *ApiClient) GetEndpoint() string {
	return apiClient.endpoint
}

// SetTimeout FIXME ...
func (apiClient *ApiClient) SetTimeout(timeout time.Duration) {
	apiClient.client.Timeout = timeout
}

// GetTimeout FIXME ...
func (apiClient *ApiClient) GetTimeout() time.Duration {
	return apiClient.client.Timeout
}

// SetHeaders FIXME ...
func (apiClient *ApiClient) SetHeaders(headers map[string]string) {
	apiClient.headers = headers
}

// GetHeaders FIXME ...
func (apiClient *ApiClient) GetHeaders() map[string]string {
	return apiClient.headers
}

// GetBeforeFunction return beforeResponseFunction
func (apiClient *ApiClient) GetBeforeFunction() common.ApiClientBeforeRequest {
	return apiClient.beforeRequest
}

// SetBeforeFunction will set beforeResponseFunction
func (apiClient *ApiClient) SetBeforeFunction(callback common.ApiClientBeforeRequest) {
	apiClient.beforeRequest = callback
}

// GetAfterFunction return afterResponseFunction
func (apiClient *ApiClient) GetAfterFunction() common.ApiClientAfterResponse {
	return apiClient.afterResponse
}

// SetAfterFunction will set afterResponseFunction
func (apiClient *ApiClient) SetAfterFunction(callback common.ApiClientAfterResponse) {
	apiClient.afterResponse = callback
}

// SetContext FIXME ...
func (apiClient *ApiClient) SetContext(ctx context.Context) {
	apiClient.ctx = ctx
}

// SetProxy FIXME ...
func (apiClient *ApiClient) SetProxy(proxyUrl string) error {
	pu, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	if pu.Scheme == "http" || pu.Scheme == "socks5" {
		apiClient.client.Transport.(*http.Transport).Proxy = http.ProxyURL(pu)
	}
	return nil
}

// SetLogger FIXME ...
func (apiClient *ApiClient) SetLogger(logger core.Logger) {
	apiClient.logger = logger
}

func (apiClient *ApiClient) logDebug(format string, a ...interface{}) {
	if apiClient.logger != nil {
		apiClient.logger.Debug(format, a...)
	}
}

func (apiClient *ApiClient) logError(err error, format string, a ...interface{}) {
	if apiClient.logger != nil {
		apiClient.logger.Error(err, format, a...)
	}
}

// Do FIXME ...
func (apiClient *ApiClient) Do(
	method string,
	path string,
	query url.Values,
	body interface{},
	headers http.Header,
) (*http.Response, error) {
	uri, err := GetURIStringPointer(apiClient.endpoint, path, query)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("Unable to construct URI from %s, %s, %s", apiClient.endpoint, path, query), errors.UserMessage("Unable to construct URI"))
	}
	// process body
	var reqBody io.Reader
	if body != nil {
		reqJson, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("unable to serialize API request body for %s", *uri), errors.UserMessage("unable to serialize API request body"))
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
		return nil, errors.Default.Wrap(err, fmt.Sprintf("unable to create API request for %s", *uri), errors.UserMessage("unable to create API request"))
	}
	req.Header.Set("Content-Type", "application/json")

	// populate headers
	if apiClient.headers != nil {
		for name, value := range apiClient.headers {
			req.Header.Set(name, value)
		}
	}
	for name, values := range headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	var res *http.Response
	// before send
	if apiClient.beforeRequest != nil {
		err = apiClient.beforeRequest(req)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error running beforeRequest for %s", req.URL.String()), errors.UserMessage("error before making API call"))
		}
	}
	apiClient.logDebug("[api-client] %v %v", method, *uri)
	res, err = apiClient.client.Do(req)
	if err != nil {
		apiClient.logError(err, "[api-client] failed to request %s with error", req.URL.String())
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error running beforeRequest for %s", req.URL.String()), errors.UserMessage("error before making API call"))
	}
	// after receive
	if apiClient.afterResponse != nil {
		err = apiClient.afterResponse(res)
		if err == ErrIgnoreAndContinue {
			res.Body.Close()
			return res, err
		}
		if err != nil {
			res.Body.Close()
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error running afterRequest for %s", req.URL.String()), errors.UserMessage("error after making API call"))
		}
	}
	return res, nil
}

// Get FIXME ...
func (apiClient *ApiClient) Get(
	path string,
	query url.Values,
	headers http.Header,
) (*http.Response, error) {
	return apiClient.Do(http.MethodGet, path, query, nil, headers)
}

// Post FIXME ...
func (apiClient *ApiClient) Post(
	path string,
	query url.Values,
	body interface{},
	headers http.Header,
) (*http.Response, error) {
	return apiClient.Do(http.MethodPost, path, query, body, headers)
}

// UnmarshalResponse FIXME ...
func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error reading response from %s", res.Request.URL.String()), errors.AsUserMessage())
	}
	err = json.Unmarshal(resBody, &v)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error decoding response from %s: raw response: %s", res.Request.URL.String(), string(resBody)), errors.AsUserMessage())
	}
	return nil
}

// GetURIStringPointer FIXME ...
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

// AddMissingSlashToURL FIXME ...
func AddMissingSlashToURL(baseUrl *string) {
	pattern := `\/$`
	isMatch, _ := regexp.Match(pattern, []byte(*baseUrl))
	if !isMatch {
		*baseUrl += "/"
	}
}

// RemoveStartingSlashFromPath FIXME ...
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
