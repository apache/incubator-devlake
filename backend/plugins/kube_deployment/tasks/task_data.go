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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type KubeDeploymentApiParams struct {
	ConnectionId uint64 `json:"connectionId"`
}

type KubeDeploymentOptions struct {
	// TODO add some custom options here if necessary
	// options means some custom params required by plugin running.
	// Such As How many rows do your want
	// You can use it in subtasks, and you need to pass it to main.go and pipelines.
	ConnectionId     uint64   `json:"connectionId"`
	Tasks            []string `json:"tasks,omitempty"`
	CreatedDateAfter string   `json:"createdDateAfter" mapstructure:"createdDateAfter,omitempty"`
}

type KubeDeploymentTaskData struct {
	Options          *KubeDeploymentOptions
	ApiClient        *api.ApiAsyncClient
	CreatedDateAfter *time.Time
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*KubeDeploymentOptions, errors.Error) {
	var op KubeDeploymentOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}
	return &op, nil
}

// NOT BEING USED
func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *KubeDeploymentTaskData) {
	data := taskCtx.GetData().(*KubeDeploymentTaskData)
	filteredData := *data
	filteredData.Options = &KubeDeploymentOptions{}
	*filteredData.Options = *data.Options
	var params = KubeDeploymentApiParams{
		ConnectionId: data.Options.ConnectionId,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}
