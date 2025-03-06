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

package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ArgoTaskData struct {
	Options   ArgoOptions
	ApiClient *helper.ApiAsyncClient
}

type ArgoOptions struct {
	ConnectionId uint64
	Project      string
}

type ArgoApiParams struct {
	ConnectionId uint64
	Project      string
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
	totalPages := res.Header.Get("X-Total-Pages")
	if totalPages == "" {
		return 1, nil
	}
	pages, err := strconv.Atoi(totalPages)
	if err != nil {
		return 0, errors.Convert(err)
	}
	return pages, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	rawMessages := []json.RawMessage{}

	if res == nil {
		return nil, errors.Default.New("response is nil")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error reading response from %s", res.Request.URL.String()))
	}

	err = errors.Convert(json.Unmarshal(resBody, &rawMessages))
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error decoding response from %s. raw response was: %s", res.Request.URL.String(), string(resBody)))
	}

	return rawMessages, nil
}
