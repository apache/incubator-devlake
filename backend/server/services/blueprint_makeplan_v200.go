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

package services

import (
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

// GeneratePlanJsonV200 generates pipeline plan according v2.0.0 definition
func GeneratePlanJsonV200(
	projectName string,
	syncPolicy plugin.BlueprintSyncPolicy,
	sources *models.BlueprintSettings,
	metrics map[string]json.RawMessage,
) (plugin.PipelinePlan, errors.Error) {
	// generate plan and collect scopes
	plan, scopes, err := genPlanJsonV200(projectName, syncPolicy, sources, metrics)
	if err != nil {
		return nil, err
	}
	// save scopes to database
	if len(scopes) > 0 {
		for _, scope := range scopes {
			err = db.CreateOrUpdate(scope)
			if err != nil {
				scopeInfo := fmt.Sprintf("[Id:%s][Name:%s][TableName:%s]", scope.ScopeId(), scope.ScopeName(), scope.TableName())
				return nil, errors.Default.Wrap(err, fmt.Sprintf("failed to create scopes:[%s]", scopeInfo))
			}
		}
	}
	return plan, err
}

func genPlanJsonV200(
	projectName string,
	syncPolicy plugin.BlueprintSyncPolicy,
	sources *models.BlueprintSettings,
	metrics map[string]json.RawMessage,
) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	connections := make([]*plugin.BlueprintConnectionV200, 0)
	err := errors.Convert(json.Unmarshal(sources.Connections, &connections))
	if err != nil {
		return nil, nil, err
	}

	// make plan for data-source plugins fist. generate plan for each
	// connections, then merge them into one legitimate plan and collect the
	// scopes produced by the data-source plugins
	sourcePlans := make([]plugin.PipelinePlan, len(connections))
	scopes := make([]plugin.Scope, 0, len(connections))
	for i, connection := range connections {
		if len(connection.Scopes) == 0 && connection.Plugin != `webhook` && connection.Plugin != `jenkins` {
			// webhook needn't scopes
			// jenkins may upgrade from v100 and its' scope is empty
			return nil, nil, errors.Default.New(fmt.Sprintf("connections[%d].scopes is empty", i))
		}
		p, err := plugin.GetPlugin(connection.Plugin)
		if err != nil {
			return nil, nil, err
		}
		if pluginBp, ok := p.(plugin.DataSourcePluginBlueprintV200); ok {
			var pluginScopes []plugin.Scope
			sourcePlans[i], pluginScopes, err = pluginBp.MakeDataSourcePipelinePlanV200(
				connection.ConnectionId,
				connection.Scopes,
				syncPolicy,
			)
			if err != nil {
				return nil, nil, err
			}
			// collect scopes for the project. a github repository may produce
			// 2 scopes, 1 repo and 1 board
			scopes = append(scopes, pluginScopes...)
		} else {
			return nil, nil, errors.Default.New(
				fmt.Sprintf("plugin %s does not support DataSourcePluginBlueprintV200", connection.Plugin),
			)
		}
	}
	// make plans for metric plugins
	metricPlans := make([]plugin.PipelinePlan, len(metrics))
	i := 0
	for metricPluginName, metricPluginOptJson := range metrics {
		p, err := plugin.GetPlugin(metricPluginName)
		if err != nil {
			return nil, nil, err
		}
		if pluginBp, ok := p.(plugin.MetricPluginBlueprintV200); ok {
			// If we enable one metric plugin, even if it has nil option, we still process it
			if len(metricPluginOptJson) == 0 {
				metricPluginOptJson = json.RawMessage("{}")
			}
			metricPlans[i], err = pluginBp.MakeMetricPluginPipelinePlanV200(projectName, metricPluginOptJson)
			if err != nil {
				return nil, nil, err
			}
			i += 1
		} else {
			return nil, nil, errors.Default.New(
				fmt.Sprintf("plugin %s does not support MetricPluginBlueprintV200", metricPluginName),
			)
		}
	}
	var planForProjectMapping plugin.PipelinePlan
	if projectName != "" {
		p, err := plugin.GetPlugin("org")
		if err != nil {
			return nil, nil, err
		}
		if pluginBp, ok := p.(plugin.ProjectMapper); ok {
			planForProjectMapping, err = pluginBp.MapProject(projectName, scopes)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	plan := SequencializePipelinePlans(
		planForProjectMapping,
		ParallelizePipelinePlans(sourcePlans...),
		ParallelizePipelinePlans(metricPlans...),
	)
	return plan, scopes, err
}
