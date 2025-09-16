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

package api

import (
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	helperapi "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/slack/tasks"
)

func MakeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// Map blueprint scopes to actual Slack channels
	scopeDetails, err := dsHelper.ScopeSrv.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}
	// Build one stage per selected channel
	plan := make(coreModels.PipelinePlan, len(scopeDetails))
	for i, scopeDetail := range scopeDetails {
		stage := plan[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}
		// Only include CROSS domain subtasks; Slack subtasks define DomainTypes accordingly.
		entities := []string{plugin.DOMAIN_TYPE_CROSS}
		scope := scopeDetail.Scope // *models.SlackChannel
		task, err := helperapi.MakePipelinePlanTask(
			"slack",
			subtaskMetas,
			entities,
			tasks.SlackOptions{ConnectionId: connectionId, ChannelId: scope.ScopeId()},
		)
		if err != nil {
			return nil, nil, err
		}
		stage = append(stage, task)
		plan[i] = stage
	}
	// No domain scopes emitted by Slack for now
	return plan, nil, nil
}
