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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

func (pa *pluginAPI) TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	err := pa.invoker.Call("test-connection", bridge.DefaultContext, input.Body).Err
	if err != nil {
		body := shared.ApiBody{
			Success: false,
			Message: err.Error(),
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
	} else {
		body := shared.ApiBody{Success: true}
		return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
	}
}

func (pa *pluginAPI) PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.helper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connections := pa.connType.NewSlice()
	err := pa.helper.List(connections)
	if err != nil {
		return nil, err
	}
	conns := connections.Unwrap()
	return &plugin.ApiResourceOutput{Body: conns}, nil
}

func (pa *pluginAPI) GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.helper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn}, nil
}

func (pa *pluginAPI) PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.helper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.helper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = pa.helper.Delete(connection)
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn}, err
}
