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

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/mitchellh/mapstructure"
)

type JenkinsApiParams struct {
	ConnectionId uint64
}

type JenkinsOptions struct {
	ConnectionId uint64 `json:"connectionId"`
	Since        string
	Tasks        []string `json:"tasks,omitempty"`
}

type JenkinsTaskData struct {
	Options    *JenkinsOptions
	ApiClient  *helper.ApiAsyncClient
	Connection *models.JenkinsConnection
	Since      *time.Time
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*JenkinsOptions, error) {
	var op JenkinsOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters", errors.AsUserMessage())
	}
	// find the needed Jenkins now
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	return &op, nil
}
