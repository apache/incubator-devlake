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
	"github.com/apache/incubator-devlake/api/shared"
	"net/http"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/github/models"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/plugins/core"
)

/*
POST /plugins/github/test
REQUEST BODY:
{
	"endpoint": "https://api.github.com/",
	"token": "github api access token"
}
RESPONSE BODY:
Success:
{
	"login": "xxx"
}
Failure:
{
	"success": false,
	"message": "invalid token"
}
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
	tokens := strings.Split(params.Token, ",")

	// verify multiple token in parallel
	type VerifyResult struct {
		err   error
		login string
	}
	results := make(chan VerifyResult)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		j := i + 1
		go func() {
			apiClient, err := helper.NewApiClient(
				context.TODO(),
				params.Endpoint,
				map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", token),
				},
				3*time.Second,
				params.Proxy,
				basicRes,
			)
			if err != nil {
				results <- VerifyResult{err: fmt.Errorf("verify token failed for #%v %s %w", j, token, err)}
				return
			}
			res, err := apiClient.Get("user", nil, nil)
			if err != nil {
				results <- VerifyResult{err: fmt.Errorf("verify token failed for #%v %s %w", j, token, err)}
				return
			}
			githubUserOfToken := &models.GithubUserOfToken{}
			err = helper.UnmarshalResponse(res, githubUserOfToken)
			if err != nil {
				results <- VerifyResult{err: fmt.Errorf("verify token failed for #%v %s %w", j, token, err)}
				return
			} else if githubUserOfToken.Login == "" {
				results <- VerifyResult{err: fmt.Errorf("invalid token for #%v %s", j, token)}
				return
			}
			results <- VerifyResult{login: githubUserOfToken.Login}
		}()
	}

	// collect verification results
	logins := make([]string, 0)
	msgs := make([]string, 0)
	i := 0
	for result := range results {
		if result.err != nil {
			msgs = append(msgs, result.err.Error())
		}
		logins = append(logins, result.login)
		i++
		if i == len(tokens) {
			close(results)
		}
	}
	if len(msgs) > 0 {
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}

	var githubApiResponse struct {
		shared.ApiBody
		Login string `json:"login"`
	}
	githubApiResponse.Success = true
	githubApiResponse.Message = "success"
	githubApiResponse.Login = strings.Join(logins, `,`)
	return &core.ApiResourceOutput{Body: githubApiResponse, Status: http.StatusOK}, nil
}

/*
POST /plugins/github/connections
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/github/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
DELETE /plugins/github/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/github/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.GithubConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections}, nil
}

/*
GET /plugins/github/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, err
}
