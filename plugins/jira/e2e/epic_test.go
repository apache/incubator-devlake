package e2e

import (
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEpicDataflow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)
	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId:        1,
			BoardId:             93,
			TransformationRules: tasks.TransformationRules{StoryPointField: "customfield_10024"},
		},
	}

	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issue_types.csv", "_raw_jira_api_issue_types")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issues.csv", "_raw_jira_api_issues")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_external_epics.csv", "_raw_jira_api_epics")

	dataflowTester.FlushTabler(&models.JiraIssue{})
	dataflowTester.FlushTabler(&models.JiraBoardIssue{})
	dataflowTester.FlushTabler(&models.JiraSprintIssue{})
	dataflowTester.FlushTabler(&models.JiraIssueChangelogs{})
	dataflowTester.FlushTabler(&models.JiraIssueChangelogItems{})
	dataflowTester.FlushTabler(&models.JiraWorklog{})
	dataflowTester.FlushTabler(&models.JiraAccount{})
	dataflowTester.FlushTabler(&models.JiraIssueType{})

	ctx := dataflowTester.SubtaskContext(taskData)

	// run pre-req subtasks
	require.NoError(t, tasks.ExtractIssueTypesMeta.EntryPoint(ctx))
	require.NoError(t, tasks.ExtractIssuesMeta.EntryPoint(ctx))
	dataflowTester.VerifyTableWithOptions(
		models.JiraIssue{}, e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_jira_issues_for_external_epics.csv",
			TargetFields: nil,
			IgnoreFields: nil,
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
		},
	)
	dataflowTester.VerifyTableWithOptions(
		models.JiraBoardIssue{}, e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_jira_board_issues_for_external_epics.csv",
			TargetFields: []string{"connection_id", "board_id", "issue_id"},
			IgnoreFields: nil,
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
		},
	)

	// run the part of the collector that queries tools data
	keys, err := tasks.GetEpicKeys(ctx.GetDal(), taskData)
	require.NoError(t, err)
	require.Contains(t, keys, "K5-1") //epic not on the board

	require.NoError(t, tasks.ExtractEpicsMeta.EntryPoint(ctx))

	dataflowTester.VerifyTableWithOptions(
		models.JiraBoardIssue{}, e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_jira_board_issues_for_external_epics.csv",
			TargetFields: nil,
			IgnoreFields: nil,
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
		},
	)
	dataflowTester.VerifyTableWithOptions(
		models.JiraIssue{}, e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_jira_issues_for_external_epics.csv",
			TargetFields: nil,
			IgnoreFields: nil,
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
		},
	)
}
