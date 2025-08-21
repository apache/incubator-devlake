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

package api

import (
	gocontext "context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
)

// DsRemoteApiProxyHelper is a helper to proxy api request to remote servers
type DsRemoteApiProxyHelper[C plugin.ToolLayerApiConnection] struct {
	*ModelApiHelper[C]
	logger               log.Logger
	httpClientCache      map[string]*ApiClient
	httpClientCacheMutex *sync.RWMutex
}

// NewDsRemoteApiProxyHelper creates a new DsRemoteApiProxyHelper
func NewDsRemoteApiProxyHelper[C plugin.ToolLayerApiConnection](modelApiHelper *ModelApiHelper[C]) *DsRemoteApiProxyHelper[C] {
	return &DsRemoteApiProxyHelper[C]{
		ModelApiHelper:       modelApiHelper,
		logger:               modelApiHelper.basicRes.GetLogger().Nested("remote_api_helper"),
		httpClientCache:      make(map[string]*ApiClient),
		httpClientCacheMutex: &sync.RWMutex{},
	}
}

func (rap *DsRemoteApiProxyHelper[C]) prepare(input *plugin.ApiResourceInput) (*C, *ApiClient, errors.Error) {
	connection, err := rap.FindByPk(input)
	if err != nil {
		return nil, nil, err
	}
	apiClient, err := rap.getApiClient(connection)
	if err != nil {
		return nil, nil, err
	}
	return connection, apiClient, nil
}

func (rap *DsRemoteApiProxyHelper[C]) getApiClient(connection *C) (*ApiClient, errors.Error) {
	c := interface{}(connection)
	key := ""
	if cacheableConn, ok := c.(plugin.CacheableConnection); ok {
		key = cacheableConn.GetHash()
	}
	// try to reuse api client
	if key != "" {
		rap.httpClientCacheMutex.RLock()
		client, ok := rap.httpClientCache[key]
		rap.httpClientCacheMutex.RUnlock()
		if ok {
			rap.logger.Info("Reused api client")
			return client, nil
		}
	}
	// create new client if cache missed
	client, err := NewApiClientFromConnection(gocontext.TODO(), rap.basicRes, c.(plugin.ApiConnection))
	if err != nil {
		return nil, err
	}
	// cache the client if key is not empty
	if key != "" {
		rap.httpClientCacheMutex.Lock()
		rap.httpClientCache[key] = client
		rap.httpClientCacheMutex.Unlock()
	} else {
		rap.logger.Info("No api client reuse")
	}
	return client, nil
}

// Proxy forwards api request to a specific remote server
func (rap *DsRemoteApiProxyHelper[C]) Proxy(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	_, apiClient, err := rap.prepare(input)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.Get(input.Params["path"], input.Query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := errors.Convert01(io.ReadAll(resp.Body))
	if err != nil {
		return nil, err
	}
	// verify response body is json
	var tmp interface{}
	err = errors.Convert(json.Unmarshal(body, &tmp))
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	for k, vs := range resp.Header {
		// skip headers doesn't start with "x-"
		if !strings.HasPrefix(strings.ToLower(k), "x-") {
			continue
		}
		for _, v := range vs {
			headers.Add(k, v)
		}
	}
	return &plugin.ApiResourceOutput{Status: resp.StatusCode, Body: json.RawMessage(body), Header: headers}, nil
}
