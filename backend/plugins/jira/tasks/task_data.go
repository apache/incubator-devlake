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
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
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

type JiraTransformationRule struct {
	ConnectionId               uint64       `mapstructure:"connectionId" json:"connectionId"`
	Name                       string       `gorm:"type:varchar(255)" validate:"required"`
	EpicKeyField               string       `json:"epicKeyField"`
	StoryPointField            string       `json:"storyPointField"`
	RemotelinkCommitShaPattern string       `json:"remotelinkCommitShaPattern"`
	RemotelinkRepoPattern      []string     `json:"remotelinkRepoPattern"`
	TypeMappings               TypeMappings `json:"typeMappings"`
}

func (r *JiraTransformationRule) ToDb() (*models.JiraTransformationRule, errors.Error) {
	blob, err := json.Marshal(r.TypeMappings)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error marshaling TypeMappings")
	}
	remotelinkRepoPattern, err := json.Marshal(r.RemotelinkRepoPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error marshaling RemotelinkRepoPattern")
	}
	rule := &models.JiraTransformationRule{
		ConnectionId:               r.ConnectionId,
		Name:                       r.Name,
		EpicKeyField:               r.EpicKeyField,
		StoryPointField:            r.StoryPointField,
		RemotelinkCommitShaPattern: r.RemotelinkCommitShaPattern,
		RemotelinkRepoPattern:      remotelinkRepoPattern,
		TypeMappings:               blob,
	}
	if err1 := rule.VerifyRegexp(); err1 != nil {
		return nil, err1
	}
	return rule, nil
}

func MakeTransformationRules(rule models.JiraTransformationRule) (*JiraTransformationRule, errors.Error) {
	var typeMapping TypeMappings
	var err error
	if len(rule.TypeMappings) > 0 {
		err = json.Unmarshal(rule.TypeMappings, &typeMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the typeMapping")
		}
	}
	var remotelinkRepoPattern []string
	if len(rule.RemotelinkRepoPattern) > 0 {
		err = json.Unmarshal(rule.RemotelinkRepoPattern, &remotelinkRepoPattern)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error unMarshaling RemotelinkRepoPattern")
		}
	}
	result := &JiraTransformationRule{
		ConnectionId:               rule.ConnectionId,
		Name:                       rule.Name,
		EpicKeyField:               rule.EpicKeyField,
		StoryPointField:            rule.StoryPointField,
		RemotelinkCommitShaPattern: rule.RemotelinkCommitShaPattern,
		RemotelinkRepoPattern:      remotelinkRepoPattern,
		TypeMappings:               typeMapping,
	}
	return result, nil
}

type JiraOptions struct {
	ConnectionId         uint64 `json:"connectionId"`
	BoardId              uint64 `json:"boardId"`
	TimeAfter            string
	TransformationRules  *JiraTransformationRule `json:"transformationRules"`
	ScopeId              string
	TransformationRuleId uint64
	PageSize             int
}

type JiraTaskData struct {
	Options        *JiraOptions
	ApiClient      *api.ApiAsyncClient
	TimeAfter      *time.Time
	JiraServerInfo models.JiraServerInfo
}

type JiraApiParams struct {
	ConnectionId uint64
	BoardId      uint64
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*JiraOptions, errors.Error) {
	var op JiraOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid connectionId:%d", op.ConnectionId))
	}
	if op.BoardId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid boardId:%d", op.BoardId))
	}
	return &op, nil
}
