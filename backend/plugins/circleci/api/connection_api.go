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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type CircleciTestConnResponse struct {
	shared.ApiBody
}

func testConnection(ctx context.Context, connection *models.CircleciConn) (*CircleciTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, connection)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("/v2/me", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}

	body := CircleciTestConnResponse{}
	body.Success = true
	body.Message = "success"
	// output
	return &body, nil
}

// TestConnection test circleci connection
// @Summary test circleci connection
// @Description Test circleci Connection
// @Tags plugins/circleci
// @Param body body models.CircleciConnection true "json body"
// @Success 200  {object} CircleciTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	connection := &models.CircleciConn{}
	err := api.Decode(input.Body, connection, vld)
	if err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test circleci connection
// @Summary test circleci connection
// @Description Test circleci Connection
// @Tags plugins/circleci
// @Param connectionId path int true "connection ID"
// @Success 200  {object} CircleciTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), &connection.CircleciConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// PostConnections @Summary create circleci connection
// @Description Create circleci connection
// @Tags plugins/circleci
// @Param body body models.CircleciConnection true "json body"
// @Success 200  {object} models.CircleciConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection @Summary patch circleci connection
// @Description Patch circleci connection
// @Tags plugins/circleci
// @Param body body models.CircleciConnection true "json body"
// @Success 200  {object} models.CircleciConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection @Summary delete a circleci connection
// @Description Delete a circleci connection
// @Tags plugins/circleci
// @Success 200  {object} models.CircleciConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections @Summary get all circleci connections
// @Description Get all circleci connections
// @Tags plugins/circleci
// @Success 200  {object} []models.CircleciConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection @Summary get circleci connection detail
// @Description Get circleci connection detail
// @Tags plugins/circleci
// @Success 200  {object} models.CircleciConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/circleci/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

// GetConnectionTransformToDeployments return one connection deployments
// @Summary return one connection deployments
// @Description return one connection deployments
// @Tags plugins/circleci
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} map[string]interface{}
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/circleci/connections/{connectionId}/transform-to-deployments [POST]
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
		SELECT DISTINCT pipeline_number, name, project_slug, created_date
		FROM(
			SELECT pipeline_number, name, project_slug, created_date
			FROM _tool_circleci_workflows
			WHERE connection_id = ? 
			    AND (name REGEXP ?)
    			AND (? = '' OR name REGEXP ?)
			UNION
			SELECT w.pipeline_number, w.name, w.project_slug, w.created_date
			FROM _tool_circleci_jobs j 
			LEFT JOIN _tool_circleci_workflows w on w.id = j.workflow_id
			WHERE j.connection_id = ? 
			    AND (j.name REGEXP ?)
    			AND (? = '' OR j.name REGEXP ?)
		) AS t
		ORDER BY created_date DESC
	`, connectionId, deploymentPattern, productionPattern, productionPattern, connectionId, deploymentPattern, productionPattern, productionPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get")
	}
	defer cursor.Close()

	type selectFileds struct {
		PipelineNumber int
		Name           string
		ProjectSlug    string
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
			Name: fmt.Sprintf("#%d - %s", sf.PipelineNumber, sf.Name),
			URL:  CIRCLECI_URL + sf.ProjectSlug,
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

const CIRCLECI_URL = "https://app.circleci.com/pipelines/"
