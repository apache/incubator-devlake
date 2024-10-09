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

	"github.com/apache/incubator-devlake/core/dal"
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
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
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
	if _, ok := input.Body["enableGraphql"]; !ok {
		input.Body["enableGraphql"] = true
	}
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
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
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
	if conn.AuthMethod == models.AppKey {
		if tokenTestResult, err := testGithubConnAppKeyAuth(ctx, conn); err != nil {
			return nil, errors.Convert(err)
		} else {
			githubApiResponse.Success = tokenTestResult.Success
			githubApiResponse.Message = tokenTestResult.Message
			githubApiResponse.Login = tokenTestResult.Login
			githubApiResponse.Installations = tokenTestResult.Installations
		}
	} else if conn.AuthMethod == models.AccessToken {
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
	AppId          string `json:"appId,omitempty"`
	InstallationID int    `json:"installationId,omitempty"`

	// AccessToken
	Token string `json:"token,omitempty"`

	Success       bool                           `json:"success"`
	Message       string                         `json:"message,omitempty"`
	Login         string                         `json:"login"`
	Warning       bool                           `json:"warning"`
	Installations []models.GithubAppInstallation `json:"installations,omitempty"`
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
		AuthMethod: models.AccessToken,
		Token:      sanitizeConn.Token,
		Success:    success,
		Message:    strings.Join(messages, ";\n"),
		Login:      githubUserOfToken.Login,
		Warning:    warning,
	}
	return &tokenTestResult, nil
}

func getInstallationsWithGithubConnAppKeyAuth(ctx context.Context, conn models.GithubConn) (*GitHubTestConnResult, error) {
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
		return nil, errors.BadInput.Wrap(err, "verify token(get app) failed")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	githubApp := &models.GithubApp{}
	err = api.UnmarshalResponse(res, githubApp)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token resp failed")
	} else if githubApp.Slug == "" {
		return nil, errors.BadInput.Wrap(err, "invalid token")
	}
	res, err = apiClient.Get("app/installations", nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)},
	})
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token(get app installations) failed")
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
		AuthMethod:     models.AppKey,
		AppId:          conn.AppId,
		InstallationID: conn.InstallationID,
		Success:        true,
		Message:        "success",
		Login:          githubApp.Slug,
		Installations:  *githubAppInstallations,
	}
	return &tokenTestResult, nil
}

func testGithubConnAppKeyAuth(ctx context.Context, conn models.GithubConn) (*GitHubTestConnResult, error) {
	// AppKey can only have one secretKey, can shouldn't have tokens.
	conn.Token = ""
	// I think connection with InstallationID needs another test
	// But it's to be determined. So just ignore it temporarily.
	conn.InstallationID = 0
	return getInstallationsWithGithubConnAppKeyAuth(ctx, conn)
}

func testExistingConnection(ctx context.Context, conn models.GithubConn) (*GithubMultiTestConnResponse, errors.Error) {
	if vld != nil {
		if err := vld.StructExcept(conn, "GithubAppKey", "GithubAccessToken"); err != nil {
			return nil, errors.Convert(err)
		}
	}
	githubApiResponse := &GithubMultiTestConnResponse{}
	if conn.AuthMethod == models.AppKey {
		if tokenTestResult, err := testGithubConnAppKeyAuth(ctx, conn); err != nil {
			return nil, errors.Convert(err)
		} else {
			githubApiResponse.Tokens = append(githubApiResponse.Tokens, tokenTestResult)
		}
	} else if conn.AuthMethod == models.AccessToken {
		tokens := strings.Split(conn.Token, ",")
		for _, token := range tokens {
			testGithubConn := conn
			testGithubConn.Token = token
			tokenTestResult, err := testGithubConnAccessTokenAuth(ctx, testGithubConn)
			if err != nil {
				// generate a failed message for current token
				tokenTestResult = &GitHubTestConnResult{
					AuthMethod:     models.AccessToken,
					AppId:          testGithubConn.AppId,
					InstallationID: testGithubConn.InstallationID,
					Token:          testGithubConn.Sanitize().Token,
					Success:        false,
					Message:        err.Error(),
					Login:          "",
					Warning:        false,
					Installations:  nil,
				}
			}
			githubApiResponse.Tokens = append(githubApiResponse.Tokens, tokenTestResult)
		}
	} else {
		return nil, errors.BadInput.New("invalid authentication method")
	}

	// resp.success is true by default
	githubApiResponse.Success = true
	githubApiResponse.Message = "success"
	for _, token := range githubApiResponse.Tokens {
		if !token.Success {
			githubApiResponse.Success = false
			githubApiResponse.Message = token.Message
			githubApiResponse.Causes = append(githubApiResponse.Causes, token.Message)
		}
	}

	return githubApiResponse, nil
}

