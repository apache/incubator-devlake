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
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/server/api/shared"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

type JiraTestConnResponse struct {
	shared.ApiBody
	Connection *models.JiraConn
}

// @Summary test jira connection
// @Description Test Jira Connection
// @Tags plugins/jira
// @Param body body models.JiraConn true "json body"
// @Success 200  {object} JiraTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.JiraConn
	e := mapstructure.Decode(input.Body, &connection)
	if e != nil {
		return nil, errors.Convert(e)
	}
	e = vld.StructExcept(connection, "BasicAuth", "AccessToken")
	if e != nil {
		return nil, errors.Convert(e)
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	// serverInfo checking
	res, err := apiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, err
	}
	serverInfoFail := "Failed testing the serverInfo: [ " + res.Request.URL.String() + " ]"
	// check if `/rest/` was missing
	if res.StatusCode == http.StatusNotFound && !strings.HasSuffix(connection.Endpoint, "/rest/") {
		endpointUrl, err := url.Parse(connection.Endpoint)
		if err != nil {
			return nil, errors.Convert(err)
		}
		refUrl, err := url.Parse("/rest/")
		if err != nil {
			return nil, errors.Convert(err)
		}
		restUrl := endpointUrl.ResolveReference(refUrl)
		return nil, errors.NotFound.New(fmt.Sprintf("Seems like an invalid Endpoint URL, please try %s", restUrl.String()))
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("Error username/password")
	}

	resBody := &models.JiraServerInfo{}
	err = api.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// check version
	if resBody.DeploymentType == models.DeploymentServer {
		// only support 8.x.x or higher
		if versions := resBody.VersionNumbers; len(versions) == 3 && versions[0] < 7 {
			return nil, errors.Default.New(fmt.Sprintf("%s Support JIRA Server 8+ only", serverInfoFail))
		}
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("%s unexpected status code: %d", serverInfoFail, res.StatusCode))
	}

	// verify credential
	getStatusFail := "an error occurred while making request to `/rest/agile/1.0/board`"
	res, err = apiClient.Get("agile/1.0/board", nil, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, getStatusFail)
	}
	getStatusFail += ": [ " + res.Request.URL.String() + " ]"

	errMsg := ""
	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("it might you use the right token(password) but with the wrong username.please check your username/password")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("%s Unexpected [%s] status code: %d %s", getStatusFail, res.Request.URL, res.StatusCode, errMsg))
	}
	body := JiraTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// @Summary create jira connection
// @Description Create Jira connection
// @Tags plugins/jira
// @Param body body models.JiraConnection true "json body"
// @Success 200  {object} models.JiraConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.JiraConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch jira connection
// @Description Patch Jira connection
// @Tags plugins/jira
// @Param body body models.JiraConnection true "json body"
// @Success 200  {object} models.JiraConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a jira connection
// @Description Delete a Jira connection
// @Tags plugins/jira
// @Success 200  {object} models.JiraConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all jira connections
// @Description Get all Jira connections
// @Tags plugins/jira
// @Success 200  {object} []models.JiraConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.JiraConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get jira connection detail
// @Description Get Jira connection detail
// @Tags plugins/jira
// @Success 200  {object} models.JiraConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jira/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, err
}
