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

package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type PagingInfo struct {
	Next  int
	Last  int
	First int
	Prev  int
}

type RateLimitInfo struct {
	Date      time.Time
	ResetTime time.Time
	Remaining int
}

type GiteaApiParams struct {
	ConnectionId uint64
	Repo         string
	Owner        string
}

type GiteaInput struct {
	Repo  string
	Owner string
	Iid   int
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	total := res.Header.Get("X-PageCount")
	if total == "" {
		return 0, nil
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		return 0, err
	}
	return totalInt, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, error) {
	var rawMessages []json.RawMessage

	if res == nil {
		return nil, fmt.Errorf("res is nil")
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, res.Request.URL.String())
	}

	err = json.Unmarshal(resBody, &rawMessages)
	if err != nil {
		return nil, fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
	}

	return rawMessages, nil
}

func CreateRawDataSubTaskArgs(taskCtx core.SubTaskContext, Table string) (*helper.RawDataSubTaskArgs, *GiteaTaskData) {
	data := taskCtx.GetData().(*GiteaTaskData)
	RawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GiteaApiParams{
			ConnectionId: data.Options.ConnectionId,
			Repo:         data.Options.Repo,
			Owner:        data.Options.Owner,
		},
		Table: Table,
	}
	return RawDataSubTaskArgs, data
}
