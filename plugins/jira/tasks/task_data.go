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
	"github.com/apache/incubator-devlake/errors"
	"time"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

type StatusMapping struct {
	StandardStatus string `json:"standardStatus"`
}

type StatusMappings map[string]StatusMapping

type TypeMapping struct {
	StandardType   string         `json:"standardType"`
	StatusMappings StatusMappings `json:"statusMappings"`
}

type TypeMappings map[string]TypeMapping

type TransformationRules struct {
	EpicKeyField               string       `json:"epicKeyField"`
	StoryPointField            string       `json:"storyPointField"`
	RemotelinkCommitShaPattern string       `json:"remotelinkCommitShaPattern"`
	TypeMappings               TypeMappings `json:"typeMappings"`
}

type JiraOptions struct {
	ConnectionId        uint64 `json:"connectionId"`
	BoardId             uint64 `json:"boardId"`
	Since               string
	TransformationRules TransformationRules `json:"transformationRules"`
}

type JiraTaskData struct {
	Options        *JiraOptions
	ApiClient      *helper.ApiAsyncClient
	Since          *time.Time
	JiraServerInfo models.JiraServerInfo
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*JiraOptions, error) {
	var op JiraOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid connectionId:%d", op.ConnectionId), errors.AsUserMessage())
	}
	if op.BoardId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid boardId:%d", op.BoardId), errors.AsUserMessage())
	}
	return &op, nil
}
