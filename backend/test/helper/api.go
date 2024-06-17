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

package helper

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/services"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/server/api/blueprints"
	apiProject "github.com/apache/incubator-devlake/server/api/project"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/stretchr/testify/require"
)

// CreateConnection FIXME
func (d *DevlakeClient) TestConnection(pluginName string, connection any) {
	d.testCtx.Helper()
	_ = sendHttpRequest[Connection](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/test", d.Endpoint, pluginName), nil, connection)
}

// CreateConnection FIXME
func (d *DevlakeClient) CreateConnection(pluginName string, connection any) *Connection {
	d.testCtx.Helper()
	created := sendHttpRequest[Connection](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/connections", d.Endpoint, pluginName), nil, connection)
	return &created
}

// ListConnections FIXME
func (d *DevlakeClient) ListConnections(pluginName string) []*Connection {
	d.testCtx.Helper()
	all := sendHttpRequest[[]*Connection](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections", d.Endpoint, pluginName), nil, nil)
	return all
}

// DeleteConnection FIXME
func (d *DevlakeClient) DeleteConnection(pluginName string, connectionId uint64) services.BlueprintProjectPairs {
	d.testCtx.Helper()
	refs := sendHttpRequest[shared.TypedApiBody[services.BlueprintProjectPairs]](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodDelete, fmt.Sprintf("%s/plugins/%s/connections/%d", d.Endpoint, pluginName, connectionId), nil, nil)
	return refs.Data
}

// CreateBasicBlueprintV2 FIXME
func (d *DevlakeClient) CreateBasicBlueprintV2(name string, config *BlueprintV2Config) models.Blueprint {
	blueprint := models.Blueprint{
		Name:        name,
		ProjectName: config.ProjectName,
		Mode:        models.BLUEPRINT_MODE_NORMAL,
		Plan:        nil,
		Enable:      true,
		CronConfig:  "manual",
		IsManual:    true,
		SyncPolicy: models.SyncPolicy{
			SkipOnFail: config.SkipOnFail,
		},
		Labels: []string{"test-label"},
		Connections: []*models.BlueprintConnection{
			config.Connection,
		},
	}
	d.testCtx.Helper()
	blueprint = sendHttpRequest[models.Blueprint](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/blueprints", d.Endpoint), nil, &blueprint)
	return blueprint
}

// PatchBasicBlueprintV2 FIXME
func (d *DevlakeClient) PatchBasicBlueprintV2(blueprintId uint64, name string, config *BlueprintV2Config) models.Blueprint {
	blueprint := models.Blueprint{
		Name:        name,
		ProjectName: config.ProjectName,
		Mode:        models.BLUEPRINT_MODE_NORMAL,
		Plan:        nil,
		Enable:      true,
		CronConfig:  "manual",
		IsManual:    true,
		SyncPolicy: models.SyncPolicy{
			SkipOnFail: config.SkipOnFail,
			TimeAfter: func() *time.Time {
				t, _ := time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
				return &t
			}(),
		},
		Labels: []string{"test-label"},
		Connections: []*models.BlueprintConnection{
			config.Connection,
		},
	}
	d.testCtx.Helper()
	blueprint = sendHttpRequest[models.Blueprint](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPatch, fmt.Sprintf("%s/blueprints/%d", d.Endpoint, blueprintId), nil, &blueprint)
	return blueprint
}

func (d *DevlakeClient) ListBlueprints() blueprints.PaginatedBlueprint {
	return sendHttpRequest[blueprints.PaginatedBlueprint](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/blueprints", d.Endpoint), nil, nil)
}

func (d *DevlakeClient) GetBlueprint(blueprintId uint64) models.Blueprint {
	return sendHttpRequest[models.Blueprint](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/blueprints/%d", d.Endpoint, blueprintId), nil, nil)
}

func (d *DevlakeClient) DeleteBlueprint(blueprintId uint64) {
	sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodDelete, fmt.Sprintf("%s/blueprints/%d", d.Endpoint, blueprintId), nil, nil)
}

func (d *DevlakeClient) CreateProject(project *ProjectConfig) models.ApiOutputProject {
	var metrics []*models.BaseMetric
	doraSeen := false
	for _, p := range project.MetricPlugins {
		if p.Name == "dora" {
			doraSeen = true
		}
		metrics = append(metrics, &models.BaseMetric{
			PluginName:   p.Name,
			PluginOption: ToJson(p.Options),
			Enable:       true,
		})
	}
	if project.EnableDora && !doraSeen {
		metrics = append(metrics, &models.BaseMetric{
			PluginName:   "dora",
			PluginOption: ToJson(nil),
			Enable:       true,
		})
	}
	return sendHttpRequest[models.ApiOutputProject](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/projects", d.Endpoint), nil, &models.ApiInputProject{
		BaseProject: models.BaseProject{
			Name:        project.ProjectName,
			Description: project.ProjectDescription,
		},
		Enable:    Val(true),
		Metrics:   metrics,
		Blueprint: project.Blueprint,
	})
}

func (d *DevlakeClient) GetProject(projectName string) models.ApiOutputProject {
	return sendHttpRequest[models.ApiOutputProject](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/projects/%s", d.Endpoint, projectName), nil, nil)
}

func (d *DevlakeClient) ListProjects() apiProject.PaginatedProjects {
	return sendHttpRequest[apiProject.PaginatedProjects](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/projects", d.Endpoint), nil, nil)
}

func (d *DevlakeClient) DeleteProject(projectName string) {
	sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodDelete, fmt.Sprintf("%s/projects/%s", d.Endpoint, projectName), nil, nil)
}

func (d *DevlakeClient) CreateScopes(pluginName string, connectionId uint64, scopes ...any) any {
	request := map[string]any{
		"data": scopes,
	}
	return sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPut, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes", d.Endpoint, pluginName, connectionId), nil, request)
}

func (d *DevlakeClient) UpdateScope(pluginName string, connectionId uint64, scopeId string, scope any) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPatch, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes/%s", d.Endpoint, pluginName, connectionId, scopeId), nil, scope)
}

func (d *DevlakeClient) ListScopes(pluginName string, connectionId uint64, listBlueprints bool) ScopeListResponseOut {
	scopesRaw := sendHttpRequest[ScopeListResponseIn](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes?blueprints=%v", d.Endpoint, pluginName, connectionId, listBlueprints), nil, nil)
	var responses []ScopeResponse
	for _, scopeRaw := range scopesRaw.Scopes {
		responses = append(responses, getScopeResponse(scopeRaw))
	}
	return ScopeListResponseOut{
		Scopes: responses,
		Count:  scopesRaw.Count,
	}
}

func (d *DevlakeClient) GetScope(pluginName string, connectionId uint64, scopeId string, listBlueprints bool) ScopeResponse {
	scopeRaw := sendHttpRequest[map[string]any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes/%s?blueprints=%v", d.Endpoint, pluginName, connectionId, scopeId, listBlueprints), nil, nil)
	return getScopeResponse(scopeRaw)
}

func (d *DevlakeClient) DeleteScope(pluginName string, connectionId uint64, scopeId string, deleteDataOnly bool) services.BlueprintProjectPairs {
	return sendHttpRequest[services.BlueprintProjectPairs](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodDelete, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes/%s?delete_data_only=%v", d.Endpoint, pluginName, connectionId, scopeId, deleteDataOnly), nil, nil)
}

func (d *DevlakeClient) CreateScopeConfig(pluginName string, connectionId uint64, scopeConfig any) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/connections/%d/scope-configs",
		d.Endpoint, pluginName, connectionId), nil, scopeConfig)
}

func (d *DevlakeClient) PatchScopeConfig(pluginName string, connectionId uint64, scopeConfigId uint64, scopeConfig any) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPatch, fmt.Sprintf("%s/plugins/%s/connections/%d/scope-configs/%d",
		d.Endpoint, pluginName, connectionId, scopeConfigId), nil, scopeConfig)
}

func (d *DevlakeClient) ListScopeConfigs(pluginName string, connectionId uint64) []any {
	return sendHttpRequest[[]any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/scope-configs?pageSize=20&page=1",
		d.Endpoint, pluginName, connectionId), nil, nil)
}

func (d *DevlakeClient) GetScopeConfig(pluginName string, connectionId uint64, scopeConfigId uint64) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/scope-configs/%d",
		d.Endpoint, pluginName, connectionId, scopeConfigId), nil, nil)
}

