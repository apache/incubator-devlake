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
	"time"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	apiProject "github.com/apache/incubator-devlake/server/api/project"
)

type Connection struct {
	api.BaseConnection `mapstructure:",squash"`
	api.RestConnection `mapstructure:",squash"`
	api.AccessToken    `mapstructure:",squash"`
}

// CreateConnection FIXME
func (d *DevlakeClient) TestConnection(pluginName string, connection any) {
	d.testCtx.Helper()
	_ = sendHttpRequest[Connection](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/test", d.Endpoint, pluginName), connection)
}

// CreateConnection FIXME
func (d *DevlakeClient) CreateConnection(pluginName string, connection any) *Connection {
	d.testCtx.Helper()
	created := sendHttpRequest[Connection](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/connections", d.Endpoint, pluginName), connection)
	return &created
}

// ListConnections FIXME
func (d *DevlakeClient) ListConnections(pluginName string) []*Connection {
	d.testCtx.Helper()
	all := sendHttpRequest[[]*Connection](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections", d.Endpoint, pluginName), nil)
	return all
}

type BlueprintV2Config struct {
	Connection *plugin.BlueprintConnectionV200
	// Deprecated: to be deleted
	CreatedDateAfter *time.Time
	TimeAfter        *time.Time
	SkipOnFail       bool
	ProjectName      string
}

// CreateBasicBlueprintV2 FIXME
func (d *DevlakeClient) CreateBasicBlueprintV2(name string, config *BlueprintV2Config) models.Blueprint {
	settings := &models.BlueprintSettings{
		Version: "2.0.0",
		// Deprecated: to be deleted
		CreatedDateAfter: config.CreatedDateAfter,
		TimeAfter:        config.TimeAfter,
		Connections: ToJson([]*plugin.BlueprintConnectionV200{
			config.Connection,
		}),
	}
	blueprint := models.Blueprint{
		Name:        name,
		ProjectName: config.ProjectName,
		Mode:        models.BLUEPRINT_MODE_NORMAL,
		Plan:        nil,
		Enable:      true,
		CronConfig:  "manual",
		IsManual:    true,
		SkipOnFail:  config.SkipOnFail,
		Labels:      []string{"test-label"},
		Settings:    ToJson(settings),
	}
	d.testCtx.Helper()
	blueprint = sendHttpRequest[models.Blueprint](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/blueprints", d.Endpoint), &blueprint)
	return blueprint
}

type (
	ProjectPlugin struct {
		Name    string
		Options any
	}
	ProjectConfig struct {
		ProjectName        string
		ProjectDescription string
		EnableDora         bool
		MetricPlugins      []ProjectPlugin
	}
)

func (d *DevlakeClient) CreateProject(project *ProjectConfig) models.ApiOutputProject {
	var metrics []models.BaseMetric
	doraSeen := false
	for _, p := range project.MetricPlugins {
		if p.Name == "dora" {
			doraSeen = true
		}
		metrics = append(metrics, models.BaseMetric{
			PluginName:   p.Name,
			PluginOption: string(ToJson(p.Options)),
			Enable:       true,
		})
	}
	if project.EnableDora && !doraSeen {
		metrics = append(metrics, models.BaseMetric{
			PluginName:   "dora",
			PluginOption: string(ToJson(nil)),
			Enable:       true,
		})
	}
	return sendHttpRequest[models.ApiOutputProject](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/projects", d.Endpoint), &models.ApiInputProject{
		BaseProject: models.BaseProject{
			Name:        project.ProjectName,
			Description: project.ProjectDescription,
		},
		Enable:  Val(true),
		Metrics: &metrics,
	})
}

func (d *DevlakeClient) GetProject(projectName string) models.ApiOutputProject {
	return sendHttpRequest[models.ApiOutputProject](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/projects/%s", d.Endpoint, projectName), nil)
}

func (d *DevlakeClient) ListProjects() apiProject.PaginatedProjects {
	return sendHttpRequest[apiProject.PaginatedProjects](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/projects", d.Endpoint), nil)
}

func (d *DevlakeClient) CreateScope(pluginName string, connectionId uint64, scope any) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPut, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes", d.Endpoint, pluginName, connectionId), scope)
}

func (d *DevlakeClient) ListScopes(pluginName string, connectionId uint64) []any {
	return sendHttpRequest[[]any](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/connections/%d/scopes", d.Endpoint, pluginName, connectionId), nil)
}

func (d *DevlakeClient) CreateTransformRule(pluginName string, rules any) any {
	return sendHttpRequest[any](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s/transformation_rules", d.Endpoint, pluginName), rules)
}

func (d *DevlakeClient) ListTransformRules(pluginName string) []any {
	return sendHttpRequest[[]any](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s/transformation_rules?pageSize=20?page=1", d.Endpoint, pluginName), nil)
}

// CreateBasicBlueprint FIXME
func (d *DevlakeClient) CreateBasicBlueprint(name string, connection *plugin.BlueprintConnectionV100) models.Blueprint {
	settings := &models.BlueprintSettings{
		Version:     "1.0.0",
		Connections: ToJson([]*plugin.BlueprintConnectionV100{connection}),
	}
	blueprint := models.Blueprint{
		Name:       name,
		Mode:       models.BLUEPRINT_MODE_NORMAL,
		Plan:       nil,
		Enable:     true,
		CronConfig: "manual",
		IsManual:   true,
		Settings:   ToJson(settings),
	}
	d.testCtx.Helper()
	blueprint = sendHttpRequest[models.Blueprint](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/blueprints", d.Endpoint), &blueprint)
	return blueprint
}

// TriggerBlueprint FIXME
func (d *DevlakeClient) TriggerBlueprint(blueprintId uint64) models.Pipeline {
	d.testCtx.Helper()
	pipeline := sendHttpRequest[models.Pipeline](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/blueprints/%d/trigger", d.Endpoint, blueprintId), nil)
	return d.monitorPipeline(pipeline.ID)
}

// RunPipeline FIXME
func (d *DevlakeClient) RunPipeline(pipeline models.NewPipeline) models.Pipeline {
	d.testCtx.Helper()
	pipelineResult := sendHttpRequest[models.Pipeline](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/pipelines", d.Endpoint), &pipeline)
	return d.monitorPipeline(pipelineResult.ID)
}

// MonitorPipeline FIXME
func (d *DevlakeClient) monitorPipeline(id uint64) models.Pipeline {
	d.testCtx.Helper()
	var previousResult models.Pipeline
	endpoint := fmt.Sprintf("%s/pipelines/%d", d.Endpoint, id)
	coloredPrintf("calling:\n\t%s %s\nwith:\n%s\n", http.MethodGet, endpoint, string(ToCleanJson(false, nil)))
	for {
		time.Sleep(1 * time.Second)
		pipelineResult := sendHttpRequest[models.Pipeline](d.testCtx, d.timeout, debugInfo{
			print: false,
		}, http.MethodGet, fmt.Sprintf("%s/pipelines/%d", d.Endpoint, id), nil)
		if pipelineResult.Status == models.TASK_COMPLETED || pipelineResult.Status == models.TASK_FAILED {
			coloredPrintf("result: %s\n", ToCleanJson(true, &pipelineResult))
			return pipelineResult
		}
		if !reflect.DeepEqual(pipelineResult, previousResult) {
			coloredPrintf("result: %s\n", ToCleanJson(true, &pipelineResult))
		}
		previousResult = pipelineResult
	}
}
