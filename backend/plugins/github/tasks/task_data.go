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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

type GithubOptions struct {
	ConnectionId  uint64                    `json:"connectionId" mapstructure:"connectionId,omitempty"`
	ScopeConfigId uint64                    `json:"scopeConfigId" mapstructure:"scopeConfigId,omitempty"`
	GithubId      int                       `json:"githubId" mapstructure:"githubId,omitempty"`
	TimeAfter     string                    `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	Owner         string                    `json:"owner" mapstructure:"owner,omitempty"`
	Repo          string                    `json:"repo"  mapstructure:"repo,omitempty"`
	Name          string                    `json:"name"  mapstructure:"name,omitempty"`
	ScopeConfig   *models.GithubScopeConfig `mapstructure:"scopeConfig,omitempty" json:"scopeConfig"`
}

type GithubTaskData struct {
	Options       *GithubOptions
	ApiClient     *helper.ApiAsyncClient
	GraphqlClient *helper.GraphqlAsyncClient
	TimeAfter     *time.Time
	RegexEnricher *helper.RegexEnricher
}

// TODO: avoid touching too many files, should be removed in the future
type GithubApiParams models.GithubApiParams

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
	if op.Name == "" {
		op.Name = strings.Trim(fmt.Sprintf("%s/%s", op.Owner, op.Repo), " ")
	}
	if op.Name == "" && op.GithubId == 0 {
		return errors.BadInput.New("no enough info for GitHub execution")
	}
	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
