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

type CommitUrlPattern struct {
	Pattern string `json:"pattern"`
	Regex   string `json:"regex"`
}

type TypeMappings map[string]TypeMapping

type JiraScopeConfig struct {
	Entities                   []string           `json:"entities"`
	ConnectionId               uint64             `mapstructure:"connectionId" json:"connectionId"`
	Name                       string             `gorm:"type:varchar(255)" validate:"required"`
	EpicKeyField               string             `json:"epicKeyField"`
	StoryPointField            string             `json:"storyPointField"`
	RemotelinkCommitShaPattern string             `json:"remotelinkCommitShaPattern"`
	RemotelinkRepoPattern      []CommitUrlPattern `json:"remotelinkRepoPattern"`
	TypeMappings               TypeMappings       `json:"typeMappings"`
	ApplicationType            string             `json:"applicationType"`
}

func (r *JiraScopeConfig) ToDb() (*models.JiraScopeConfig, errors.Error) {
	blob, err := json.Marshal(r.TypeMappings)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error marshaling TypeMappings")
	}
	if r.ApplicationType != "" && len(r.RemotelinkRepoPattern) == 0 {
		return nil, errors.Default.New("error remotelinkRepoPattern is empty")
	}
	remotelinkRepoPattern, err := json.Marshal(r.RemotelinkRepoPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error marshaling RemotelinkRepoPattern")
	}
	scopeConfig := &models.JiraScopeConfig{
		ConnectionId:               r.ConnectionId,
		Name:                       r.Name,
		EpicKeyField:               r.EpicKeyField,
		StoryPointField:            r.StoryPointField,
		RemotelinkCommitShaPattern: r.RemotelinkCommitShaPattern,
		RemotelinkRepoPattern:      remotelinkRepoPattern,
		TypeMappings:               blob,
		ApplicationType:            r.ApplicationType,
	}
	scopeConfig.Entities = r.Entities
	if err1 := scopeConfig.VerifyRegexp(); err1 != nil {
		return nil, err1
	}
	return scopeConfig, nil
}

func MakeScopeConfig(rule models.JiraScopeConfig) (*JiraScopeConfig, errors.Error) {
	var typeMapping TypeMappings
	var err error
	if len(rule.TypeMappings) > 0 {
		err = json.Unmarshal(rule.TypeMappings, &typeMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the typeMapping")
		}
	}
	var remotelinkRepoPattern []CommitUrlPattern
	if len(rule.RemotelinkRepoPattern) > 0 {
		err = json.Unmarshal(rule.RemotelinkRepoPattern, &remotelinkRepoPattern)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error unMarshaling RemotelinkRepoPattern")
		}
	}
	result := &JiraScopeConfig{
		Entities:                   rule.Entities,
		ConnectionId:               rule.ConnectionId,
		Name:                       rule.Name,
		EpicKeyField:               rule.EpicKeyField,
		StoryPointField:            rule.StoryPointField,
		RemotelinkCommitShaPattern: rule.RemotelinkCommitShaPattern,
		RemotelinkRepoPattern:      remotelinkRepoPattern,
		TypeMappings:               typeMapping,
		ApplicationType:            rule.ApplicationType,
	}
	return result, nil
}

type JiraOptions struct {
	ConnectionId  uint64 `json:"connectionId"`
	BoardId       uint64 `json:"boardId"`
	TimeAfter     string
	ScopeConfig   *JiraScopeConfig `json:"scopeConfig"`
	ScopeId       string
	ScopeConfigId uint64
	PageSize      int
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
