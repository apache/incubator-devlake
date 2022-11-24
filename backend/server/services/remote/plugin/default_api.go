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
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

type ConnectionAPI struct {
	invoker  bridge.Invoker
	connType *models.DynamicTabler
	helper   *api.ConnectionApiHelper
}

func GetDefaultAPI(invoker bridge.Invoker, connType *models.DynamicTabler, helper *api.ConnectionApiHelper) map[string]map[string]plugin.ApiResourceHandler {
	api := &ConnectionAPI{
		invoker:  invoker,
		connType: connType,
		helper:   helper,
	}
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (c *ConnectionAPI) TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	err := c.invoker.Call("test-connection", bridge.DefaultContext, input.Body).Get()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ConnectionAPI) PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := c.connType.New()
	err := c.helper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn, Status: http.StatusOK}, nil
}

func (c *ConnectionAPI) ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connections := c.connType.NewSlice()
	err := c.helper.List(connections)
	if err != nil {
		return nil, err
	}
	conns := connections.Unwrap()
	return &plugin.ApiResourceOutput{Body: conns}, nil
}

func (c *ConnectionAPI) GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := c.connType.New()
	err := c.helper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn}, nil
}

func (c *ConnectionAPI) PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := c.connType.New()
	err := c.helper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn, Status: http.StatusOK}, nil
}

func (c *ConnectionAPI) DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := c.connType.New()
	err := c.helper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = c.helper.Delete(connection)
	conn := connection.Unwrap()
	return &plugin.ApiResourceOutput{Body: conn}, err
}
