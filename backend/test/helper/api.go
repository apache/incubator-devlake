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
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"reflect"
	"time"
)

// CreateConnection FIXME
func (d *DevlakeClient) CreateConnection(dest string, connection interface{}) *api.BaseConnection {
	d.testCtx.Helper()
	created := sendHttpRequest[api.BaseConnection](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodPost, fmt.Sprintf("%s/plugins/%s", d.Endpoint, dest), connection)
	return &created
}

// ListConnections FIXME
func (d *DevlakeClient) ListConnections(dest string) []*api.BaseConnection {
	d.testCtx.Helper()
	all := sendHttpRequest[[]*api.BaseConnection](d.testCtx, d.timeout, debugInfo{
		print:      true,
		inlineJson: false,
	}, http.MethodGet, fmt.Sprintf("%s/plugins/%s", d.Endpoint, dest), nil)
	return all
}

// CreateBasicBlueprintV2 FIXME This has to actually be adapted to v2
func (d *DevlakeClient) CreateBasicBlueprintV2(name string, connection *plugin.BlueprintConnectionV200) models.Blueprint {
	settings := &models.BlueprintSettings{
		Version:     "2.0.0",
		Connections: ToJson([]*plugin.BlueprintConnectionV200{connection}),
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
