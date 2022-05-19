package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugChangelog

var ExtractBugChangelogMeta = core.SubTaskMeta{
	Name:             "extractBugChangelog",
	EntryPoint:       ExtractBugChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_bug_changelogs",
}

type TapdBugChangelogRes struct {
	BugChange models.TapdBugChangelog
}

func ExtractBugChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			results := make([]interface{}, 0, 2)
			var bugChangelogBody TapdBugChangelogRes
			err := json.Unmarshal(row.Data, &bugChangelogBody)
			if err != nil {
				return nil, err
			}
			bugChangelog := bugChangelogBody.BugChange

			bugChangelog.ConnectionId = data.Connection.ID
			bugChangelog.WorkspaceID = data.Options.WorkspaceID
			item := &models.TapdBugChangelogItem{
				ConnectionId:      data.Connection.ID,
				ChangelogId:       bugChangelog.ID,
				Field:             bugChangelog.Field,
				ValueBeforeParsed: bugChangelog.OldValue,
				ValueAfterParsed:  bugChangelog.NewValue,
			}
			if item.Field == "iteration_id" {
				iterationFrom, iterationTo, err := parseIterationChangelog(taskCtx, item.ValueBeforeParsed, item.ValueAfterParsed)
				if err != nil {
					return nil, err
				}
				item.IterationIdFrom = iterationFrom
				item.IterationIdTo = iterationTo
			}
			results = append(results, &bugChangelog, item)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
