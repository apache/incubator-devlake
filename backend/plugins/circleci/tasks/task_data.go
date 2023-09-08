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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
)

type CircleciOptions struct {
	ConnectionId  uint64                      `json:"connectionId"`
	ProjectSlug   string                      `json:"projectSlug"`
	PageSize      uint64                      `mapstruct:"pageSize"`
	TimeAfter     string                      `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	ScopeConfigId uint64                      `json:"scopeConfigId" mapstructure:"scopeConfigId,omitempty"`
	ScopeConfig   *models.CircleciScopeConfig `json:"scopeConfig" mapstructure:"scopeConfig,omitempty"`
}

type CircleciTaskData struct {
	Options   *CircleciOptions
	ApiClient *helper.ApiAsyncClient
	TimeAfter *time.Time
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*CircleciOptions, errors.Error) {
	var op CircleciOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}
	return &op, nil
}
