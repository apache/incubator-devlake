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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

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
			"Authorization": fmt.Sprintf("Basic %v", connection.GetEncodedToken()),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}
	// serverInfo checking
	res, err := apiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, err
	}
	serverInfoFail := "You are failed on test the serverInfo: [ " + res.Request.URL.String() + " ]"
	// check if `/rest/` was missing
	if res.StatusCode == http.StatusNotFound && !strings.HasSuffix(connection.Endpoint, "/rest/") {
		endpointUrl, err := url.Parse(connection.Endpoint)
		if err != nil {
			return nil, err
		}
		refUrl, err := url.Parse("/rest/")
		if err != nil {
			return nil, err
		}
		restUrl := endpointUrl.ResolveReference(refUrl)
		return nil, errors.NewNotFound(fmt.Sprintf("Seems like an invalid Endpoint URL, please try %s", restUrl.String()))
	}
	resBody := &models.JiraServerInfo{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}
	// check version
	if resBody.DeploymentType == models.DeploymentServer {
		// only support 8.x.x or higher
		if versions := resBody.VersionNumbers; len(versions) == 3 && versions[0] < 8 {
			return nil, fmt.Errorf("%s Support JIRA Server 8+ only", serverInfoFail)
		}
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s unexpected status code: %d", serverInfoFail, res.StatusCode)
	}

	// usersSearch checking
	usersSearchFail := "You are suceess on test the serverInfo but failed to test on userSearch"
	res, err = apiClient.Get("api/3/users/search", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%s %s", usersSearchFail, err)
	}
	usersSearchFail += ": [ " + res.Request.URL.String() + " ]"

	errMsg := ""
	if res.StatusCode == http.StatusForbidden {
		resErrBody := &models.JiraErrorInfo{}
		err = helper.UnmarshalResponse(res, resErrBody)
		if err != nil {
			return nil, fmt.Errorf("%s Unexpected status code: %d,and UnmarshalResponse error %s", usersSearchFail, res.StatusCode, err)
		}
		for _, em := range resErrBody.ErrorMessages {
			if em == "error.no-permission" {
				return nil, fmt.Errorf("%s We get the error %s ,it might you use the right token(password) but with the wrong username.please check your password", usersSearchFail, em)
			}
			errMsg += em + " \r\n"
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s Unexpected [%s] status code: %d %s", usersSearchFail, res.Request.URL, res.StatusCode, errMsg)
	}

	return nil, nil
}

/*
POST /plugins/jira/connections
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jira api access token"
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// update from request and save to database
	connection := &models.JiraConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/jira/connections/:connectionId
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jira api access token"
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/jira/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/jira/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.JiraConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/jira/connections/:connectionId


{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jira api access token"
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JiraConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}
