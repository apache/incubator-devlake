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
	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

type AzureApiParams struct {
	ConnectionId uint64
	Project      string
}

type AzureOptions struct {
	ConnectionId uint64 `json:"connectionId"`
	Project      string
	Since        string
	Tasks        []string `json:"tasks,omitempty"`
}

type AzureTaskData struct {
	Options    *AzureOptions
	ApiClient  *helper.ApiAsyncClient
	Connection *models.AzureConnection
	Repo       *models.AzureRepo
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*AzureOptions, error) {
	var op AzureOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	// find the needed Azure now
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid", errors.AsUserMessage())
	}
	return &op, nil
}
