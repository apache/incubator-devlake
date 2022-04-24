package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugChangelog

var ExtractBugChangelogMeta = core.SubTaskMeta{
	Name:             "extractBugChangelog",
	EntryPoint:       ExtractBugChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_bug_changelogs",
}

type TapdBugChangelogRes struct {
	BugChange models.TapdBugChangelogApiRes
}

func ExtractBugChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
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
			bugChangelogRes := bugChangelogBody.BugChange

			i, err := VoToDTO(&bugChangelogRes, &models.TapdBugChangelog{})
			if err != nil {
				return nil, err
			}
			v := i.(*models.TapdBugChangelog)
			v.SourceId = data.Source.ID
			v.WorkspaceId = data.Source.WorkspaceId
			item := &models.TapdBugChangelogItem{
				SourceId:          data.Source.ID,
				ChangelogId:       v.ID,
				Field:             v.Field,
				ValueBeforeParsed: v.OldValue,
				ValueAfterParsed:  v.NewValue,
			}
			if item.Field == "iteration_id" {
				iterationFrom, iterationTo, err := parseIterationChangelog(taskCtx, item.ValueBeforeParsed, item.ValueAfterParsed)
				if err != nil {
					return nil, err
				}
				item.IterationIdFrom = iterationFrom
				item.IterationIdTo = iterationTo
			}
			results = append(results, v, item)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
