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
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

var publicPermissions = []string{"repo:status", "repo_deployment", "read:user", "read:org"}
var privatePermissions = []string{"repo"}
var parentPermissions = map[string]string{
	"repo:status":     "repo",
	"repo_deployment": "repo",
	"read:user":       "user",
	"read:org":        "admin:org",
}

const (
	AppKey      = "AppKey"
	AccessToken = "AccessToken"
)

// findMissingPerms returns the missing required permissions from the given user permissions
func findMissingPerms(userPerms map[string]bool, requiredPerms []string) []string {
	missingPerms := make([]string, 0)
	for _, pp := range requiredPerms {
		// either the specific permission or its parent permission(larger) is granted
		if !userPerms[pp] && !userPerms[parentPermissions[pp]] {
			missingPerms = append(missingPerms, pp)
		}
	}
	return missingPerms
}

type GithubTestConnResponse struct {
	shared.ApiBody
	Login         string                         `json:"login"`
	Warning       bool                           `json:"warning"`
	Installations []models.GithubAppInstallation `json:"installations"`
}

// TestConnection test github connection
// @Summary test github connection
// @Description Test github Connection
// @Tags plugins/github
// @Param body body models.GithubConn true "json body"
// @Success 200  {object} GithubTestConnResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var conn models.GithubConn
	e := mapstructure.Decode(input.Body, &conn)
	if e != nil {
		return nil, errors.Convert(e)
	}
	testConnectionResult, err := testConnection(context.TODO(), conn)
	if err != nil {
		return nil, errors.Convert(err)
	}
	return &plugin.ApiResourceOutput{Body: testConnectionResult, Status: http.StatusOK}, nil
}

// @Summary create github connection
// @Description Create github connection
// @Tags plugins/github
// @Param body body models.GithubConnection true "json body"
// @Success 200  {object} models.GithubConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch github connection
// @Description Patch github connection
// @Tags plugins/github
// @Param body body models.GithubConnection true "json body"
// @Success 200  {object} models.GithubConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a github connection
// @Description Delete a github connection
// @Tags plugins/github
// @Success 200  {object} models.GithubConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all github connections
// @Description Get all github connections
// @Tags plugins/github
// @Success 200  {object} []models.GithubConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get github connection detail
// @Description Get github connection detail
// @Tags plugins/github
// @Success 200  {object} models.GithubConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

func testConnection(ctx context.Context, conn models.GithubConn) (*GithubTestConnResponse, errors.Error) {
	if vld != nil {
		if err := vld.StructExcept(conn, "GithubAppKey", "GithubAccessToken"); err != nil {
			return nil, errors.Convert(err)
		}
	}
	githubApiResponse := &GithubTestConnResponse{}
	if conn.AuthMethod == AppKey {
		if tokenTestResult, err := testGithubConnAppKeyAuth(ctx, conn); err != nil {
			return nil, errors.Convert(err)
		} else {
			githubApiResponse.Success = tokenTestResult.Success
			githubApiResponse.Message = tokenTestResult.Message
			githubApiResponse.Login = tokenTestResult.Login
			githubApiResponse.Installations = tokenTestResult.Installations
		}
	} else if conn.AuthMethod == AccessToken {
		if tokenTestResult, err := testGithubConnAccessTokenAuth(ctx, conn); err != nil {
			return nil, errors.Convert(err)
		} else {
			githubApiResponse.Success = tokenTestResult.Success
			githubApiResponse.Warning = tokenTestResult.Warning
			githubApiResponse.Message = tokenTestResult.Message
			githubApiResponse.Login = tokenTestResult.Login
		}
	} else {
		return nil, errors.BadInput.New("invalid authentication method")
	}

	return githubApiResponse, nil
}

type GitHubTestConnResult struct {
	AuthMethod string `json:"auth_method"`

	// AppKey
	AppId          string `json:"appId"`
	InstallationID int    `json:"installationId"`

	// AccessToken
	Token string `json:"token"`

	Success       bool                           `json:"success"`
	Message       string                         `json:"message"`
	Login         string                         `json:"login"`
	Warning       bool                           `json:"warning"`
	Installations []models.GithubAppInstallation `json:"installations"`
}

type GithubMultiTestConnResponse struct {
	shared.ApiBody
	Tokens []*GitHubTestConnResult `json:"tokens"`
}

