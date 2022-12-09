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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
)

// GeneratePlanJsonV200 generates pipeline plan according v2.0.0 definition
func GeneratePlanJsonV200(
	projectName string,
	sources *models.BlueprintSettings,
	metrics map[string]json.RawMessage,
) (core.PipelinePlan, errors.Error) {
	// generate plan and collect scopes
	plan, scopes, err := genPlanJsonV200(projectName, sources, metrics)
	if err != nil {
		return nil, err
	}
	// refresh project_mapping table to reflect project/scopes relationship
	if len(scopes) > 0 {
		e := db.Where("project_name = ?", projectName).Delete(&crossdomain.ProjectMapping{}).Error
		if e != nil {
			return nil, errors.Convert(err)
		}
		for _, scope := range scopes {
			e = basicRes.GetDal().CreateOrUpdate(scope)
			if e != nil {
				return nil, errors.Convert(err)
			}
		}
	}
	return plan, err
}

func genPlanJsonV200(
	projectName string,
	sources *models.BlueprintSettings,
	metrics map[string]json.RawMessage,
) (core.PipelinePlan, []core.Scope, errors.Error) {
	connections := make([]*core.BlueprintConnectionV200, 0)
	err := errors.Convert(json.Unmarshal(sources.Connections, &connections))
	if err != nil {
		return nil, nil, err
	}

	// make plan for data-source plugins fist. generate plan for each
	// connections, then merge them into one legitimate plan and collect the
	// scopes produced by the data-source plugins
	sourcePlans := make([]core.PipelinePlan, len(connections))
	scopes := make([]core.Scope, 0, len(connections))
	for i, connection := range connections {
		if len(connection.Scopes) == 0 {
			return nil, nil, errors.Default.New(fmt.Sprintf("connections[%d].scope is empty", i))
		}
		plugin, err := core.GetPlugin(connection.Plugin)
		if err != nil {
			return nil, nil, err
		}
		if pluginBp, ok := plugin.(core.DataSourcePluginBlueprintV200); ok {
			var pluginScopes []core.Scope
			sourcePlans[i], pluginScopes, err = pluginBp.MakeDataSourcePipelinePlanV200(
				connection.ConnectionId,
				connection.Scopes,
			)
			if err != nil {
				return nil, nil, err
			}
			// collect scopes for the project. a github repository may produces
			// 2 scopes, 1 repo and 1 board
			scopes = append(scopes, pluginScopes...)
		} else {
			return nil, nil, errors.Default.New(
				fmt.Sprintf("plugin %s does not support DataSourcePluginBlueprintV200", connection.Plugin),
			)
		}
	}
	// make plans for metric plugins
	metricPlans := make([]core.PipelinePlan, len(metrics))
	i := 0
	for metricPluginName, metricPluginOptJson := range metrics {
		plugin, err := core.GetPlugin(metricPluginName)
		if err != nil {
			return nil, nil, err
		}
		if pluginBp, ok := plugin.(core.MetricPluginBlueprintV200); ok {
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
	plan := SequencializePipelinePlans(
		ParallelizePipelinePlans(sourcePlans...),
		ParallelizePipelinePlans(metricPlans...),
	)
	return plan, scopes, err
}
