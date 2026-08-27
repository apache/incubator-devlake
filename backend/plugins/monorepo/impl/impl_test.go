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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Core's blueprint plan builder (server/services/blueprint_makeplan_v200.go) merges every
// enabled metric plugin's plan with ParallelizePipelinePlans, which zips stages together by
// index and does NOT consult RunAfter(). If monorepo's real work sat in stage 0 like a naive
// single-stage plan would, it would run concurrently with dora's stage 0 instead of after
// dora's stage 2 (where project_pr_metrics/cicd_deployment_commits actually get written) —
// silently producing nil-metric output. This test locks in the stage-padding workaround so a
// future edit can't accidentally collapse the plan back to one stage.
func TestMakeMetricPluginPipelinePlanV200_StagePadding(t *testing.T) {
	options, err := json.Marshal(map[string]interface{}{
		"subProjects": []map[string]interface{}{
			{"name": "serviceA", "prLabels": []string{"serviceA"}, "deployJobPattern": "^deploy-serviceA$"},
		},
	})
	require.NoError(t, err)

	var p Monorepo
	plan, err2 := p.MakeMetricPluginPipelinePlanV200("test-project", options)
	require.NoError(t, err2)

	require.Greater(t, len(plan), 3, "plan must have more stages than dora's plan (3), or monorepo's "+
		"work would run concurrently with dora instead of after it")

	for i := 0; i < len(plan)-1; i++ {
		assert.Emptyf(t, plan[i], "stage %d should be empty padding, not real work", i)
	}

	lastStage := plan[len(plan)-1]
	require.Len(t, lastStage, 1)
	assert.Equal(t, "monorepo", lastStage[0].Plugin)
	assert.ElementsMatch(t, []string{
		"attributeDeployments", "attributePullRequests", "updateProjectPrMetricsSubProject",
	}, lastStage[0].Subtasks)
	assert.Equal(t, true, lastStage[0].Options["includeUnattributed"],
		"includeUnattributed must default to true when the caller doesn't set it")
}

// TestMakeMetricPluginPipelinePlanV200_IncludeUnattributedExplicitFalse locks in that an
// explicit `"includeUnattributed": false` survives into the task options unchanged,
// distinguishing it from "not set" (which defaults to true).
func TestMakeMetricPluginPipelinePlanV200_IncludeUnattributedExplicitFalse(t *testing.T) {
	options, err := json.Marshal(map[string]interface{}{
		"subProjects": []map[string]interface{}{
			{"name": "serviceA", "prLabels": []string{"serviceA"}, "deployJobPattern": "^deploy-serviceA$"},
		},
		"includeUnattributed": false,
	})
	require.NoError(t, err)

	var p Monorepo
	plan, err2 := p.MakeMetricPluginPipelinePlanV200("test-project", options)
	require.NoError(t, err2)

	lastStage := plan[len(plan)-1]
	assert.Equal(t, false, lastStage[0].Options["includeUnattributed"])
}
