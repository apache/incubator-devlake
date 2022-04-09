package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractIssueCommits

var ExtractIssueCommitMeta = core.SubTaskMeta{
	Name:             "extractIssueCommits",
	EntryPoint:       ExtractIssueCommits,
	EnabledByDefault: true,
	Description:      "Extract raw IssueCommits data into tool layer table _tool_tapd_issue_commits",
}

func ExtractIssueCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_ISSUE_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var issueCommitBody models.TapdIssueCommitApiRes
			err := json.Unmarshal(row.Data, &issueCommitBody)
			if err != nil {
				return nil, err
			}

			i, err := VoToDTO(&issueCommitBody, &models.TapdIssueCommit{})
			if err != nil {
				return nil, err
			}
			toolL := i.(*models.TapdIssueCommit)
			toolL.SourceId = data.Source.ID
			issue := &models.IssueTypeAndId{}
			err = json.Unmarshal(row.Input, issue)
			if err != nil {
				return nil, err
			}
			toolL.IssueId = issue.IssueId
			toolL.IssueType = issue.Type
			toolL.WorkspaceId = data.Options.WorkspaceId
			results := make([]interface{}, 0, 1)
			results = append(results, toolL)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
