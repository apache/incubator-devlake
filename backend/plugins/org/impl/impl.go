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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/org/api"
	"github.com/apache/incubator-devlake/plugins/org/tasks"
)

var _ plugin.PluginMeta = (*Org)(nil)
var _ plugin.PluginInit = (*Org)(nil)
var _ plugin.PluginTask = (*Org)(nil)
var _ plugin.PluginModel = (*Org)(nil)
var _ plugin.ProjectMapper = (*Org)(nil)

type Org struct {
	handlers *api.Handlers
}

func (p *Org) Init(basicRes context.BasicRes) errors.Error {
	p.handlers = api.NewHandlers(basicRes)
	return nil
}

func (p Org) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (p Org) Description() string {
	return "collect data related to team and organization"
}

func (p Org) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ConnectUserAccountsExactMeta,
		tasks.SetProjectMappingMeta,
	}
}

func (p Org) MapProject(projectName string, scopes []plugin.Scope) (plugin.PipelinePlan, errors.Error) {
	var plan plugin.PipelinePlan
	var stage plugin.PipelineStage

	// construct task options for Org
	options := make(map[string]interface{})
	options["projectMappings"] = []tasks.ProjectMapping{tasks.NewProjectMapping(projectName, scopes)}

	subtasks, err := helper.MakePipelinePlanSubtasks([]plugin.SubTaskMeta{tasks.SetProjectMappingMeta}, []string{plugin.DOMAIN_TYPE_CROSS})
	if err != nil {
		return nil, err
	}
	stage = append(stage, &plugin.PipelineTask{
		Plugin:   "org",
		Subtasks: subtasks,
		Options:  options,
	})
	plan = append(plan, stage)
	return plan, nil
}

func (p Org) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.Options
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode options")
	}
	taskData := &tasks.TaskData{
		Options: &op,
	}
	return taskData, nil
}

func (p Org) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/org"
}

func (p Org) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"teams.csv": {
			"GET": p.handlers.GetTeam,
			"PUT": p.handlers.CreateTeam,
		},
		"users.csv": {
			"GET": p.handlers.GetUser,
			"PUT": p.handlers.CreateUser,
		},

		"user_account_mapping.csv": {
			"GET": p.handlers.GetUserAccountMapping,
			"PUT": p.handlers.CreateUserAccountMapping,
		},
		"project_mapping.csv": {
			"GET": p.handlers.GetProjectMapping,
			"PUT": p.handlers.CreateProjectMapping,
		},
	}
}
