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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/plugins/core"
)

var vld = validator.New()

/*
POST /plugins/tapd/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	var params models.TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		return nil, err
	}
	err = vld.Struct(params)
	if err != nil {
		return nil, err
	}
	// verify multiple token in parallel
	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	token := params.Auth
	apiClient, err := helper.NewApiClient(
		params.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", token),
		},
		3*time.Second,
		params.Proxy,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("verify token failed for %s %w", token, err)
	}
	res, err := apiClient.Get("/quickstart/testauth", nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("verify token failed for %s", token)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	// output
	return nil, nil
}

func findConnectionByInputParam(input *core.ApiResourceInput) (*models.TapdConnection, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return nil, fmt.Errorf("missing connectionsid")
	}
	tapdConnectionId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid connectionId")
	}

	return getTapdConnectionById(tapdConnectionId)
}

func getTapdConnectionById(id uint64) (*models.TapdConnection, error) {
	tapdConnection := &models.TapdConnection{}
	err := db.First(tapdConnection, id).Error
	if err != nil {
		return nil, err
	}

	// decrypt
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	tapdConnection.BasicAuthEncoded, err = core.Decrypt(encKey, tapdConnection.BasicAuthEncoded)
	if err != nil {
		return nil, err
	}

	return tapdConnection, nil
}

func mergeFieldsToTapdConnection(tapdConnection *models.TapdConnection, connections ...map[string]interface{}) error {
	// decode
	for _, connections := range connections {
		err := mapstructure.Decode(connections, tapdConnection)
		if err != nil {
			return err
		}
	}

	// validate
	vld := validator.New()
	err := vld.Struct(tapdConnection)
	if err != nil {
		return err
	}

	return nil
}

func refreshAndSaveTapdConnection(tapdConnection *models.TapdConnection, data map[string]interface{}) error {
	var err error
	// update fields from request body
	err = mergeFieldsToTapdConnection(tapdConnection, data)
	if err != nil {
		return err
	}

	// encrypt
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	if encKey == "" {
		// Randomly generate a bunch of encryption keys and set them to config
		encKey = core.RandomEncKey()
		v.Set(core.EncodeKeyEnvStr, encKey)
		err := config.WriteConfig(v)
		if err != nil {
			return err
		}
	}
	tapdConnection.BasicAuthEncoded, err = core.Encrypt(encKey, tapdConnection.BasicAuthEncoded)
	if err != nil {
		return err
	}

	// transaction for nested operations
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if tapdConnection.RateLimit == 0 {
		tapdConnection.RateLimit = 10800
	}
	if tapdConnection.ID > 0 {
		err = tx.Save(tapdConnection).Error
	} else {
		err = tx.Create(tapdConnection).Error
	}
	if err != nil {
		if common.IsDuplicateError(err) {
			return fmt.Errorf("tapd connections with name %s already exists", tapdConnection.Name)
		}
		return err
	}
	tapdConnection.BasicAuthEncoded, err = core.Decrypt(encKey, tapdConnection.BasicAuthEncoded)
	if err != nil {
		return err
	}
	return nil
}

/*
POST /plugins/tapd/connections
{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <tapd login email>:<tapd token> | base64`",
	"rateLimit": 10800,
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connections
	tapdConnection := &models.TapdConnection{}

	// update from request and save to database
	err := refreshAndSaveTapdConnection(tapdConnection, input.Body)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: tapdConnection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/tapd/connections/:connectionId
{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <tapd login email>:<tapd token> | base64`",
	"rateLimit": 10800,
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	tapdConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}

	// update from request and save to database
	err = refreshAndSaveTapdConnection(tapdConnection, input.Body)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: tapdConnection}, nil
}

/*
DELETE /plugins/tapd/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	tapdConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	err = db.Delete(tapdConnection).Error
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: tapdConnection}, nil
}

/*
GET /plugins/tapd/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdConnections := make([]models.TapdConnection, 0)
	err := db.Find(&tapdConnections).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdConnections}, nil
}

/*
GET /plugins/tapd/connections/:connectionId


{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <tapd login email>:<tapd token> | base64`",
	"rateLimit": 10800,
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}

	detail := &models.TapdConnectionDetail{
		TapdConnection: *tapdConnection,
	}
	return &core.ApiResourceOutput{Body: detail}, nil
}

// GET /plugins/tapd/connections/:connectionId/boards

func GetBoardsByConnectionId(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return nil, fmt.Errorf("missing connectionId")
	}
	tapdConnectionId, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid connectionId")
	}
	var tapdWorkspaces []models.TapdWorkspace
	err = db.Where("connection_Id = ?", tapdConnectionId).Find(&tapdWorkspaces).Error
	if err != nil {
		return nil, err
	}
	var workSpaceResponses []models.WorkspaceResponse
	for _, workSpace := range tapdWorkspaces {
		workSpaceResponses = append(workSpaceResponses, models.WorkspaceResponse{
			Id:    uint64(workSpace.ID),
			Title: workSpace.Name,
			Value: fmt.Sprintf("%v", workSpace.ID),
		})
	}
	return &core.ApiResourceOutput{Body: workSpaceResponses}, nil
}
