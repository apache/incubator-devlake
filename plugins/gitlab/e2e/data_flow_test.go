package e2e

import (
	"testing"

	"github.com/merico-dev/lake/plugins/gitlab/impl"
	"github.com/merico-dev/lake/plugins/gitlab/tasks"
	"github.com/merico-dev/lake/testhelper"
)

func TestGitlabDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := testhelper.NewDataFlowTester(t, gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ProjectId: 3472737,
		},
	}
	dataflowTester.ImportCsv("./rawdata/_raw_gitlab_api_projects.csv", "_raw_gitlab_api_project")
	dataflowTester.Subtask(tasks.ExtractProjectMeta, taskData)
	dataflowTester.VerifyTable(
		"_tool_gitlab_projects",
		"rawdata/_tool_gitlab_projects.csv",
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
