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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/mitchellh/mapstructure"
)

type ZentaoApiParams models.ZentaoApiParams

type ZentaoOptions struct {
	// options means some custom params required by plugin running.
	// Such As How many rows do your want
	// You can use it in subtasks, and you need to pass it to main.go and pipelines.
	ConnectionId uint64 `json:"connectionId"`
	ProjectId    int64  `json:"projectId" mapstructure:"projectId"`
	// TODO not support now
	TimeAfter     string              `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	ScopeConfigId uint64              `json:"scopeConfigId" mapstructure:"scopeConfigId,omitempty"`
	ScopeConfigs  *ZentaoScopeConfigs `json:"scopeConfigs" mapstructure:"scopeConfigs,omitempty"`
}

func (o *ZentaoOptions) GetParams() any {
	return models.ZentaoApiParams{
		ConnectionId: o.ConnectionId,
		ProjectId:    o.ProjectId,
	}
}

type TypeMappings map[string]string

type StatusMappings map[string]string

type ZentaoScopeConfigs struct {
	TypeMappings        TypeMappings   `json:"typeMappings"`
	BugStatusMappings   StatusMappings `json:"bugStatusMappings"`
	StoryStatusMappings StatusMappings `json:"storyStatusMappings"`
	TaskStatusMappings  StatusMappings `json:"taskStatusMappings"`
}

func MakeScopeConfigs(rule models.ZentaoScopeConfig) (*ZentaoScopeConfigs, errors.Error) {
	var bugStatusMapping StatusMappings
	var storyStatusMapping StatusMappings
	var taskStatusMapping StatusMappings
	var typeMapping TypeMappings
	var err error
	if len(rule.TypeMappings) > 0 {
		err = json.Unmarshal(rule.TypeMappings, &typeMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the typeMapping")
		}
	}
	if len(rule.BugStatusMappings) > 0 {
		err = json.Unmarshal(rule.BugStatusMappings, &bugStatusMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the statusMapping")
		}
	}
	if len(rule.StoryStatusMappings) > 0 {
		err = json.Unmarshal(rule.StoryStatusMappings, &storyStatusMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the statusMapping")
		}
	}
	if len(rule.TaskStatusMappings) > 0 {
		err = json.Unmarshal(rule.TaskStatusMappings, &taskStatusMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the statusMapping")
		}
	}
	result := &ZentaoScopeConfigs{
		TypeMappings:        typeMapping,
		BugStatusMappings:   bugStatusMapping,
		StoryStatusMappings: storyStatusMapping,
		TaskStatusMappings:  taskStatusMapping,
	}
	return result, nil
}

type ZentaoTaskData struct {
	Options  *ZentaoOptions
	RemoteDb dal.Dal

	TimeAfter    *time.Time
	ProjectName  string
	Stories      map[int64]struct{}
	Tasks        map[int64]struct{}
	Bugs         map[int64]struct{}
	AccountCache *AccountCache
	ApiClient    *helper.ApiAsyncClient
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*ZentaoOptions, error) {
	var op ZentaoOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}
	if op.ProjectId == 0 {
		return nil, fmt.Errorf("please set projectId")
	}
	return &op, nil
}

func EncodeTaskOptions(op *ZentaoOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := helper.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
