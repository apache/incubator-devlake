package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractStoryCommits

var ExtractStoryCommitMeta = core.SubTaskMeta{
	Name:             "extractStoryCommits",
	EntryPoint:       ExtractStoryCommits,
	EnabledByDefault: true,
	Description:      "Extract raw StoryCommits data into tool layer table _tool_tapd_issue_commits",
}

func ExtractStoryCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_STORY_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var issueCommitBody models.TapdStoryCommit
			err := json.Unmarshal(row.Data, &issueCommitBody)
			if err != nil {
				return nil, err
			}
			toolL := issueCommitBody
			toolL.ConnectionId = data.Connection.ID
			issue := SimpleStory{}
			err = json.Unmarshal(row.Input, &issue)
			if err != nil {
				return nil, err
			}
			toolL.StoryId = issue.Id
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
