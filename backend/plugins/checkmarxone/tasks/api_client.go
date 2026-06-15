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

package tasks

import (
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/plugins/checkmarxone/models"
	"net/http"
	"time"
)

type CheckmarxoneApiClient struct {
	client     *http.Client
	headers    map[string]string
	logger     log.Logger
	serverUrl  string
	username   string
	password   string
	clientId   string
	clientSecret string
	token      string
	tokenExpire time.Time
}

func NewCheckmarxoneApiClient(logger log.Logger, connection *models.CheckmarxoneConnection) (*CheckmarxoneApiClient, errors.Error) {
	client := &CheckmarxoneApiClient{
		client:       &http.Client{Timeout: 30 * time.Second},
		logger:       logger,
		serverUrl:    connection.ServerUrl,
		username:     connection.Username,
		password:     connection.Password,
		clientId:     connection.ClientId,
		clientSecret: connection.ClientSecret,
		headers: map[string]string{
			"Accept": "application/json",
		},
	}

	err := client.authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *CheckmarxoneApiClient) authenticate() errors.Error {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.clientId, c.clientSecret)))

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", auth),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	res, err := utils.HTTPRequest(
		"POST",
		fmt.Sprintf("%s/auth/oauth/token", c.serverUrl),
		nil,
		map[string]string{"grant_type": "client_credentials"},
		headers,
		c.client,
		false,
	)
	if err != nil {
		return errors.Default.Wrap(err, "failed to authenticate with CheckmarxOne")
	}

	if res.StatusCode != 200 {
		return errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("failed to authenticate: %s", res.Body))
	}

	type TokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var tokenResp TokenResponse
	err = utils.UnmarshalResponse(res, &tokenResp)
	if err != nil {
		return errors.Default.Wrap(err, "failed to parse token response")
	}

	c.token = tokenResp.AccessToken
	c.tokenExpire = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	c.headers["Authorization"] = fmt.Sprintf("Bearer %s", c.token)

	return nil
}

func (c *CheckmarxoneApiClient) GetProjects() ([]map[string]interface{}, errors.Error) {
	url := fmt.Sprintf("%s/api/projects", c.serverUrl)
	return c.fetch(url)
}

func (c *CheckmarxoneApiClient) GetFindings(projectId string) ([]map[string]interface{}, errors.Error) {
	url := fmt.Sprintf("%s/api/projects/%s/results-summary", c.serverUrl, projectId)
	return c.fetch(url)
}

func (c *CheckmarxoneApiClient) fetch(url string) ([]map[string]interface{}, errors.Error) {
	err := c.checkAndRefreshToken()
	if err != nil {
		return nil, err
	}

	res, err := utils.HTTPRequest("GET", url, nil, nil, c.headers, c.client, false)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to fetch data from CheckmarxOne")
	}

	if res.StatusCode != 200 {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("failed to fetch: %s", res.Body))
	}

	type DataResponse struct {
		Results []map[string]interface{} `json:"results"`
		Data    []map[string]interface{} `json:"data"`
	}

	var dataResp DataResponse
	err = utils.UnmarshalResponse(res, &dataResp)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to parse response")
	}

	if len(dataResp.Results) > 0 {
		return dataResp.Results, nil
	}
	return dataResp.Data, nil
}

func (c *CheckmarxoneApiClient) checkAndRefreshToken() errors.Error {
	if time.Now().Before(c.tokenExpire) {
		return nil
	}
	return c.authenticate()
}

func (c *CheckmarxoneApiClient) Close() {
	c.client.CloseIdleConnections()
}