func (d *DevlakeClient) DeleteScopeConfig(pluginName string, connectionId uint64, scopeConfigId uint64) {
	sendHttpRequest[any](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodDelete, fmt.Sprintf("%s/plugins/%s/connections/%d/scope-configs/%d",
		d.Endpoint, pluginName, connectionId, scopeConfigId), nil, nil)
}

func (d *DevlakeClient) RemoteScopes(query RemoteScopesQuery) RemoteScopesOutput {
	url := fmt.Sprintf("%s/plugins/%s/connections/%d/remote-scopes",
		d.Endpoint,
		query.PluginName,
		query.ConnectionId,
	)
	if query.Params == nil {
		query.Params = make(map[string]string)
	}
	if query.GroupId != "" {
		query.Params["groupId"] = query.GroupId
	}
	if query.PageToken != "" {
		query.Params["pageToken"] = query.PageToken
	}
	if len(query.Params) > 0 {
		url = url + "?" + mapToQueryString(query.Params)
	}
	return sendHttpRequest[RemoteScopesOutput](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, url, nil, nil)
}

// SearchRemoteScopes makes calls to the "scope API" indirectly. "Search" is the remote endpoint to hit.
func (d *DevlakeClient) SearchRemoteScopes(query SearchRemoteScopesQuery) SearchRemoteScopesOutput {
	return sendHttpRequest[SearchRemoteScopesOutput](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/search-remote-scopes?search=%s&page=%d&pageSize=%d&%s",
		d.Endpoint,
		query.PluginName,
		query.ConnectionId,
		query.Search,
		query.Page,
		query.PageSize,
		mapToQueryString(query.Params)),
		nil, nil)
}

