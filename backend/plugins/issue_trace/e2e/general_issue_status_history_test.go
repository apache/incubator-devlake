package e2e

import (
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/issue_trace/impl"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
	"github.com/apache/incubator-devlake/plugins/issue_trace/tasks"
)

func TestConvertIssueStatusHistory(t *testing.T) {
	var plugin impl.IssueTrace
	dataflowTester := e2ehelper.NewDataFlowTester(t, "issue_trace", plugin)
	dataflowTester.ImportCsvIntoTabler("./raw_tables/board_issues.csv", &ticket.BoardIssue{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issues.csv", &ticket.Issue{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issue_changelogs.csv", &ticket.IssueChangelogs{})

	dataflowTester.FlushTabler(models.IssueStatusHistory{})

	dataflowTester.Subtask(tasks.ConvertIssueStatusHistoryMeta, TaskData)

	dataflowTester.VerifyTable(
		models.IssueStatusHistory{},
		"./snapshot_tables/issue_status_history.csv",
		[]string{
			"issue_id",
			"status",
			"original_status",
			"start_date",
			"is_current_status",
			"is_first_status",
		},
	)
}
