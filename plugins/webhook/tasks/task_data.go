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
	"github.com/mitchellh/mapstructure"
)

type WebhookApiParams struct {
}

type WebhookOptions struct {
	ConnectionId uint64   `json:"connectionId"`
	Tasks        []string `json:"tasks,omitempty"`
}

type WebhookTaskData struct {
	Options *WebhookOptions
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*WebhookOptions, error) {
	var op WebhookOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}
	return &op, nil
}
