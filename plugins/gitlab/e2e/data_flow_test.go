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

	dataflowTester.ImportCsv("./rawdata/_raw_gitlab_api_project.csv", "_raw_gitlab_api_project")
	dataflowTester.Subtask(tasks.ExtractProjectMeta)
	dataflowTester.VerifyTable("_tool_gitlab_project", "rawdata/_tool_gitlab_projects.csv")
}
