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
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
)

type OpsgenieOptions struct {
	ConnectionId uint64   `json:"connectionId"`
	ServiceId    string   `json:"service_id,omitempty"`
	ServiceName  string   `json:"service_name,omitempty"`
	Tasks        []string `json:"tasks,omitempty"`
	*models.OpsenieScopeConfig
}

type OpsgenieTaskData struct {
	Options *OpsgenieOptions
	Client  api.RateLimitedApiClient
}

func (p *OpsgenieOptions) GetParams() any {
	return models.OpsgenieParams{
		ConnectionId: p.ConnectionId,
		ScopeId:      p.ServiceId,
	}
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*OpsgenieOptions, errors.Error) {
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

func DecodeTaskOptions(options map[string]interface{}) (*OpsgenieOptions, errors.Error) {
	var op OpsgenieOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func EncodeTaskOptions(op *OpsgenieOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := api.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidateTaskOptions(op *OpsgenieOptions) errors.Error {
	if op.ServiceName == "" {
		return errors.BadInput.New("not enough info for Opsgenie execution")
	}
	if op.ServiceId == "" {
		return errors.BadInput.New("not enough info for Opsgenie execution")
	}
	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
