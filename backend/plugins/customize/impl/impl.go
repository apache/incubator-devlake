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
	"github.com/apache/incubator-devlake/plugins/customize/api"
	"github.com/apache/incubator-devlake/plugins/customize/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/customize/tasks"
	"github.com/mitchellh/mapstructure"
)

var _ plugin.PluginMeta = (*Customize)(nil)
var _ plugin.PluginInit = (*Customize)(nil)
var _ plugin.PluginApi = (*Customize)(nil)
var _ plugin.PluginModel = (*Customize)(nil)
var _ plugin.PluginMigration = (*Customize)(nil)

type Customize struct {
	handlers *api.Handlers
}

func (p *Customize) Init(basicRes context.BasicRes) errors.Error {
	p.handlers = api.NewHandlers(basicRes.GetDal())
	return nil
}

func (p Customize) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (p Customize) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ExtractCustomizedFieldsMeta,
	}
}

func (p Customize) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
}
func (p Customize) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.Options
	var err error
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not decode Jira options")
	}
	taskData := &tasks.TaskData{
		Options: &op,
	}
	return taskData, nil
}

func (p Customize) Description() string {
	return "To customize table fields"
}

func (p Customize) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Customize) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/customize"
}

func (p *Customize) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		":table/fields": {
			"GET":  p.handlers.ListFields,
			"POST": p.handlers.CreateFields,
		},
		":table/fields/:field": {
			"DELETE": p.handlers.DeleteField,
		},
		"csvfiles/issues.csv": {
			"POST": p.handlers.ImportIssue,
		},
		"csvfiles/issue_commits.csv": {
			"POST": p.handlers.ImportIssueCommit,
		},
		"csvfiles/issue_repo_commits.csv": {
			"POST": p.handlers.ImportIssueRepoCommit,
		},
	}
}
