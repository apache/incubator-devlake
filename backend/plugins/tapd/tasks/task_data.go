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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

type TapdOptions struct {
	ConnectionId         uint64 `mapstruct:"connectionId"`
	WorkspaceId          uint64 `mapstruct:"workspaceId"`
	PageSize             uint64 `mapstruct:"pageSize"`
	TimeAfter            string `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	CstZone              *time.Location
	TransformationRuleId uint64
	TransformationRules  *TransformationRules `json:"transformationRules"`
}

func MakeTransformationRules(rule models.TapdTransformationRule) (*TransformationRules, errors.Error) {
	var statusMapping StatusMappings
	var typeMapping TypeMappings
	var err error
	if len(rule.TypeMappings) > 0 {
		err = json.Unmarshal(rule.TypeMappings, &typeMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the typeMapping")
		}
	}
	if len(rule.StatusMappings) > 0 {
		err = json.Unmarshal(rule.StatusMappings, &statusMapping)
		if err != nil {
			return nil, errors.Default.Wrap(err, "unable to unmarshal the statusMapping")
		}
	}
	result := &TransformationRules{
		StatusMappings: statusMapping,
		TypeMappings:   typeMapping,
	}
	return result, nil
}

type TapdTaskData struct {
	Options    *TapdOptions
	ApiClient  *helper.ApiAsyncClient
	TimeAfter  *time.Time
	Connection *models.TapdConnection
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*TapdOptions, errors.Error) {
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

func DecodeTaskOptions(options map[string]interface{}) (*TapdOptions, errors.Error) {
	var op TapdOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func ValidateTaskOptions(op *TapdOptions) errors.Error {
	if op.WorkspaceId == 0 {
		return errors.BadInput.New("no enough info for tapd execution")
	}
	// find the needed tapd now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}

type StatusMappings map[string]string

type TypeMappings map[string]string

type TransformationRules struct {
	TypeMappings   TypeMappings   `json:"typeMappings"`
	StatusMappings StatusMappings `json:"statusMappings"`
}
