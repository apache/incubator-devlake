package testhelper

import (
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
			ProjectId: 3472737,
		},
	}

	// import raw data table
	dataflowTester.ImportCsv("./tables/_raw_gitlab_api_projects.csv", "_raw_gitlab_api_project")

	// verify extraction
	dataflowTester.FlushTable("_tool_gitlab_projects")
	dataflowTester.Subtask(tasks.ExtractProjectMeta, taskData)
	dataflowTester.VerifyTable(
		"_tool_gitlab_projects",
		"tables/_tool_gitlab_projects.csv",
		[]string{"gitlab_id"},
		[]string{
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
