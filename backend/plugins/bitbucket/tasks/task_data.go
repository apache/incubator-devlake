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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

type BitbucketOptions struct {
	ConnectionId                        uint64   `json:"connectionId" mapstructure:"connectionId,omitempty"`
	Tasks                               []string `json:"tasks,omitempty" mapstructure:",omitempty"`
	FullName                            string   `json:"fullName" mapstructure:"fullName"`
	TimeAfter                           string   `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	TransformationRuleId                uint64   `json:"transformationRuleId" mapstructure:"transformationRuleId,omitempty"`
	*models.BitbucketTransformationRule `mapstructure:"transformationRules,omitempty" json:"transformationRules"`
}

type BitbucketTaskData struct {
	Options       *BitbucketOptions
	ApiClient     *api.ApiAsyncClient
	TimeAfter     *time.Time
	RegexEnricher *api.RegexEnricher
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*BitbucketOptions, errors.Error) {
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

func DecodeTaskOptions(options map[string]interface{}) (*BitbucketOptions, errors.Error) {
	var op BitbucketOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func EncodeTaskOptions(op *BitbucketOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := api.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidateTaskOptions(op *BitbucketOptions) errors.Error {
	if op.FullName == "" {
		return errors.BadInput.New("no enough info for Bitbucket execution")
	}
	// find the needed Bitbucket now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
