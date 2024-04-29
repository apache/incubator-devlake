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

package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

type TestConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (pa *pluginAPI) TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var result TestConnectionResult
	err := pa.invoker.Call("test-connection", bridge.DefaultContext, input.Body).Get(&result)
	if err != nil {
		body := shared.ApiBody{
			Success: false,
			Message: fmt.Sprintf("Error while testing connection: %s", err.Error()),
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 500}, nil
	} else {
		body := shared.ApiBody{
			Success: result.Success,
			Message: result.Message,
		}
		return &plugin.ApiResourceOutput{Body: body, Status: result.Status}, nil
	}
}

func (pa *pluginAPI) TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := pa.dsHelper.ConnApi.GetMergedConnectionAny(input)
	if err != nil {
		return nil, err
	}
	conn := connection.(models.DynamicTabler).Unwrap()
	params := make(map[string]interface{})
	if data, err := json.Marshal(conn); err != nil {
		return nil, errors.Convert(err)
	} else {
		if err := json.Unmarshal(data, &params); err != nil {
			return nil, errors.Convert(err)
		}
	}

	necessaryParams := make(map[string]string)
	necessaryParams["proxy"] = cast.ToString(params["proxy"])
	necessaryParams["token"] = cast.ToString(params["token"])

	var result TestConnectionResult
	rpcCallErr := pa.invoker.Call("test-connection", bridge.DefaultContext, necessaryParams).Get(&result)
	if rpcCallErr != nil {
		body := shared.ApiBody{
			Success: false,
			Message: fmt.Sprintf("Error while testing connection: %s", rpcCallErr.Error()),
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 500}, nil
	} else {
		body := shared.ApiBody{
			Success: result.Success,
			Message: result.Message,
		}
		return &plugin.ApiResourceOutput{Body: body, Status: result.Status}, nil
	}
}

func (pa *pluginAPI) PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ConnApi.Post(input)
}

func (pa *pluginAPI) ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ConnApi.GetAll(input)
}

func (pa *pluginAPI) GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ConnApi.GetDetail(input)
}

func (pa *pluginAPI) PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ConnApi.Patch(input)
}

func (pa *pluginAPI) DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ConnApi.Delete(input)
}
