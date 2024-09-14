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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type GitlabTestConnResponse struct {
	shared.ApiBody
	Connection *models.GitlabConn
}

func testConnection(ctx context.Context, connection models.GitlabConn) (*GitlabTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}

	// check API/read_api permissions
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", 1))
	query.Set("per_page", fmt.Sprintf("%v", 1))
	res, err := apiClient.Get("projects", query, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing api or read_api permissions")
	}

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.BadInput.New("token need api or read_api permissions scope")
	}

	connection = connection.Sanitize()
	body := GitlabTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection

	return &body, nil
}

// TestConnection test gitlab connection
// @Summary test gitlab connection
// @Description Test gitlab Connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConn true "json body"
// @Success 200  {object} GitlabTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.GitlabConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test gitlab connection
// @Summary test gitlab connection
// @Description Test gitlab Connection
// @Tags plugins/gitlab
// @Param connectionId path int true "connection ID"
// @Success 200  {object} GitlabTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.Convert(err)
	}
	if result, err := testConnection(context.TODO(), connection.GitlabConn); err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	} else {
		return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
	}
}

// @Summary create gitlab connection
// @Description Create gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch gitlab connection
// @Description Patch gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a gitlab connection
// @Description Delete a gitlab connection
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all gitlab connections
// @Description Get all gitlab connections
// @Tags plugins/gitlab
// @Success 200  {object} []models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get gitlab connection detail
// @Description Get gitlab connection detail
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

// GetConnectionDeployments return one connection deployments
// @Summary return one connection deployments
// @Description return one connection deployments
// @Tags plugins/gitlab
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {array} string "List of Environment Names"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/deployments [GET]
func GetConnectionDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	db := basicRes.GetDal()
	connectionId := input.Params["connectionId"]
	var environments []string
	err := db.All(&environments,
		dal.From(&models.GitlabDeployment{}),
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
// @Tags plugins/gitlab
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} map[string]interface{}
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/transform-to-deployments [POST]
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
		SELECT DISTINCT name, gitlab_id, web_url, started_at
		FROM (
			SELECT r.name, p.gitlab_id, p.web_url, p.started_at
			FROM _tool_gitlab_pipelines p
			LEFT JOIN _tool_gitlab_projects r on r.gitlab_id = p.project_id
			WHERE p.connection_id = ? 
				AND (ref REGEXP ?)
				AND (? = '' OR ref REGEXP ?)
			UNION
			SELECT r.name, p.gitlab_id, p.web_url, p.started_at
			FROM _tool_gitlab_pipelines p
			LEFT JOIN _tool_gitlab_projects r on r.gitlab_id = p.project_id
			LEFT JOIN _tool_gitlab_jobs j on j.pipeline_id = p.gitlab_id
			WHERE j.connection_id = ? 
			   	AND (j.name REGEXP ?)
    		   	AND (? = '' OR j.name REGEXP ?)
		) r
		ORDER BY r.started_at DESC
	`, connectionId, deploymentPattern, productionPattern, productionPattern, connectionId, deploymentPattern, productionPattern, productionPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get")
	}
	defer cursor.Close()

	type selectFileds struct {
		Name     string
		GitlabId int
		WebUrl   string
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
			Name: fmt.Sprintf("%s-#%d", sf.Name, sf.GitlabId),
			URL:  sf.WebUrl,
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
