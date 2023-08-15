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
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type JiraOptions struct {
	ConnectionId  uint64 `json:"connectionId"`
	BoardId       uint64 `json:"boardId"`
	TimeAfter     string
	ScopeConfig   *models.JiraScopeConfig `json:"scopeConfig"`
	ScopeId       string
	ScopeConfigId uint64
	PageSize      int
}

type JiraTaskData struct {
	Options        *JiraOptions
	ApiClient      *api.ApiAsyncClient
	TimeAfter      *time.Time
	JiraServerInfo models.JiraServerInfo
}

type JiraApiParams models.JiraApiParams

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*JiraOptions, errors.Error) {
	var op JiraOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid connectionId:%d", op.ConnectionId))
	}
	if op.BoardId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid boardId:%d", op.BoardId))
	}
	return &op, nil
}
