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
	"github.com/apache/incubator-devlake/errors"
	"time"

	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

type BitbucketOptions struct {
	ConnectionId               uint64   `json:"connectionId"`
	Tasks                      []string `json:"tasks,omitempty"`
	Since                      string
	Owner                      string
	Repo                       string
	models.TransformationRules `mapstructure:"transformationRules" json:"transformationRules"`
}

type BitbucketTaskData struct {
	Options   *BitbucketOptions
	ApiClient *helper.ApiAsyncClient
	Since     *time.Time
	Repo      *models.BitbucketRepo
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*BitbucketOptions, error) {
	var op BitbucketOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.Owner == "" {
		return nil, errors.Default.New("owner is required for Bitbucket execution")
	}
	if op.Repo == "" {
		return nil, errors.Default.New("repo is required for Bitbucket execution")
	}

	// find the needed Bitbucket now
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}
	return &op, nil
}
