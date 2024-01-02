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

package plugin

import (
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
)

// ApiAsyncCallback FIXME ...
type ApiAsyncCallback func(*http.Response) errors.Error

// ApiClientBeforeRequest FIXME ...
type ApiClientBeforeRequest func(req *http.Request) errors.Error

// ApiClientAfterResponse FIXME ...
type ApiClientAfterResponse func(res *http.Response) errors.Error

// ApiClientAbstract defines the functionalities needed by all plugins for Synchronized API Request
type ApiClient interface {
	SetData(name string, data interface{})
	GetData(name string) interface{}
	SetHeaders(headers map[string]string)
	SetBeforeFunction(callback ApiClientBeforeRequest)
	GetBeforeFunction() ApiClientBeforeRequest
	SetAfterFunction(callback ApiClientAfterResponse)
	GetAfterFunction() ApiClientAfterResponse
	Get(path string, query url.Values, headers http.Header) (*http.Response, errors.Error)
	Post(path string, query url.Values, body interface{}, headers http.Header) (*http.Response, errors.Error)
}
