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
	ConnectionId                     uint64   `json:"connectionId" mapstructure:"connectionId,omitempty"`
	TransformationRuleId             uint64   `json:"transformationRuleId" mapstructure:"transformationRuleId,omitempty"`
	ScopeId                          string   `json:"scopeId" mapstructure:"scopeId,omitempty"`
	GithubId                         int      `json:"githubId" mapstructure:"githubId,omitempty"`
	Tasks                            []string `json:"tasks,omitempty" mapstructure:",omitempty"`
	CreatedDateAfter                 string   `json:"createdDateAfter" mapstructure:"createdDateAfter,omitempty"`
	Owner                            string   `json:"owner" mapstructure:"owner,omitempty"`
	Repo                             string   `json:"repo"  mapstructure:"repo,omitempty"`
	*models.GithubTransformationRule `mapstructure:"transformationRules,omitempty" json:"transformationRules"`
}

type GithubTaskData struct {
	Options          *GithubOptions
	ApiClient        *helper.ApiAsyncClient
	GraphqlClient    *helper.GraphqlAsyncClient
	CreatedDateAfter *time.Time
}

type GithubApiParams struct {
	ConnectionId uint64
	Owner        string
	Repo         string
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*GithubOptions, errors.Error) {
	op, err := DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	err = ValidateTaskOptions(op)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func DecodeTaskOptions(options map[string]interface{}) (*GithubOptions, errors.Error) {
	var op GithubOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func EncodeTaskOptions(op *GithubOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := helper.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidateTaskOptions(op *GithubOptions) errors.Error {
	if op.Owner == "" {
		return errors.BadInput.New("owner is required for GitHub execution")
	}
	if op.Repo == "" {
		return errors.BadInput.New("repo is required for GitHub execution")
	}
	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	if op.GithubTransformationRule == nil && op.TransformationRuleId == 0 {
		op.GithubTransformationRule = new(models.GithubTransformationRule)
	}
	return nil
}
