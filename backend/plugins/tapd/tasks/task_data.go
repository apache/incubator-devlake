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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

type TapdOptions struct {
	ConnectionId  uint64                  `json:"connectionId" mapstructure:"connectionId,omitempty"`
	WorkspaceId   uint64                  `json:"workspaceId" mapstructure:"workspaceId,omitempty"`
	PageSize      uint64                  `json:"pageSize,omitempty" mapstructure:"pageSize,omitempty"`
	CstZone       *time.Location          `json:"cstZone,omitempty" mapstructure:"cstZone,omitempty"`
	ScopeConfigId uint64                  `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
	ScopeConfig   *models.TapdScopeConfig `json:"scopeConfig,omitempty" mapstructure:"pageSize,omitempty"`
}

type TapdTaskData struct {
	Options    *TapdOptions
	ApiClient  *helper.ApiAsyncClient
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
