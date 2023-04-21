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
	ConnectionId                     uint64   `mapstructure:"connectionId" json:"connectionId"`
	ProjectId                        int      `mapstructure:"projectId" json:"projectId"`
	TransformationRuleId             uint64   `mapstructure:"transformationRuleId" json:"transformationRuleId"`
	Tasks                            []string `mapstructure:"tasks" json:"tasks,omitempty"`
	TimeAfter                        string
	*models.GitlabTransformationRule `mapstructure:"transformationRules" json:"transformationRules"`
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
	if op.GitlabTransformationRule == nil && op.TransformationRuleId == 0 {
		op.GitlabTransformationRule = new(models.GitlabTransformationRule)
	}
	return &op, nil
}