// testGithubConnAccessTokenAuth only works when conn has one token
func testGithubConnAccessTokenAuth(ctx context.Context, conn models.GithubConn) (*GitHubTestConnResult, error) {
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &conn)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing connection")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	githubUserOfToken := &models.GithubUserOfToken{}
	err = api.UnmarshalResponse(res, githubUserOfToken)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	} else if githubUserOfToken.Login == "" {
		return nil, errors.BadInput.Wrap(err, "invalid token")
	}
	success := false
	warning := false
	var messages []string
	// for github classic token, check permission
	if strings.HasPrefix(conn.Token, "ghp_") {
		scopes := res.Header.Get("X-OAuth-Scopes")
		// convert "X-OAuth-Scopes" header to user permissions map
		userPerms := map[string]bool{}
		for _, userPerm := range strings.Split(scopes, ", ") {
			userPerms[userPerm] = true
		}
		// check public repo permission
		missingPubPerms := findMissingPerms(userPerms, publicPermissions)
		success = len(missingPubPerms) == 0
		if !success {
			messages = append(messages, fmt.Sprintf(
				"Please check the field(s) %s",
				strings.Join(missingPubPerms, ", "),
			))
		}
		// check private repo permission
		missingPriPerms := findMissingPerms(userPerms, privatePermissions)
		warning = len(missingPriPerms) > 0
		if warning {
			msgFmt := "If you want to collect private repositories, please check the field(s) %s"
			if success {
				// @Startrekzky and @yumengwang03 firmly believe that this is critical for users to understand the message
				msgFmt = "This token is able to collect public repositories. " + msgFmt
			}
			messages = append(messages, fmt.Sprintf(
				msgFmt,
				strings.Join(missingPriPerms, ", "),
			))
		}
	}

	sanitizeConn := conn.Sanitize()
	tokenTestResult := GitHubTestConnResult{
		AuthMethod: AccessToken,
		Token:      sanitizeConn.Token,
		Success:    success,
		Message:    strings.Join(messages, ";\n"),
		Login:      githubUserOfToken.Login,
		Warning:    warning,
	}
	return &tokenTestResult, nil
}

func testGithubConnAppKeyAuth(ctx context.Context, conn models.GithubConn) (*GitHubTestConnResult, error) {
	// AppKey can only have one secretKey
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &conn)
	if err != nil {
		return nil, err
	}
	jwt, err := conn.GithubAppKey.CreateJwt()
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("app", nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)},
	})
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	githubApp := &models.GithubApp{}
	err = api.UnmarshalResponse(res, githubApp)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	} else if githubApp.Slug == "" {
		return nil, errors.BadInput.Wrap(err, "invalid token")
	}
	res, err = apiClient.Get("app/installations", nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)},
	})
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	githubAppInstallations := &[]models.GithubAppInstallation{}
	err = api.UnmarshalResponse(res, githubAppInstallations)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	}
	tokenTestResult := GitHubTestConnResult{
		AuthMethod:     AppKey,
		AppId:          conn.AppId,
		InstallationID: conn.InstallationID,
		Success:        true,
		Message:        "success",
		Login:          githubApp.Slug,
		Installations:  *githubAppInstallations,
	}
	return &tokenTestResult, nil
}

func testExistingConnection(ctx context.Context, conn models.GithubConn) (*GithubMultiTestConnResponse, errors.Error) {
	if vld != nil {
		if err := vld.StructExcept(conn, "GithubAppKey", "GithubAccessToken"); err != nil {
			return nil, errors.Convert(err)
		}
	}
	githubApiResponse := &GithubMultiTestConnResponse{}
	if conn.AuthMethod == AppKey {
		if tokenTestResult, err := testGithubConnAppKeyAuth(ctx, conn); err != nil {
			return nil, errors.Convert(err)
		} else {
			githubApiResponse.Tokens = append(githubApiResponse.Tokens, tokenTestResult)
		}
	} else if conn.AuthMethod == AccessToken {
		tokens := strings.Split(conn.Token, ",")
		for _, token := range tokens {
			testGithubConn := conn
			testGithubConn.Token = token
			if tokenTestResult, err := testGithubConnAccessTokenAuth(ctx, testGithubConn); err != nil {
				return nil, errors.Convert(err)
			} else {
				githubApiResponse.Tokens = append(githubApiResponse.Tokens, tokenTestResult)
			}
		}
	} else {
		return nil, errors.BadInput.New("invalid authentication method")
	}

	return githubApiResponse, nil
}

// TestExistingConnection test github connection options
// @Summary test github connection
// @Description Test github Connection
// @Tags plugins/github
// @Success 200  {object} GithubTestConnResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.FindByPk(input)
	if err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testExistingConnection(context.TODO(), connection.GithubConn)
	if testConnectionErr != nil {
		return nil, testConnectionErr
	}
	return &plugin.ApiResourceOutput{Body: testConnectionResult, Status: http.StatusOK}, nil
}
