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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

type JenkinsApiParams struct {
	ConnectionId uint64
	JobName      string
	JobPath      string
}

type JenkinsOptions struct {
	ConnectionId                uint64 `json:"connectionId"`
	TransformationRuleId        uint64 `json:"transformationRuleId"`
	JobFullName                 string `json:"JobFullName"`
	JobName                     string `json:"jobName"`
	JobPath                     string `json:"jobPath"`
	Since                       string
	Tasks                       []string `json:"tasks,omitempty"`
	*models.TransformationRules `mapstructure:"transformationRules" json:"transformationRules"`
}

type JenkinsTaskData struct {
	Options    *JenkinsOptions
	ApiClient  *helper.ApiAsyncClient
	Connection *models.JenkinsConnection
	Since      *time.Time
	Job        *models.JenkinsJob
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*JenkinsOptions, errors.Error) {
	var op JenkinsOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	if op.TransformationRules == nil && op.TransformationRuleId == 0 {
		op.TransformationRules = new(models.TransformationRules)
	}
	return &op, nil
}
