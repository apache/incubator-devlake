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

package api

import (
	"net/http"
	"strconv"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/checkmarxone/models"
)

func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.CheckmarxoneConnection{}
	err := input.GetBody(connection)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid request body")
	}

	basicRes := input.Ctx.Value(plugin.CTX_KEY_BASIC_RES).(plugin.BasicRes)
	if err := basicRes.GetDal().Create(connection); err != nil {
		return nil, errors.Default.Wrap(err, "failed to create connection")
	}

	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusCreated}, nil
}

func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.CheckmarxoneConnection
	basicRes := input.Ctx.Value(plugin.CTX_KEY_BASIC_RES).(plugin.BasicRes)

	if err := basicRes.GetDal().All(&connections); err != nil {
		return nil, errors.Default.Wrap(err, "failed to list connections")
	}

	return &plugin.ApiResourceOutput{Body: connections}, nil
}

func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId := input.Params["connectionId"]
	connId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, errors.BadInput.New("invalid connection id")
	}

	connection := &models.CheckmarxoneConnection{}
	basicRes := input.Ctx.Value(plugin.CTX_KEY_BASIC_RES).(plugin.BasicRes)

	if err := basicRes.GetDal().First(connection, map[string]interface{}{"id": connId}); err != nil {
		return nil, errors.NotFound.Wrap(err, "connection not found")
	}

	return &plugin.ApiResourceOutput{Body: connection}, nil
}

func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId := input.Params["connectionId"]
	connId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, errors.BadInput.New("invalid connection id")
	}

	connection := &models.CheckmarxoneConnection{}
	basicRes := input.Ctx.Value(plugin.CTX_KEY_BASIC_RES).(plugin.BasicRes)

	err2 := input.GetBody(connection)
	if err2 != nil {
		return nil, errors.BadInput.Wrap(err2, "invalid request body")
	}

	connection.ID = connId
	if err := basicRes.GetDal().Update(connection); err != nil {
		return nil, errors.Default.Wrap(err, "failed to update connection")
	}

	return &plugin.ApiResourceOutput{Body: connection}, nil
}

func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId := input.Params["connectionId"]
	connId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, errors.BadInput.New("invalid connection id")
	}

	basicRes := input.Ctx.Value(plugin.CTX_KEY_BASIC_RES).(plugin.BasicRes)
	connection := &models.CheckmarxoneConnection{}

	if err := basicRes.GetDal().Delete(connection, map[string]interface{}{"id": connId}); err != nil {
		return nil, errors.Default.Wrap(err, "failed to delete connection")
	}

	return &plugin.ApiResourceOutput{Status: http.StatusNoContent}, nil
}
