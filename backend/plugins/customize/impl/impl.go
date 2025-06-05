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
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/customize/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/customize/tasks"
	"github.com/mitchellh/mapstructure"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMigration
} = (*Customize)(nil)

var handlers *api.Handlers

type Customize struct {
}

func (p Customize) Init(basicRes context.BasicRes) errors.Error {
	handlers = api.NewHandlers(basicRes.GetDal())
	return nil
}

func (p Customize) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.CustomizedField{},
	}
}

func (p Customize) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ExtractCustomizedFieldsMeta,
	}
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

func (p Customize) Name() string {
	return "customize"
}

func (p Customize) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Customize) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/customize"
}

func (p Customize) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		":table/fields": {
			"GET":  handlers.ListFields,
			"POST": handlers.CreateFields,
		},
		":table/fields/:field": {
			"DELETE": handlers.DeleteField,
		},
		"csvfiles/issues.csv": {
			"POST": handlers.ImportIssue,
		},
		"csvfiles/issue_commits.csv": {
			"POST": handlers.ImportIssueCommit,
		},
		"csvfiles/issue_repo_commits.csv": {
			"POST": handlers.ImportIssueRepoCommit,
		},
		"csvfiles/issue_changelogs.csv": {
			"POST": handlers.ImportIssueChangelog,
		},
		"csvfiles/issue_worklogs.csv": {
			"POST": handlers.ImportIssueWorklog,
		},
		"csvfiles/sprints.csv": {
			"POST": handlers.ImportSprint,
		},
		"csvfiles/qa_apis.csv": {
			"POST": handlers.ImportQaApis,
		},
		"csvfiles/qa_test_cases.csv": {
			"POST": handlers.ImportQaTestCases,
		},
		"csvfiles/qa_test_case_executions.csv": {
			"POST": handlers.ImportQaTestCaseExecutions,
		},
	}
}
