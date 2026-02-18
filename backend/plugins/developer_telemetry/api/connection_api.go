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
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coremodels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/developer_telemetry/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type DeveloperTelemetryTestConnResponse struct {
	shared.ApiBody
	Connection *models.DeveloperTelemetryConnection `json:"connection"`
}

type DeveloperTelemetryConnectionResponse struct {
	*models.DeveloperTelemetryConnection
	ApiKey *coremodels.ApiKey `json:"apiKey,omitempty"`
}

func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.DeveloperTelemetryConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}

	if connection.Name == "" {
		return nil, errors.BadInput.New("connection name is required")
	}

	return &plugin.ApiResourceOutput{Body: DeveloperTelemetryTestConnResponse{
		Connection: connection,
		ApiBody:    shared.ApiBody{Success: true, Message: "success"},
	}, Status: http.StatusOK}, nil
}

func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.DeveloperTelemetryConnection{}
	tx := basicRes.GetDal().Begin()
	err := connectionHelper.CreateWithTx(tx, connection, input)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		if strings.Contains(err.Error(), "the connection name already exists (400)") {
			return nil, errors.BadInput.New(fmt.Sprintf("A developer telemetry connection with name %s already exists.", connection.Name))
		}
		return nil, err
	}
	logger.Info("connection: %+v", connection)
	name := apiKeyHelper.GenApiKeyNameForPlugin(pluginName, connection.ID)
	allowedPath := fmt.Sprintf("/plugins/%s/connections/%d/.*", pluginName, connection.ID)
	extra := fmt.Sprintf("connectionId:%d", connection.ID)
	apiKeyRecord, err := apiKeyHelper.CreateForPlugin(tx, input.User, name, pluginName, allowedPath, extra)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "CreateForPlugin")
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		logger.Info("transaction commit: %s", err)
	}

	response := &DeveloperTelemetryConnectionResponse{
		DeveloperTelemetryConnection: connection,
		ApiKey:                       apiKeyRecord,
	}
	logger.Info("api output connection: %+v", response)

	return &plugin.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.DeveloperTelemetryConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.DeveloperTelemetryConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		logger.Error(err, "query connection")
		return nil, err
	}

	tx := basicRes.GetDal().Begin()

	err = tx.Delete(connection, dal.Where("id = ?", connection.ID))
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "delete connection: %d", connection.ID)
		return nil, err
	}

	extra := fmt.Sprintf("connectionId:%d", connection.ID)
	err = apiKeyHelper.DeleteForPlugin(tx, pluginName, extra)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "DeleteForPlugin")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Info("transaction commit: %s", err)
	}

	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}

func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.DeveloperTelemetryConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.DeveloperTelemetryConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}
