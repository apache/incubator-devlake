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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/issue_trace/api"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/issue_trace/tasks"
	"github.com/mitchellh/mapstructure"
)

type IssueTrace struct{}

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMetric
	plugin.PluginMigration
	plugin.PluginApi
	plugin.MetricPluginBlueprintV200
} = (*IssueTrace)(nil)

func (p IssueTrace) Name() string {
	return "issue_trace"
}

func (p IssueTrace) RequiredDataEntities() (data []map[string]interface{}, err errors.Error) {
	return []map[string]interface{}{
		{
			"model": "issue_changelogs",
			"requiredFields": map[string]string{
				"column":        "type",
				"execptedValue": "Issue",
			},
		},
	}, nil
}

func (p IssueTrace) IsProjectMetric() bool {
	return true
}

func (p IssueTrace) RunAfter() ([]string, errors.Error) {
	return []string{}, nil
}

func (p IssueTrace) Settings() interface{} {
	return nil
}

func (p IssueTrace) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p IssueTrace) Description() string {
	return "To enrich data from issue tracking domain"
}

// Register all subtasks
func (p IssueTrace) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		// issue_status_history
		tasks.ConvertIssueStatusHistoryMeta,
		// issue_assignee_history
		tasks.ConvertIssueAssigneeHistoryMeta,
	}
}

// Prepare your apiClient which will be used to request remote api,
// `apiClient` is defined in `client.go` under `tasks`
// `SprintPerformanceEnricherTaskData` is defined in `task_data.go` under `tasks`
func (p IssueTrace) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.Options
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Failed to decode options")
	}
	var scopeIds []string
	if op.ScopeIds != nil {
		scopeIds = op.ScopeIds
	} else {
		db := taskCtx.GetDal()
		pmClauses := []dal.Clause{
			dal.From("project_mapping pm"),
			dal.Where("pm.project_name = ? and pm.table = ?", op.ProjectName, "boards"),
		}
		pm := []crossdomain.ProjectMapping{}
		err = db.All(&pm, pmClauses...)
		if err != nil {
			return nil, errors.Default.Wrap(err, "Failed to get project mapping")
		}
		for _, p := range pm {
			scopeIds = append(scopeIds, p.RowId)
		}
	}

	var taskData = &tasks.TaskData{
		Options:     op,
		ScopeIds:    scopeIds,
		ProjectName: op.ProjectName,
	}

	return taskData, nil
}

func (p IssueTrace) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/issue_trace"
}

func (p IssueTrace) MigrationScripts() []plugin.MigrationScript {
	return []plugin.MigrationScript{
		&migrationscripts.NewIssueTable{},
	}
}

func (p IssueTrace) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{}
}

func (p IssueTrace) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.IssueAssigneeHistory{},
		&models.IssueStatusHistory{},
	}
}

func (p IssueTrace) MakeMetricPluginPipelinePlanV200(projectName string, options json.RawMessage) (coreModels.PipelinePlan, errors.Error) {
	op := &tasks.Options{}
	if options != nil && string(options) != "\"\"" {
		err := json.Unmarshal(options, op)
		if err != nil {
			return nil, errors.Default.WrapRaw(err)
		}
	}

	plan := coreModels.PipelinePlan{
		{
			{
				Plugin: "issue_trace",
				Options: map[string]interface{}{
					"projectName": projectName,
					"scopeIds":    op.ScopeIds,
				},
				Subtasks: []string{
					"ConvertIssueStatusHistory",
					"ConvertIssueAssigneeHistory",
				},
			},
		},
	}
	return plan, nil
}
