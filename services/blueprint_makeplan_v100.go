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
	"github.com/apache/incubator-devlake/plugins/core"
)

// GeneratePlanJsonV100 generates pipeline plan according v1.0.0 definition
func GeneratePlanJsonV100(settings *models.BlueprintSettings) (core.PipelinePlan, errors.Error) {
	connections := make([]*core.BlueprintConnectionV100, 0)
	err := errors.Convert(json.Unmarshal(settings.Connections, &connections))
	if err != nil {
		return nil, err
	}
	hasDoraEnrich := false
	doraRules := make(map[string]interface{})
	plans := make([]core.PipelinePlan, len(connections))
	for i, connection := range connections {
		if len(connection.Scope) == 0 {
			return nil, errors.Default.New(fmt.Sprintf("connections[%d].scope is empty", i))
		}
		plugin, err := core.GetPlugin(connection.Plugin)
		if err != nil {
			return nil, err
		}
		if pluginBp, ok := plugin.(core.PluginBlueprintV100); ok {
			plans[i], err = pluginBp.MakePipelinePlan(connection.ConnectionId, connection.Scope)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Default.New(fmt.Sprintf("plugin %s does not support blueprint protocol version 1.0.0", connection.Plugin))
		}
		for _, stage := range plans[i] {
			for _, task := range stage {
				if task.Plugin == "dora" {
					hasDoraEnrich = true
					for k, v := range task.Options {
						doraRules[k] = v
					}
				}
			}
		}
	}
	mergedPipelinePlan := ParallelizePipelinePlans(plans...)
	if hasDoraEnrich {
		plan := core.PipelineStage{
			&core.PipelineTask{
				Plugin:   "dora",
				Subtasks: []string{"calculateChangeLeadTime", "ConnectIssueDeploy"},
				Options:  doraRules,
			},
		}
		mergedPipelinePlan = append(mergedPipelinePlan, plan)
	}
	return mergedPipelinePlan, nil
}
