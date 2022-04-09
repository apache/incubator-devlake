package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugChangelog

var ExtractBugChangelogMeta = core.SubTaskMeta{
	Name:             "extractBugChangelog",
	EntryPoint:       ExtractBugChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
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
			toolL := &models.TapdChangelog{
				SourceId:      data.Source.ID,
				ID:            v.ID,
				WorkspaceId:   v.WorkspaceId,
				Creator:       v.Author,
				Created:       v.Created,
				ChangeSummary: v.Memo,
				EntityType:    "BUG",
				ChangeType:    "BUG",
				IssueId:       v.BugID,
			}
			results := make([]interface{}, 0, 2)

			item := &models.TapdChangelogItem{
				SourceId:          data.Source.ID,
				ChangelogId:       v.ID,
				Field:             v.Field,
				ValueBeforeParsed: v.OldValue,
				ValueAfterParsed:  v.NewValue,
				NoPKModel:         common.NoPKModel{},
			}
			if item.Field == "iteration_id" {
				item, err = parseIterationChangelog(taskCtx, item)
				if err != nil {
					return nil, err
				}
			}
			results = append(results, toolL, item)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
