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

package token

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type refreshApiClient struct {
	endpoint string
	client   *http.Client
	timeout  time.Duration
}

func newRefreshApiClient(endpoint string, client *http.Client) plugin.ApiClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &refreshApiClient{
		endpoint: endpoint,
		client:   client,
		timeout:  10 * time.Second,
	}
}

func (c *refreshApiClient) SetData(name string, data interface{}) {}

func (c *refreshApiClient) GetData(name string) interface{} { return nil }

func (c *refreshApiClient) SetHeaders(headers map[string]string) {}

func (c *refreshApiClient) SetBeforeFunction(callback plugin.ApiClientBeforeRequest) {}

func (c *refreshApiClient) GetBeforeFunction() plugin.ApiClientBeforeRequest { return nil }

func (c *refreshApiClient) SetAfterFunction(callback plugin.ApiClientAfterResponse) {}

func (c *refreshApiClient) GetAfterFunction() plugin.ApiClientAfterResponse { return nil }

func (c *refreshApiClient) Get(path string, query url.Values, headers http.Header) (*http.Response, errors.Error) {
	return c.do(http.MethodGet, path, query, nil, headers)
}

func (c *refreshApiClient) Post(path string, query url.Values, body interface{}, headers http.Header) (*http.Response, errors.Error) {
	return c.do(http.MethodPost, path, query, body, headers)
}

func (c *refreshApiClient) do(method, path string, query url.Values, body interface{}, headers http.Header) (*http.Response, errors.Error) {
	uri, err := api.GetURIStringPointer(c.endpoint, path, query)
	if err != nil {
		return nil, err
	}

	var reqBody *bytes.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Convert(err)
		}
		reqBody = bytes.NewReader(payload)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req, err := errors.Convert01(http.NewRequestWithContext(ctx, method, *uri, reqBody))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for name, values := range headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	res, err := errors.Convert01(c.client.Do(req))
	if err != nil {
		return nil, err
	}
	return res, nil
}
