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

	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type GithubOptions struct {
	ConnectionId                     uint64   `json:"connectionId"`
	TransformationRuleId             uint64   `json:"transformationRuleId"`
	Tasks                            []string `json:"tasks,omitempty"`
	Since                            string
	Owner                            string
	Repo                             string
	*models.GithubTransformationRule `mapstructure:"transformationRules" json:"transformationRules"`
}

type GithubTaskData struct {
	Options       *GithubOptions
	ApiClient     *helper.ApiAsyncClient
	GraphqlClient *helper.GraphqlAsyncClient
	Since         *time.Time
	Repo          *models.GithubRepo
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*GithubOptions, errors.Error) {
	var op GithubOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.Owner == "" {
		return nil, errors.BadInput.New("owner is required for GitHub execution")
	}
	if op.Repo == "" {
		return nil, errors.BadInput.New("repo is required for GitHub execution")
	}
	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	if op.GithubTransformationRule == nil && op.TransformationRuleId == 0 {
		op.GithubTransformationRule = new(models.GithubTransformationRule)
	}
	return &op, nil
}
