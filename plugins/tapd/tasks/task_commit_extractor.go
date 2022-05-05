package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractTaskCommits

var ExtractTaskCommitMeta = core.SubTaskMeta{
	Name:             "extractTaskCommits",
	EntryPoint:       ExtractTaskCommits,
	EnabledByDefault: true,
	Description:      "Extract raw TaskCommits data into tool layer table _tool_tapd_issue_commits",
}

func ExtractTaskCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_TASK_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var issueCommitBody models.TapdTaskCommit
			err := json.Unmarshal(row.Data, &issueCommitBody)
			if err != nil {
				return nil, err
			}
			toolL := issueCommitBody
			toolL.ConnectionId = data.Connection.ID
			issue := SimpleTask{}
			err = json.Unmarshal(row.Input, &issue)
			if err != nil {
				return nil, err
			}
			toolL.TaskId = issue.Id
			toolL.WorkspaceID = data.Options.WorkspaceID
			results := make([]interface{}, 0, 1)
			results = append(results, &toolL)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
