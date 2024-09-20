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

package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

type TestConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func sanitizeConnection(connection interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if _, ok := result["token"]; ok {
		result["token"] = ""
	}
	return result, nil
}

func multiSanitizeConnections(connections []interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for _, c := range connections {
		result, err := sanitizeConnection(c)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (pa *pluginAPI) TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var result TestConnectionResult
	err := pa.invoker.Call("test-connection", bridge.DefaultContext, input.Body).Get(&result)
	if err != nil {
		body := shared.ApiBody{
			Success: false,
			Message: fmt.Sprintf("Error while testing connection: %s", err.Error()),
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 500}, nil
	} else {
		body := shared.ApiBody{
			Success: result.Success,
			Message: result.Message,
		}
		return &plugin.ApiResourceOutput{Body: body, Status: result.Status}, nil
	}
}

func (pa *pluginAPI) TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.connhelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	params := make(map[string]interface{})
	if data, err := json.Marshal(conn); err != nil {
		return nil, errors.Convert(err)
	} else {
		if err := json.Unmarshal(data, &params); err != nil {
			return nil, errors.Convert(err)
		}
	}

	necessaryParams := make(map[string]string)
	necessaryParams["proxy"] = cast.ToString(params["proxy"])
	necessaryParams["token"] = cast.ToString(params["token"])

	var result TestConnectionResult
	rpcCallErr := pa.invoker.Call("test-connection", bridge.DefaultContext, necessaryParams).Get(&result)
	if rpcCallErr != nil {
		body := shared.ApiBody{
			Success: false,
			Message: fmt.Sprintf("Error while testing connection: %s", rpcCallErr.Error()),
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 500}, nil
	} else {
		body := shared.ApiBody{
			Success: result.Success,
			Message: result.Message,
		}
		return &plugin.ApiResourceOutput{Body: body, Status: result.Status}, nil
	}
}

func (pa *pluginAPI) PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.connhelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	result, sanitizeErr := sanitizeConnection(conn)
	if sanitizeErr != nil {
		return nil, errors.Convert(sanitizeErr)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connections := pa.connType.NewSlice()
	err := pa.connhelper.List(connections)
	if err != nil {
		return nil, err
	}
	conns := connections.UnwrapSlice()
	if len(conns) == 0 {
		conns = []interface{}{}
		return &plugin.ApiResourceOutput{Body: conns}, nil
	}
	results, sanitizeErr := multiSanitizeConnections(conns)
	if sanitizeErr != nil {
		return nil, errors.Convert(sanitizeErr)
	}
	return &plugin.ApiResourceOutput{Body: results}, nil
}

func (pa *pluginAPI) GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.connhelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	result, sanitizeErr := sanitizeConnection(conn)
	if sanitizeErr != nil {
		return nil, errors.Convert(sanitizeErr)
	}
	return &plugin.ApiResourceOutput{Body: result}, nil
}

func (pa *pluginAPI) PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := pa.connType.New()
	err := pa.connhelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	conn := connection.Unwrap()
	result, sanitizeErr := sanitizeConnection(conn)
	if sanitizeErr != nil {
		return nil, errors.Convert(sanitizeErr)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.connhelper.Delete(pa.connType.New(), input)
}

func (pa *pluginAPI) GetConnectionTransformToDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
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
		SELECT DISTINCT id, name, url, start_time
		FROM(
			SELECT id, name, url, start_time
			FROM _tool_azuredevops_builds
      		WHERE connection_id = ?
		    	AND (name REGEXP ?)
    			AND (? = '' OR name REGEXP ?)
    		UNION
			SELECT b.id, b.name, b.url, b.start_time
			FROM _tool_azuredevops_jobs j
			LEFT JOIN _tool_azuredevops_builds b on CONCAT('azuredevops:Build:', b.connection_id, ':', b.id) = j.build_id
			WHERE j.connection_id = ?
			    AND (j.name REGEXP ?)
    			AND (? = '' OR j.name REGEXP ?)
		) AS t
		ORDER BY start_time DESC
	`, connectionId, deploymentPattern, productionPattern, productionPattern, connectionId, deploymentPattern, productionPattern, productionPattern)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get")
	}
	defer cursor.Close()

	type selectFileds struct {
		Id   int
		Name string
		URL  string
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
			Name: fmt.Sprintf("#%d - %s", sf.Id, sf.Name),
			URL:  sf.URL,
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
