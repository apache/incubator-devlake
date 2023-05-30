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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

type GitlabOptions struct {
	ConnectionId  uint64                    `mapstructure:"connectionId" json:"connectionId"`
	ProjectId     int                       `mapstructure:"projectId" json:"projectId"`
	ScopeConfigId uint64                    `mapstructure:"scopeConfigId" json:"scopeConfigId"`
	TimeAfter     string                    `mapstructure:"timeAfter" json:"timeAfter"`
	ScopeConfig   *models.GitlabScopeConfig `mapstructure:"scopeConfig" json:"scopeConfig"`
}

type GitlabTaskData struct {
	Options       *GitlabOptions
	ApiClient     *helper.ApiAsyncClient
	ProjectCommit *models.GitlabProjectCommit
	TimeAfter     *time.Time
	RegexEnricher *helper.RegexEnricher
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*GitlabOptions, errors.Error) {
	var op GitlabOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ProjectId == 0 {
		return nil, errors.BadInput.New("ProjectId is required for Gitlab execution")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	if op.ScopeConfig == nil && op.ScopeConfigId == 0 {
		op.ScopeConfig = new(models.GitlabScopeConfig)
	}
	return &op, nil
}
