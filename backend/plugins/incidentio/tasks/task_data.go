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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/incidentio/models"
)

type IncidentioOptions struct {
	ConnectionId     uint64                        `json:"connectionId" mapstructure:"connectionId,omitempty"`
	IncidentTypeId   string                        `json:"incidentTypeId,omitempty" mapstructure:"incidentTypeId,omitempty"`
	IncidentTypeName string                        `json:"incidentTypeName,omitempty" mapstructure:"incidentTypeName,omitempty"`
	ScopeConfigId    uint64                        `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
	ScopeConfig      *models.IncidentioScopeConfig `json:"scopeConfig,omitempty" mapstructure:"scopeConfig,omitempty"`
}

type IncidentioTaskData struct {
	Options *IncidentioOptions
	Client  api.RateLimitedApiClient
}

func (p *IncidentioOptions) GetParams() any {
	scopeId := p.IncidentTypeId
	if scopeId == "" {
		scopeId = "all"
	}
	return models.IncidentioParams{
		ConnectionId: p.ConnectionId,
		ScopeId:      scopeId,
	}
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*IncidentioOptions, errors.Error) {
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

func DecodeTaskOptions(options map[string]interface{}) (*IncidentioOptions, errors.Error) {
	var op IncidentioOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func ValidateTaskOptions(op *IncidentioOptions) errors.Error {
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
