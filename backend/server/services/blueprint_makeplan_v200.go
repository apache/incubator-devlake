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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

// GeneratePlanJsonV200 generates pipeline plan according v2.0.0 definition
func GeneratePlanJsonV200(
	projectName string,
	connections []*coreModels.BlueprintConnection,
	metrics map[string]json.RawMessage,
	skipCollectors bool,
) (coreModels.PipelinePlan, errors.Error) {
	var err errors.Error
	// make plan for data-source coreModels fist. generate plan for each
	// connection, then merge them into one legitimate plan and collect the
	// scopes produced by the data-source plugins
	sourcePlans := make([]coreModels.PipelinePlan, len(connections))
	scopes := make([]plugin.Scope, 0, len(connections))
	for i, connection := range connections {
		if len(connection.Scopes) == 0 && connection.PluginName != `webhook` && connection.PluginName != `jenkins` {
			// webhook needn't scopes
			// jenkins may upgrade from v100 and its scope is empty
			return nil, errors.Default.New(fmt.Sprintf("connections[%d].scopes is empty", i))
		}

		p, err := plugin.GetPlugin(connection.PluginName)
		if err != nil {
			return nil, err
		}
		if pluginBp, ok := p.(plugin.DataSourcePluginBlueprintV200); ok {
			var pluginScopes []plugin.Scope
			sourcePlans[i], pluginScopes, err = pluginBp.MakeDataSourcePipelinePlanV200(
				connection.ConnectionId,
				connection.Scopes,
				skipCollectors,
			)
			if err != nil {
				return nil, errors.Default.Wrap(err, "make data source pipeline plan v200")
			}
			// collect scopes for the project. a github repository may produce
			// 2 scopes, 1 repo and 1 board
			scopes = append(scopes, pluginScopes...)
		} else {
			return nil, errors.Default.New(
				fmt.Sprintf("plugin %s does not support DataSourcePluginBlueprintV200", connection.PluginName),
			)
		}
	}

	// skip collectors
	if skipCollectors {
		for i, plan := range sourcePlans {
			sourcePlans[i] = removeCollectorTasks(plan)
		}

		// remove gitextractor plugin if it's not the only task
		for i, plan := range sourcePlans {
			for j, stage := range plan {
				newStage := make(coreModels.PipelineStage, 0, len(stage))
				hasGitExtractor := false
				for _, task := range stage {
					if task.Plugin != "gitextractor" {
						newStage = append(newStage, task)
					} else {
						hasGitExtractor = true
					}
				}
				if !hasGitExtractor || len(newStage) > 0 {
					sourcePlans[i][j] = newStage
				}
			}
		}
	}

	// make plans for metric plugins
	metricPlans := make([]coreModels.PipelinePlan, len(metrics))
	i := 0
	for metricPluginName, metricPluginOptJson := range metrics {
		p, err := plugin.GetPlugin(metricPluginName)
		if err != nil {
			return nil, errors.Default.Wrap(err, "get plugin")
		}
		if pluginBp, ok := p.(plugin.MetricPluginBlueprintV200); ok {
			// If we enable one metric plugin, even if it has nil option, we still process it
			if len(metricPluginOptJson) == 0 {
				metricPluginOptJson = json.RawMessage("{}")
			}
			metricPlans[i], err = pluginBp.MakeMetricPluginPipelinePlanV200(projectName, metricPluginOptJson)
			if err != nil {
				return nil, errors.Default.Wrap(err, "make metric plugin pipeline plan v200")
			}
			i++
		} else {
			return nil, errors.Default.New(
				fmt.Sprintf("plugin %s does not support MetricPluginBlueprintV200", metricPluginName),
			)
		}
	}
	var planForProjectMapping coreModels.PipelinePlan
	if projectName != "" {
		p, err := plugin.GetPlugin("org")
		if err != nil {
			return nil, errors.Default.Wrap(err, "get plugin org")
		}
		if pluginBp, ok := p.(plugin.ProjectMapper); ok {
			planForProjectMapping, err = pluginBp.MapProject(projectName, scopes)
			if err != nil {
				return nil, errors.Default.Wrap(err, "make project")
			}
		}
	}
	plan := SequentializePipelinePlans(
		planForProjectMapping,
		ParallelizePipelinePlans(sourcePlans...),
		ParallelizePipelinePlans(metricPlans...),
	)
	return plan, err
}

func removeCollectorTasks(plan coreModels.PipelinePlan) coreModels.PipelinePlan {
	for j, stage := range plan {
		for k, task := range stage {
			newSubtasks := make([]string, 0, len(task.Subtasks))
			for _, subtask := range task.Subtasks {
				if !strings.Contains(strings.ToLower(subtask), "collect") {
					newSubtasks = append(newSubtasks, subtask)
				}
			}
			task.Subtasks = newSubtasks
			plan[j][k] = task
		}
	}
	return plan
}
