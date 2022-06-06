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
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var vld = validator.New()

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
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", connection.GetEncodedToken()),
		},
		3*time.Second,
		connection.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, err
	}
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
			return nil, fmt.Errorf("Support JIRA Server 8+ only")
		}
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
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
	// create a new connection
	jiraConnection := &models.JiraConnection{}

	// update from request and save to database
	err := helper.CreateConnection(input.Body, jiraConnection, db)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnection, Status: http.StatusOK}, nil
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
	jiraConnection := &models.JiraConnection{}
	err := helper.PatchConnection(input, jiraConnection, db)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnection}, nil
}

/*
DELETE /plugins/jira/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraConnectionID, err := helper.GetConnectionIdByInputParam(input.Params)
	if err != nil {
		return nil, err
	}
	// cascading delete
	err = db.Where("id = ?", jiraConnectionID).Delete(&models.JiraConnection{}).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraConnectionID}, nil
}

/*
GET /plugins/jira/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnections := make([]*models.JiraConnection, 0)

	err := helper.ListConnections(&jiraConnections, db)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraConnections, Status: http.StatusOK}, nil
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
	jiraConnection := &models.JiraConnection{}
	err := helper.GetConnection(input.Params, jiraConnection, db)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: jiraConnection}, nil
}