// TriggerBlueprint FIXME
func (d *DevlakeClient) TriggerBlueprint(blueprintId uint64) models.Pipeline {
	d.testCtx.Helper()
	pipeline := sendHttpRequest[models.Pipeline](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/blueprints/%d/trigger", d.Endpoint, blueprintId), nil, nil)
	return d.monitorPipeline(pipeline.ID)
}

// RunPipeline FIXME
func (d *DevlakeClient) RunPipeline(pipeline models.NewPipeline) models.Pipeline {
	d.testCtx.Helper()
	pipelineResult := sendHttpRequest[models.Pipeline](d.testCtx, d.timeout, &testContext{
		client:       d,
		printPayload: true,
		inlineJson:   false,
	}, http.MethodPost, fmt.Sprintf("%s/pipelines", d.Endpoint), nil, &pipeline)
	return d.monitorPipeline(pipelineResult.ID)
}

func mapToQueryString(queryParams map[string]string) string {
	params := make([]string, 0)
	for k, v := range queryParams {
		params = append(params, k+"="+v)
	}
	return strings.Join(params, "&")
}

// MonitorPipeline FIXME
func (d *DevlakeClient) monitorPipeline(id uint64) models.Pipeline {
	d.testCtx.Helper()
	var previousResult models.Pipeline
	endpoint := fmt.Sprintf("%s/pipelines/%d", d.Endpoint, id)
	coloredPrintf("calling:\n\t%s %s\nwith:\n%s\n", http.MethodGet, endpoint, string(ToCleanJson(false, nil)))
	var pipelineResult models.Pipeline
	require.NoError(d.testCtx, runWithTimeout(d.pipelineTimeout, func() (bool, errors.Error) {
		pipelineResult = sendHttpRequest[models.Pipeline](d.testCtx, d.pipelineTimeout, &testContext{
			client:       d,
			printPayload: false,
		}, http.MethodGet, fmt.Sprintf("%s/pipelines/%d", d.Endpoint, id), nil, nil)
		if pipelineResult.Status == models.TASK_COMPLETED {
			coloredPrintf("result: %s\n", ToCleanJson(true, &pipelineResult))
			return true, nil
		}
		if pipelineResult.Status == models.TASK_FAILED {
			coloredPrintf("result: %s\n", ToCleanJson(true, &pipelineResult))
			return true, errors.Default.New("pipeline task failed")
		}
		if !reflect.DeepEqual(pipelineResult, previousResult) {
			coloredPrintf("result: %s\n", ToCleanJson(true, &pipelineResult))
		}
		previousResult = pipelineResult
		return false, nil
	}))
	return pipelineResult
}

func getScopeResponse(scopeRaw map[string]any) ScopeResponse {
	response := Cast[ScopeResponse](scopeRaw)
	response.Scope = scopeRaw
	return response
}
