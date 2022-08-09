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
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/schema"
)

// @Summary test gitlab connection
// @Description Test gitlab Connection
// @Tags plugins/gitlab
// @Param body body models.TestConnectionRequest true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/test [POST]
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// decode
	var err error
	var connection models.TestConnectionRequest
	err = mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}
	// test connection
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", connection.Token),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, err
	}
	resBody := &models.ApiUserResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

// @Summary create gitlab connection
// @Description Create gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/connections [POST]
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// update from request and save to database
	connection := &models.GitlabConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch gitlab connection
// @Description Patch gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/connections/{connectionId} [PATCH]
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a gitlab connection
// @Description Delete a gitlab connection
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/connections/{connectionId} [DELETE]
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary get all gitlab connections
// @Description Get all gitlab connections
// @Tags plugins/gitlab
// @Success 200  {object} []models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/connections [GET]
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.GitlabConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get gitlab connection detail
// @Description Get gitlab connection detail
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitlab/connections/{connectionId} [GET]
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary get table info
// @Description Get table info map
// @Tags plugins/gitlab
// @Success 200
// @Router /plugins/gitlab/table [GET]
func GetTableInfo(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	info := make(map[string]interface{})

	err := core.GetTotalInfo().VisitTable("Gitlab", func(pluginName string, ts *core.TablesInfo) (err error) {

		tableInfo := make(map[string]interface{})
		info[pluginName] = &tableInfo

		return ts.Traversal(func(pluginName, name string, describe *string, table schema.Tabler) (err error) {
			tableData := make(map[string]reflect.StructField)
			tableInfo[name] = &tableData

			sf := utils.WalkFields(reflect.TypeOf(table), nil)
			for _, f := range sf {
				tableData[f.Name] = f
			}

			return
		})
	})

	return &core.ApiResourceOutput{Body: info}, err
}
