package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractTaskChangelog

var ExtractTaskChangelogMeta = core.SubTaskMeta{
	Name:             "extractTaskChangelog",
	EntryPoint:       ExtractTaskChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdTaskChangelogRes struct {
	WorkitemChange models.TapdTaskChangelogApiRes
}

func ExtractTaskChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_TASK_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var taskChangelogBody TapdTaskChangelogRes
			err := json.Unmarshal(row.Data, &taskChangelogBody)
			if err != nil {
				return nil, err
			}
			taskChangelogRes := taskChangelogBody.WorkitemChange

			i, err := VoToDTO(&taskChangelogRes, &models.TapdTaskChangelog{})
			if err != nil {
				return nil, err
			}
			v := i.(*models.TapdTaskChangelog)
			toolL := &models.TapdChangelog{
				SourceId:       data.Source.ID,
				ID:             v.ID,
				WorkspaceId:    v.WorkspaceId,
				WorkitemTypeID: v.WorkitemTypeID,
				Creator:        v.Creator,
				Created:        v.Created,
				ChangeSummary:  v.ChangeSummary,
				Comment:        v.Comment,
				EntityType:     v.EntityType,
				ChangeType:     v.EntityType,
				IssueId:        v.TaskID,
			}
			results := make([]interface{}, 0, 1)

			for _, item := range v.FieldChanges {
				item := &models.TapdChangelogItem{
					SourceId:          data.Source.ID,
					ChangelogId:       toolL.ID,
					Field:             item.Field,
					ValueBeforeParsed: item.ValueBeforeParsed,
					ValueAfterParsed:  item.ValueAfterParsed,
				}
				results = append(results, item)
			}
			results = append(results, toolL)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
