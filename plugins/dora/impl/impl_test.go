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

package impl

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeMetricPluginPipelinePlanV200(t *testing.T) {
	var dora Dora
	const projectName = "TestMakePlanV200-project"
	// mock dora plugin as a metric plugin
	option := map[string]interface{}{
		"projectName": projectName,
	}

	optionJson, err := json.Marshal(option)
	assert.Nil(t, err)
	plan, err := dora.MakeMetricPluginPipelinePlanV200(projectName, optionJson)
	doraOutputPlan := core.PipelinePlan{
		core.PipelineStage{
			{
				Plugin:   "refdiff",
				Subtasks: []string{"calculateDeploymentDiffs"},
				Options:  map[string]interface{}{"projectName": projectName},
			},
		},
		core.PipelineStage{
			{
				Plugin:  "dora",
				Options: map[string]interface{}{"projectName": projectName},
			},
		},
	}
	assert.Equal(t, doraOutputPlan, plan)
}