// TestExistingConnection test github connection options
// @Summary test github connection
// @Description Test github Connection
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Success 200  {object} GithubMultiTestConnResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/github/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.Convert(err)
	}
	testConnectionResult, testConnectionErr := testExistingConnection(context.TODO(), connection.GithubConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	return &plugin.ApiResourceOutput{Body: testConnectionResult, Status: http.StatusOK}, nil
}

// GetConnectionDeployments return one connection deployments
// @Summary return one connection deployments
// @Description return one connection deployments
// @Tags plugins/github
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {array} string "List of Environment Names"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/deployments [GET]
func GetConnectionDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	db := basicRes.GetDal()
	connectionId := input.Params["connectionId"]
	var environments []string
	err := db.All(&environments,
		dal.From(&models.GithubDeployment{}),
		dal.Where("connection_id = ?", connectionId),
		dal.Select("DISTINCT environment"))
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body: environments,
	}, nil
}

// GetConnectionTransformToDeployments return one connection deployments
// @Summary return one connection deployments
// @Description return one connection deployments
// @Tags plugins/github
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} map[string]interface{}
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/transform-to-deployments [POST]
func GetConnectionTransformToDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	db := basicRes.GetDal()
	connectionId := input.Params["connectionId"]
	deploymentPattern := input.Body["deploymentPattern"]
	productionPattern := input.Body["productionPattern"]
	page, err := api.ParsePageParam(input.Body, "page", 1)
	if err != nil {
		return nil, errors.Default.New("invalid page value")
	}
	pageSize, err := api.ParsePageParam(input.Body, "pageSize", 10)
	if err != nil {
		return nil, errors.Default.New("invalid pageSize value")
	}

	cursor, err := db.RawCursor(`
		SELECT DISTINCT r.run_number, r.name, r.head_branch, r.html_url, r.run_started_at
		FROM (
			SELECT id, run_number, name, head_branch, html_url, run_started_at
			FROM _tool_github_runs
			WHERE connection_id = ? 
				AND (name REGEXP ?)
				AND (? = '' OR name REGEXP ? OR head_branch REGEXP ?)
			UNION
			SELECT r.id, r.run_number, r.name, r.head_branch, r.html_url, r.run_started_at
			FROM _tool_github_jobs j
			LEFT JOIN _tool_github_runs r ON j.run_id = r.id
			WHERE j.connection_id = ?
			    AND (j.name REGEXP ?)
    			AND (? = '' OR j.name REGEXP ?)
		) r
		ORDER BY r.run_started_at DESC
	`, connectionId, deploymentPattern, productionPattern, productionPattern, productionPattern, connectionId, deploymentPattern, productionPattern, productionPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get")
	}
	defer cursor.Close()

	type selectFileds struct {
		RunNumber  int
		Name       string
		HeadBranch string
		HTMLURL    string
	}
	type transformedFields struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	var allRuns []transformedFields
	for cursor.Next() {
		sf := &selectFileds{}
		err = db.Fetch(cursor, sf)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error on fetch")
		}
		// Directly transform and append to allRuns
		transformed := transformedFields{
			Name: fmt.Sprintf("#%d - %s", sf.RunNumber, sf.Name),
			URL:  sf.HTMLURL,
		}
		allRuns = append(allRuns, transformed)
	}
	// Calculate total count
	totalCount := len(allRuns)

	// Paginate in memory
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}
	pagedRuns := allRuns[start:end]

	// Return result containing paged runs and total count
	result := map[string]interface{}{
		"total": totalCount,
		"data":  pagedRuns,
	}
	return &plugin.ApiResourceOutput{
		Body: result,
	}, nil
}
