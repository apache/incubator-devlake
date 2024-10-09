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
	"net/http"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type BambooTestConnResponse struct {
	shared.ApiBody
	Connection *models.BambooConn
}

func testConnection(ctx context.Context, connection models.BambooConn) (*BambooTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// test connection
	_, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	connection = connection.Sanitize()
	if err != nil {
		return nil, err
	}
	body := BambooTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	return &body, nil
}

// TestConnection test bamboo connection
// @Summary test bamboo connection
// @Description Test bamboo Connection
// @Tags plugins/bamboo
// @Param body body models.BambooConn true "json body"
// @Success 200  {object} BambooTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bamboo/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.BambooConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test bamboo connection
// @Summary test bamboo connection
// @Description Test bamboo Connection
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Success 200  {object} BambooTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.FindByPk(input)
	if err != nil {
		return nil, err
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection.BambooConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// @Summary create bamboo connection
// @Description Create bamboo connection
// @Tags plugins/bamboo
// @Param body body models.BambooConnection true "json body"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch bamboo connection
// @Description Patch bamboo connection
// @Tags plugins/bamboo
// @Param body body models.BambooConnection true "json body"
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a bamboo connection
// @Description Delete a bamboo connection
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all bamboo connections
// @Description Get all bamboo connections
// @Tags plugins/bamboo
// @Success 200  {object} []models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get bamboo connection detail
// @Description Get bamboo connection detail
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Success 200  {object} models.BambooConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/bamboo/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

// GetConnectionDeployments return one connection deployments
// @Summary return one connection deployments
// @Description return one connection deployments
// @Tags plugins/bamboo
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {array} string "List of Environment Names"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/deployments [GET]
func GetConnectionDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	db := basicRes.GetDal()
	connectionId := input.Params["connectionId"]
	var environments []string
	err := db.All(&environments,
		dal.From(&models.BambooDeployBuild{}),
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
// @Tags plugins/bamboo
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} map[string]interface{}
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/transform-to-deployments [POST]
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
		SELECT DISTINCT plan_build_key, link_href, build_started_time
		FROM(
			SELECT plan_build_key, link_href, build_started_time
			FROM _tool_bamboo_plan_builds
			WHERE connection_id = ? 
			    AND (plan_name REGEXP ?)
    			AND (? = '' OR plan_name REGEXP ?)
			UNION
			SELECT p.plan_build_key, p.link_href, p.build_started_time
			FROM _tool_bamboo_job_builds j
			LEFT JOIN _tool_bamboo_plan_builds p on p.plan_build_key = j.plan_build_key
			WHERE j.connection_id = ? 
			    AND (j.job_name REGEXP ?)
    			AND (? = '' OR j.job_name REGEXP ?)
			ORDER BY build_started_time DESC
		) AS t
		ORDER BY build_started_time DESC
	`, connectionId, deploymentPattern, productionPattern, productionPattern, connectionId, deploymentPattern, productionPattern, productionPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get")
	}
	defer cursor.Close()

	type selectFileds struct {
		PlanBuildKey string
		LinkHref     string
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
		url, err := tasks.GetBambooHomePage(sf.LinkHref)
		if err != nil {
			url = sf.LinkHref
		}
		transformed := transformedFields{
			Name: sf.PlanBuildKey,
			URL:  url + "/browse/" + sf.PlanBuildKey,
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
