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

package e2ehelper

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/github/models"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func ExampleDataFlowTester() {
	var t *testing.T // stub

	var gitlab core.PluginMeta
	dataflowTester := NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ProjectId: 666888,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./tables/_raw_gitlab_api_projects.csv", "_raw_gitlab_api_project")

	// verify extraction
	dataflowTester.FlushTabler(gitlabModels.GitlabProject{})
	dataflowTester.Subtask(tasks.ExtractProjectMeta, taskData)
	dataflowTester.VerifyTable(
		gitlabModels.GitlabProject{},
		"tables/_tool_gitlab_projects.csv",
		[]string{
			"gitlab_id",
			"name",
			"description",
			"default_branch",
			"path_with_namespace",
			"web_url",
			"creator_id",
			"visibility",
			"open_issues_count",
			"star_count",
			"forked_from_project_id",
			"forked_from_project_web_url",
			"created_date",
			"updated_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}

func TestGetTableMetaData(t *testing.T) {
	var gitlab core.PluginMeta
	dataflowTester := NewDataFlowTester(t, "github", gitlab)

	t.Run("get_fields", func(t *testing.T) {
		fields := dataflowTester.getFields(&models.GithubIssueLabel{}, func(column gorm.ColumnType) bool {
			return true
		})
		assert.Equal(t, 9, len(fields))
		for _, e := range []string{
			"connection_id",
			"issue_id",
			"label_name",
			"created_at",
			"updated_at",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("extract_columns", func(t *testing.T) {
		columns := dataflowTester.extractColumns(&common.NoPKModel{})
		assert.Equal(t, 4, len(columns))
		for _, e := range []string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		} {
			assert.Contains(t, columns, e)
		}
	})
	t.Run("get_pk_fields", func(t *testing.T) {
		fields := dataflowTester.getPkFields(&models.GithubIssueLabel{})
		assert.Equal(t, 3, len(fields))
		for _, e := range []string{
			"connection_id",
			"issue_id",
			"label_name",
		} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_targetFieldsOnly", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: []string{"connection_id"},
			IgnoreFields: nil,
			IgnoreTypes:  nil,
		})
		assert.Equal(t, 1, len(fields))
		for _, e := range []string{"connection_id"} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_ignoreFieldsOnly", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: nil,
			IgnoreFields: []string{
				"label_name",
				"created_at",
				"updated_at",
				"_raw_data_params",
				"_raw_data_table",
				"_raw_data_id",
				"_raw_data_remark",
			},
			IgnoreTypes: nil,
		})
		assert.Equal(t, 2, len(fields))
		for _, e := range []string{"connection_id", "issue_id"} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_ignoreFieldsOnly", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: nil,
			IgnoreFields: []string{
				"label_name",
				"created_at",
				"updated_at",
				"_raw_data_params",
				"_raw_data_table",
				"_raw_data_id",
				"_raw_data_remark",
			},
			IgnoreTypes: nil,
		})
		assert.Equal(t, 2, len(fields))
		for _, e := range []string{"connection_id", "issue_id"} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_ignoreType", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: nil,
			IgnoreFields: nil,
			IgnoreTypes:  []interface{}{&common.NoPKModel{}},
		})
		assert.Equal(t, 3, len(fields))
		for _, e := range []string{
			"connection_id",
			"issue_id",
			"label_name",
		} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_ignoreType_ignoreFields", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: nil,
			IgnoreFields: []string{"label_name"},
			IgnoreTypes:  []interface{}{&common.NoPKModel{}},
		})
		assert.Equal(t, 2, len(fields))
		for _, e := range []string{
			"connection_id",
			"issue_id",
		} {
			assert.Contains(t, fields, e)
		}
	})
	t.Run("resolve_fields_targetFields_ignoreType_ignoreFields", func(t *testing.T) {
		fields := dataflowTester.resolveTargetFields(&models.GithubIssueLabel{}, TableOptions{
			TargetFields: []string{"label_name", "createdAt", "connection_id"},
			IgnoreFields: []string{"label_name"},
			IgnoreTypes:  []interface{}{&common.NoPKModel{}},
		})
		assert.Equal(t, 1, len(fields))
		for _, e := range []string{
			"connection_id",
		} {
			assert.Contains(t, fields, e)
		}
	})
}
